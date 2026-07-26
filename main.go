package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Arup3201/shorturl/core"
)

type Config struct {
	ServerHost string
	ServerPort string
}

func load() *Config {
	return &Config{
		ServerHost: "localhost",
		ServerPort: "8000",
	}
}

func urlShortener(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "url value missing, pass ?url=... to get short url", http.StatusBadRequest)
		return
	}

	res := core.ShortenURL(url)
	json.NewEncoder(w).Encode(map[string]string{
		"url": res,
	})
}

func main() {
	config := load()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /shorten", urlShortener)

	address := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)
	server := http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	fmt.Println("Starting server...")
	fmt.Printf("Address: %s\n", address)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to start server.\nError: %s\n", err)
	}
}
