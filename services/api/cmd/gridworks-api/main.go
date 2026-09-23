package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xOMAIKOx/gridworkx/services/api/internal/api"
	"github.com/xOMAIKOx/gridworkx/services/api/internal/postgres"
)

func main() {
	defaultAddress := os.Getenv("GRIDWORKS_API_BIND")
	if defaultAddress == "" {
		defaultAddress = "127.0.0.1:18080"
	}
	address := flag.String("bind", defaultAddress, "API bind address")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	repo, err := postgres.Open(ctx, os.Getenv("GRIDWORKS_DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	server := &http.Server{Addr: *address, Handler: api.NewServer(repo, "0.1.0").Mux(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 32 << 10}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-stop
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdown); err != nil {
		log.Printf("API shutdown: %v", err)
	}
}
