package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"calculator/internal/calculator"
)

func TestHandler(t *testing.T) {
	svc := calculator.NewCalculatorService()
	allowedOrigin := "http://localhost:5173"
	h := NewHandler(svc, allowedOrigin)

	tests := []struct {
		name           string
		method         string
		origin         string
		body           string
		expectedStatus int
		expectedResult *float64
		expectedError  string
	}{
		{
			name:           "successful calculation",
			method:         http.MethodPost,
			origin:         "http://localhost:5173",
			body:           `{"operation":"add","a":10,"b":5}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 15.0; return &v }(),
		},
		{
			name:           "cors preflight allowed",
			method:         http.MethodOptions,
			origin:         "http://localhost:5173",
			body:           "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "cors forbidden origin",
			method:         http.MethodPost,
			origin:         "http://malicious-site.com",
			body:           `{"operation":"add","a":1,"b":1}`,
			expectedStatus: http.StatusForbidden,
			expectedError:  "cors origin not allowed",
		},
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			origin:         "http://localhost:5173",
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "method not allowed",
		},
		{
			name:           "invalid json payload",
			method:         http.MethodPost,
			origin:         "http://localhost:5173",
			body:           `{"operation": "add", "a": "not-a-number"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid json payload",
		},
		{
			name:           "unknown json field rejected",
			method:         http.MethodPost,
			origin:         "http://localhost:5173",
			body:           `{"operation":"add","a":1,"b":2,"unexpected":true}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid json payload",
		},
		{
			name:           "missing operation field",
			method:         http.MethodPost,
			origin:         "http://localhost:5173",
			body:           `{"a":10,"b":5}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "operation field is required",
		},
		{
			name:           "division by zero unprocessable",
			method:         http.MethodPost,
			origin:         "http://localhost:5173",
			body:           `{"operation":"divide","a":10,"b":0}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "division by zero is undefined",
		},
		{
			name:           "negative square root unprocessable",
			method:         http.MethodPost,
			origin:         "http://localhost:5173",
			body:           `{"operation":"sqrt","a":-4,"b":0}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "square root of negative number is undefined",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/v1/calculate", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d. Body: %s", tc.expectedStatus, rec.Code, rec.Body.String())
			}

			if tc.expectedResult != nil {
				var res CalculateResponse
				if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if res.Result != *tc.expectedResult {
					t.Errorf("expected result %v, got %v", *tc.expectedResult, res.Result)
				}
			}

			if tc.expectedError != "" {
				var errRes ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&errRes); err != nil {
					t.Fatalf("failed to decode error body: %v", err)
				}
				if !strings.Contains(errRes.Error, tc.expectedError) {
					t.Errorf("expected error containing %q, got %q", tc.expectedError, errRes.Error)
				}
			}
		})
	}
}
