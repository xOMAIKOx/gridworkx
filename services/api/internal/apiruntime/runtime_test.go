package apiruntime

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

type fakeServer struct {
	started     chan struct{}
	stopped     chan struct{}
	serveErr    error
	shutdownErr error
	shutdown    bool
}

func newFakeServer() *fakeServer {
	return &fakeServer{started: make(chan struct{}), stopped: make(chan struct{})}
}
func (f *fakeServer) ListenAndServe() error {
	close(f.started)
	if f.serveErr != nil {
		return f.serveErr
	}
	<-f.stopped
	return http.ErrServerClosed
}
func (f *fakeServer) Shutdown(context.Context) error {
	f.shutdown = true
	if f.shutdownErr != nil {
		return f.shutdownErr
	}
	close(f.stopped)
	return nil
}

func TestServerConfigurationUsesLoopbackAndBoundedLimits(t *testing.T) {
	s := NewServer("", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	if s.Addr != DefaultBind {
		t.Fatalf("bind=%q", s.Addr)
	}
	if s.ReadHeaderTimeout <= 0 || s.ReadTimeout <= 0 || s.WriteTimeout <= 0 || s.IdleTimeout <= 0 || s.MaxHeaderBytes <= 0 {
		t.Fatalf("unbounded server config: %+v", s)
	}
	if got := BindAddress("127.0.0.1:19000"); got != "127.0.0.1:19000" {
		t.Fatalf("override=%q", got)
	}
}

func TestRunGracefullyShutsDownWithoutPublicPort(t *testing.T) {
	server := newFakeServer()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, server) }()
	<-server.started
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("graceful shutdown did not complete")
	}
	if !server.shutdown {
		t.Fatal("shutdown was not invoked")
	}
}

func TestRunSurfacesServeFailure(t *testing.T) {
	want := errors.New("listen failed")
	server := newFakeServer()
	server.serveErr = want
	if err := Run(context.Background(), server); !errors.Is(err, want) {
		t.Fatalf("serve failure=%v", err)
	}
}

func TestRunSurfacesShutdownFailure(t *testing.T) {
	want := errors.New("shutdown failed")
	server := newFakeServer()
	server.shutdownErr = want
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, server) }()
	<-server.started
	cancel()
	if err := <-done; !errors.Is(err, want) {
		t.Fatalf("shutdown failure=%v", err)
	}
}
