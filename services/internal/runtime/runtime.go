package runtime

import (
	"encoding/json"
	"net/http"
)

type Config struct {
	Role    string
	Version string
}

type statusResponse struct {
	Role    string `json:"role"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

func NewMux(config Config) *http.ServeMux {
	mux := http.NewServeMux()
	writeStatus := func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.Header().Set("Allow", http.MethodGet)
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(statusResponse{
			Role:    config.Role,
			Status:  "ok",
			Version: config.Version,
		})
	}
	mux.HandleFunc("/healthz", writeStatus)
	mux.HandleFunc("/version", writeStatus)
	return mux
}

func Serve(config Config, address string) error {
	server := &http.Server{
		Addr:    address,
		Handler: NewMux(config),
	}
	return server.ListenAndServe()
}
