package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	serviceName = "billing"
	defaultPort = "8088"
)

var startedAt = time.Now()

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Time    string `json:"time"`
}

type statusResponse struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Status  string `json:"status"`
	Uptime  string `json:"uptime"`
}

type resourceResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	State    string `json:"state"`
	Provider string `json:"provider"`
}

func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{
			Status:  "ok",
			Service: serviceName,
			Time:    time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func readyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{
			Status:  "ready",
			Service: serviceName,
			Time:    time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func statusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, statusResponse{
			Service: serviceName,
			Version: "0.1.0",
			Status:  "skeleton",
			Uptime:  time.Since(startedAt).Round(time.Second).String(),
		})
	}
}

func listResourceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := []resourceResponse{
			{
				ID:       serviceName + "-0001",
				Name:     serviceName + "-demo-invoices",
				Region:   "br-sul-1",
				State:    "available",
				Provider: serviceName,
			},
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"items": items,
			"count": len(items),
		})
	}
}

func getResourceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/invoices/")
		if id == "" || strings.Contains(id, "/") {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid identifier"})
			return
		}
		writeJSON(w, http.StatusOK, resourceResponse{
			ID:       id,
			Name:     serviceName + "-demo-invoices",
			Region:   "br-sul-1",
			State:    "available",
			Provider: serviceName,
		})
	}
}

func newServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler())
	mux.HandleFunc("/ready", readyHandler())
	mux.HandleFunc("/api/v1/status", statusHandler())
	mux.HandleFunc("/api/v1/invoices", listResourceHandler())
	mux.HandleFunc("/api/v1/invoices/", getResourceHandler())

	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = defaultPort
	}

	srv := newServer(":" + port)

	go func() {
		log.Printf("[BRLIA-billing] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
