package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xOMAIKOx/gridworkx/services/api/internal/api"
	"github.com/xOMAIKOx/gridworkx/services/api/internal/apiruntime"
	"github.com/xOMAIKOx/gridworkx/services/api/internal/postgres"
)

func main() {
	defaultAddress := apiruntime.BindAddress(os.Getenv("GRIDWORKS_API_BIND"))
	address := flag.String("bind", defaultAddress, "API bind address")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	repo, err := postgres.Open(ctx, os.Getenv("GRIDWORKS_DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	server := apiruntime.NewServer(*address, api.NewServer(repo, "0.1.0").Mux())
	runCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := apiruntime.Run(runCtx, server); err != nil {
		log.Fatal(err)
	}
}
