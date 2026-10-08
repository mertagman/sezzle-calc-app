package server

import (
	"log/slog"
	"net/http"

	"calculator/internal/handler"
)

func NewRouter(h *handler.Handler, allowedOrigin string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/add", h.Add)
	mux.HandleFunc("POST /api/v1/subtract", h.Subtract)
	mux.HandleFunc("POST /api/v1/multiply", h.Multiply)
	mux.HandleFunc("POST /api/v1/divide", h.Divide)
	mux.HandleFunc("POST /api/v1/power", h.Power)
	mux.HandleFunc("POST /api/v1/sqrt", h.Sqrt)
	mux.HandleFunc("POST /api/v1/percentage", h.Percentage)
	mux.HandleFunc("GET /health", healthHandler)

	return recoveryMiddleware(corsMiddleware(allowedOrigin, loggingMiddleware(mux)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "error", rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("http request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			if origin != "" && origin == allowedOrigin {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
