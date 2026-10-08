package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"calculator/internal/calculator"
	"calculator/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}

	calcService := calculator.NewCalculatorService()
	calcHandler := handler.NewHandler(calcService, allowedOrigin)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/calculate", calcHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Calculator API running on port %s (Allowed Origin: %s)\n", port, allowedOrigin)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
