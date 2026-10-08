package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"calculator/internal/calculator"
	"calculator/internal/handler"
)

const testOrigin = "http://localhost:5173"

func TestRouter(t *testing.T) {
	router := NewRouter(handler.NewHandler(calculator.NewCalculatorService()), testOrigin)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		origin         string
		expectedStatus int
		expectedCORS   bool
	}{
		{name: "add route", method: http.MethodPost, path: "/api/v1/add", body: `{"a":1,"b":2}`, expectedStatus: http.StatusOK},
		{name: "subtract route", method: http.MethodPost, path: "/api/v1/subtract", body: `{"a":1,"b":2}`, expectedStatus: http.StatusOK},
		{name: "multiply route", method: http.MethodPost, path: "/api/v1/multiply", body: `{"a":1,"b":2}`, expectedStatus: http.StatusOK},
		{name: "divide route", method: http.MethodPost, path: "/api/v1/divide", body: `{"a":1,"b":2}`, expectedStatus: http.StatusOK},
		{name: "power route", method: http.MethodPost, path: "/api/v1/power", body: `{"a":1,"b":2}`, expectedStatus: http.StatusOK},
		{name: "sqrt route", method: http.MethodPost, path: "/api/v1/sqrt", body: `{"a":4}`, expectedStatus: http.StatusOK},
		{name: "percentage route", method: http.MethodPost, path: "/api/v1/percentage", body: `{"a":50}`, expectedStatus: http.StatusOK},
		{name: "health route", method: http.MethodGet, path: "/health", expectedStatus: http.StatusOK},
		{name: "allowed origin gets cors headers", method: http.MethodPost, path: "/api/v1/add", body: `{"a":1,"b":2}`, origin: testOrigin, expectedStatus: http.StatusOK, expectedCORS: true},
		{name: "preflight from allowed origin", method: http.MethodOptions, path: "/api/v1/add", origin: testOrigin, expectedStatus: http.StatusNoContent, expectedCORS: true},
		{name: "preflight from disallowed origin", method: http.MethodOptions, path: "/api/v1/add", origin: "http://evil.example", expectedStatus: http.StatusForbidden},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d. Body: %s", tc.expectedStatus, rec.Code, rec.Body.String())
			}

			allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
			if tc.expectedCORS && allowOrigin != tc.origin {
				t.Errorf("expected Access-Control-Allow-Origin %q, got %q", tc.origin, allowOrigin)
			}
			if !tc.expectedCORS && allowOrigin != "" {
				t.Errorf("expected no Access-Control-Allow-Origin header, got %q", allowOrigin)
			}
		})
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	panicking := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	rec := httptest.NewRecorder()
	recoveryMiddleware(panicking).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if body := rec.Body.String(); body != `{"error":"internal server error"}` {
		t.Errorf("unexpected body: %s", body)
	}
}
