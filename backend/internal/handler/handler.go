package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"calculator/internal/calculator"
)

type CalculatorService interface {
	Add(a, b float64) (float64, error)
	Subtract(a, b float64) (float64, error)
	Multiply(a, b float64) (float64, error)
	Divide(a, b float64) (float64, error)
	Power(a, b float64) (float64, error)
	Sqrt(a float64) (float64, error)
	Percentage(a float64) (float64, error)
}

type BinaryRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

type UnaryRequest struct {
	A *float64 `json:"a"`
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

func decodeJSON[T any](h *Handler, w http.ResponseWriter, r *http.Request) (*T, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req T
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			h.writeError(w, http.StatusRequestEntityTooLarge, "request payload exceeds size limit")
			return nil, false
		}
		if errors.Is(err, io.EOF) {
			h.writeError(w, http.StatusBadRequest, "request body cannot be empty")
			return nil, false
		}
		h.writeError(w, http.StatusBadRequest, "invalid json payload")
		return nil, false
	}

	if decoder.Decode(&struct{}{}) != io.EOF {
		h.writeError(w, http.StatusBadRequest, "request body must only contain a single JSON object")
		return nil, false
	}

	return &req, true
}

func (h *Handler) decodeBinary(w http.ResponseWriter, r *http.Request) (*BinaryRequest, bool) {
	req, ok := decodeJSON[BinaryRequest](h, w, r)
	if !ok {
		return nil, false
	}

	if req.A == nil {
		h.writeError(w, http.StatusBadRequest, "field 'a' is required")
		return nil, false
	}

	if req.B == nil {
		h.writeError(w, http.StatusBadRequest, "field 'b' is required")
		return nil, false
	}

	return req, true
}

func (h *Handler) decodeUnary(w http.ResponseWriter, r *http.Request) (*UnaryRequest, bool) {
	req, ok := decodeJSON[UnaryRequest](h, w, r)
	if !ok {
		return nil, false
	}

	if req.A == nil {
		h.writeError(w, http.StatusBadRequest, "field 'a' is required")
		return nil, false
	}

	return req, true
}

func (h *Handler) handleResult(w http.ResponseWriter, result float64, err error) {
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

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeBinary(w, r)
	if !ok {
		return
	}
	result, err := h.service.Add(*req.A, *req.B)
	h.handleResult(w, result, err)
}

func (h *Handler) Subtract(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeBinary(w, r)
	if !ok {
		return
	}
	result, err := h.service.Subtract(*req.A, *req.B)
	h.handleResult(w, result, err)
}

func (h *Handler) Multiply(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeBinary(w, r)
	if !ok {
		return
	}
	result, err := h.service.Multiply(*req.A, *req.B)
	h.handleResult(w, result, err)
}

func (h *Handler) Divide(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeBinary(w, r)
	if !ok {
		return
	}
	result, err := h.service.Divide(*req.A, *req.B)
	h.handleResult(w, result, err)
}

func (h *Handler) Power(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeBinary(w, r)
	if !ok {
		return
	}
	result, err := h.service.Power(*req.A, *req.B)
	h.handleResult(w, result, err)
}

func (h *Handler) Sqrt(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeUnary(w, r)
	if !ok {
		return
	}
	result, err := h.service.Sqrt(*req.A)
	h.handleResult(w, result, err)
}

func (h *Handler) Percentage(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeUnary(w, r)
	if !ok {
		return
	}
	result, err := h.service.Percentage(*req.A)
	h.handleResult(w, result, err)
}

func (h *Handler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
