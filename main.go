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
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	DBSSLMode  string
	DBTimeZone string
}

func load() *Config {
	return &Config{
		ServerHost: "localhost",
		ServerPort: "8000",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPass:     "secret",
		DBName:     "shortly",
		DBSSLMode:  "disable",
		DBTimeZone: "Asia/Kolkata",
	}
}

type handler struct {
	db *core.PostgresDB
}

func (h *handler) urlShortener(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "url value missing, pass ?url=... to get short url", http.StatusBadRequest)
		return
	}

	res, err := core.ShortenURL(r.Context(), url, h.db)
	if err != nil {
		fmt.Println("[Error] ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"url": res,
	})
}

func main() {
	config := load()

	db, err := core.NewPostgresDB(config.DBHost, config.DBPort,
		config.DBUser, config.DBPass, config.DBName,
		config.DBSSLMode, config.DBTimeZone)
	if err != nil {
		fmt.Println(err)
		return
	}

	h := &handler{db}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /shorten", h.urlShortener)

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
