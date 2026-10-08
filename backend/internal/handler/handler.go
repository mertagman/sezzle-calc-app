package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"calculator/internal/calculator"
)

type CalculatorService interface {
	Calculate(operation string, a float64, b float64) (float64, error)
	IsUnary(operation string) bool
}

type CalculateRequest struct {
	Operation string   `json:"operation"`
	A         *float64 `json:"a"`
	B         *float64 `json:"b"`
}

type CalculateResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	service CalculatorService
}

func NewHandler(service CalculatorService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req CalculateRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			h.writeError(w, http.StatusRequestEntityTooLarge, "request payload exceeds size limit")
			return
		}
		h.writeError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	if req.Operation == "" {
		h.writeError(w, http.StatusBadRequest, "operation field is required")
		return
	}

	if req.A == nil {
		h.writeError(w, http.StatusBadRequest, "field 'a' is required")
		return
	}

	var bVal float64
	if !h.service.IsUnary(req.Operation) {
		if req.B == nil {
			h.writeError(w, http.StatusBadRequest, "field 'b' is required")
			return
		}
		bVal = *req.B
	}

	result, err := h.service.Calculate(req.Operation, *req.A, bVal)
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
