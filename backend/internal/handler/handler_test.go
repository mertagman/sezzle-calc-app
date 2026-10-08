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
		handlerFunc    func(http.ResponseWriter, *http.Request)
		path           string
		method         string
		body           string
		expectedStatus int
		expectedResult *float64
		expectedError  string
	}{
		{
			name:           "successful add",
			handlerFunc:    h.Add,
			path:           "/api/v1/add",
			method:         http.MethodPost,
			body:           `{"a":10,"b":5}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 15.0; return &v }(),
		},
		{
			name:           "successful subtract",
			handlerFunc:    h.Subtract,
			path:           "/api/v1/subtract",
			method:         http.MethodPost,
			body:           `{"a":10,"b":4}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 6.0; return &v }(),
		},
		{
			name:           "successful multiply",
			handlerFunc:    h.Multiply,
			path:           "/api/v1/multiply",
			method:         http.MethodPost,
			body:           `{"a":6,"b":7}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 42.0; return &v }(),
		},
		{
			name:           "successful divide",
			handlerFunc:    h.Divide,
			path:           "/api/v1/divide",
			method:         http.MethodPost,
			body:           `{"a":20,"b":4}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 5.0; return &v }(),
		},
		{
			name:           "successful power",
			handlerFunc:    h.Power,
			path:           "/api/v1/power",
			method:         http.MethodPost,
			body:           `{"a":2,"b":3}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 8.0; return &v }(),
		},
		{
			name:           "successful sqrt",
			handlerFunc:    h.Sqrt,
			path:           "/api/v1/sqrt",
			method:         http.MethodPost,
			body:           `{"a":16}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 4.0; return &v }(),
		},
		{
			name:           "successful percentage",
			handlerFunc:    h.Percentage,
			path:           "/api/v1/percentage",
			method:         http.MethodPost,
			body:           `{"a":75}`,
			expectedStatus: http.StatusOK,
			expectedResult: func() *float64 { v := 0.75; return &v }(),
		},
		{
			name:           "empty request body",
			handlerFunc:    h.Add,
			path:           "/api/v1/add",
			method:         http.MethodPost,
			body:           "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "request body cannot be empty",
		},
		{
			name:           "trailing json tokens rejected",
			handlerFunc:    h.Add,
			path:           "/api/v1/add",
			method:         http.MethodPost,
			body:           `{"a":10,"b":5} {"extra":true}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "request body must only contain a single JSON object",
		},
		{
			name:           "invalid json payload syntax",
			handlerFunc:    h.Add,
			path:           "/api/v1/add",
			method:         http.MethodPost,
			body:           `{"a":`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid json payload",
		},
		{
			name:           "invalid json payload type",
			handlerFunc:    h.Add,
			path:           "/api/v1/add",
			method:         http.MethodPost,
			body:           `{"a": "not-a-number", "b": 2}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid json payload",
		},
		{
			name:           "unknown json field rejected",
			handlerFunc:    h.Add,
			path:           "/api/v1/add",
			method:         http.MethodPost,
			body:           `{"a":1,"b":2,"unexpected":true}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid json payload",
		},
		{
			name:           "missing field a binary",
			handlerFunc:    h.Add,
			path:           "/api/v1/add",
			method:         http.MethodPost,
			body:           `{"b":5}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "field 'a' is required",
		},
		{
			name:           "missing field b binary",
			handlerFunc:    h.Add,
			path:           "/api/v1/add",
			method:         http.MethodPost,
			body:           `{"a":5}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "field 'b' is required",
		},
		{
			name:           "missing field a unary",
			handlerFunc:    h.Sqrt,
			path:           "/api/v1/sqrt",
			method:         http.MethodPost,
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "field 'a' is required",
		},
		{
			name:           "division by zero unprocessable",
			handlerFunc:    h.Divide,
			path:           "/api/v1/divide",
			method:         http.MethodPost,
			body:           `{"a":10,"b":0}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "division by zero is undefined",
		},
		{
			name:           "negative square root unprocessable",
			handlerFunc:    h.Sqrt,
			path:           "/api/v1/sqrt",
			method:         http.MethodPost,
			body:           `{"a":-4}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "square root of negative number is undefined",
		},
		{
			name:           "undefined power calculation unprocessable",
			handlerFunc:    h.Power,
			path:           "/api/v1/power",
			method:         http.MethodPost,
			body:           `{"a":0,"b":-1}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "operation resulted in undefined or infinite value",
		},
		{
			name:           "payload exceeds size limit",
			handlerFunc:    h.Add,
			path:           "/api/v1/add",
			method:         http.MethodPost,
			body:           `{"a":1,"b":2,"extra":"` + strings.Repeat("x", 1048576+10) + `"}`,
			expectedStatus: http.StatusRequestEntityTooLarge,
			expectedError:  "request payload exceeds size limit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			tc.handlerFunc(rec, req)

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
