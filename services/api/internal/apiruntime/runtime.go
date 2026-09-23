package apiruntime

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

const DefaultBind = "127.0.0.1:18080"

type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func BindAddress(override string) string {
	if strings.TrimSpace(override) == "" {
		return DefaultBind
	}
	return override
}

func NewServer(address string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              BindAddress(address),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    32 << 10,
	}
}

func Run(ctx context.Context, server Server) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		err := <-errCh
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
