package main

import (
	"flag"
	"log"
	"os"

	"github.com/xOMAIKOx/gridworkx/services/internal/runtime"
)

func main() {
	defaultAddress := os.Getenv("GRIDWORKS_SIM_VALIDATOR_BIND")
	if defaultAddress == "" {
		defaultAddress = "127.0.0.1:18083"
	}
	address := flag.String("bind", defaultAddress, "loopback bind address")
	flag.Parse()
	if err := runtime.Serve(runtime.Config{Role: "gridworks-sim-validator", Version: "0.1.0"}, *address); err != nil {
		log.Fatal(err)
	}
}
