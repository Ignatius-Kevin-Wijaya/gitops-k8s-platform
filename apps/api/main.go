package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	serviceName string
	version     string
	environment string
	port        string
}

type statusResponse struct {
	Status string `json:"status"`
}

type infoResponse struct {
	Service     string `json:"service"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

func main() {
	cfg := loadConfig()
	metrics := newMetricsRegistry()

	mux := http.NewServeMux()
	mux.Handle("/healthz", instrument("healthz", metrics, http.HandlerFunc(handleHealthz)))
	mux.Handle("/readyz", instrument("readyz", metrics, http.HandlerFunc(handleReadyz)))
	mux.Handle("/api/v1/info", instrument("info", metrics, http.HandlerFunc(handleInfo(cfg))))
	mux.Handle("/metrics", metrics.handler(cfg))

	server := &http.Server{
		Addr:              ":" + cfg.port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go handleShutdown(server)

	log.Printf("starting %s on :%s", cfg.serviceName, cfg.port)

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}

func loadConfig() config {
	return config{
		serviceName: getEnv("SERVICE_NAME", "gitops-api"),
		version:     getEnv("APP_VERSION", "dev"),
		environment: getEnv("APP_ENV", "local"),
		port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func handleShutdown(server *http.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, statusResponse{Status: "ok"})
}

func handleReadyz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, statusResponse{Status: "ready"})
}

func handleInfo(cfg config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, infoResponse{
			Service:     cfg.serviceName,
			Version:     cfg.version,
			Environment: cfg.environment,
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write json response: %v", err)
	}
}
