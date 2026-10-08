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
	h := NewHandler(svc)

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedResult *float64
		expectedError  string
	}{
		{
			name:           "successful binary calculation",
			method:         http.MethodPost,
			body:           `{"operation":"add","a":10,"b":5}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 15.0; return &v }(),
		},
		{
			name:           "successful unary sqrt calculation without b",
			method:         http.MethodPost,
			body:           `{"operation":"sqrt","a":16}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 4.0; return &v }(),
		},
		{
			name:           "successful unary percentage calculation without b",
			method:         http.MethodPost,
			body:           `{"operation":"percentage","a":75}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 0.75; return &v }(),
		},
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "method not allowed",
		},
		{
			name:           "invalid json payload syntax",
			method:         http.MethodPost,
			body:           `{"operation": "add", "a":`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid json payload",
		},
		{
			name:           "invalid json payload type",
			method:         http.MethodPost,
			body:           `{"operation": "add", "a": "not-a-number", "b": 2}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid json payload",
		},
		{
			name:           "unknown json field rejected",
			method:         http.MethodPost,
			body:           `{"operation":"add","a":1,"b":2,"unexpected":true}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid json payload",
		},
		{
			name:           "missing operation field",
			method:         http.MethodPost,
			body:           `{"a":10,"b":5}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "operation field is required",
		},
		{
			name:           "missing field a",
			method:         http.MethodPost,
			body:           `{"operation":"add","b":5}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "field 'a' is required",
		},
		{
			name:           "missing field b for binary operation",
			method:         http.MethodPost,
			body:           `{"operation":"multiply","a":5}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "field 'b' is required",
		},
		{
			name:           "unsupported operation",
			method:         http.MethodPost,
			body:           `{"operation":"invalid_op","a":10,"b":2}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "unsupported operation",
		},
		{
			name:           "division by zero unprocessable",
			method:         http.MethodPost,
			body:           `{"operation":"divide","a":10,"b":0}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "division by zero is undefined",
		},
		{
			name:           "negative square root unprocessable",
			method:         http.MethodPost,
			body:           `{"operation":"sqrt","a":-4}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "square root of negative number is undefined",
		},
		{
			name:           "undefined power calculation unprocessable",
			method:         http.MethodPost,
			body:           `{"operation":"power","a":0,"b":-1}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "operation resulted in undefined or infinite value",
		},
		{
			name:           "payload exceeds size limit",
			method:         http.MethodPost,
			body:           `{"operation":"add","a":1,"b":2,"extra":"` + strings.Repeat("x", 1048576+10) + `"}`,
			expectedStatus: http.StatusRequestEntityTooLarge,
			expectedError:  "request payload exceeds size limit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/v1/calculate", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")

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
