package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"calculator/internal/calculator"
)

type CalculatorService interface {
	Calculate(operation string, a float64, b float64) (float64, error)
}

type CalculateRequest struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

type CalculateResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	service       CalculatorService
	allowedOrigin string
}

func NewHandler(service CalculatorService, allowedOrigin string) *Handler {
	return &Handler{
		service:       service,
		allowedOrigin: allowedOrigin,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if origin != "" && origin != h.allowedOrigin {
		h.writeError(w, http.StatusForbidden, "cors origin not allowed")
		return
	}

	if origin == h.allowedOrigin {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Vary", "Origin")
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 1 MB limit (DoS korumasi)
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req CalculateRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	if req.Operation == "" {
		h.writeError(w, http.StatusBadRequest, "operation field is required")
		return
	}

	result, err := h.service.Calculate(req.Operation, req.A, req.B)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, calculator.ErrDivisionByZero) ||
			errors.Is(err, calculator.ErrNegativeSqrt) ||
			errors.Is(err, calculator.ErrUndefinedResult) {
			status = http.StatusUnprocessableEntity
		}
		h.writeError(w, status, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(CalculateResponse{Result: result})
}

func (h *Handler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
