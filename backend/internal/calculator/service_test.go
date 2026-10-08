package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	svc := NewCalculatorService()

	tests := []struct {
		name      string
		operation string
		a         float64
		b         float64
		expected  float64
		err       error
	}{
		{name: "addition positive", operation: "add", a: 5, b: 3, expected: 8, err: nil},
		{name: "addition negative", operation: "add", a: -5, b: -3, expected: -8, err: nil},
		{name: "subtraction", operation: "subtract", a: 10, b: 4, expected: 6, err: nil},
		{name: "multiplication", operation: "multiply", a: 6, b: 7, expected: 42, err: nil},
		{name: "division valid", operation: "divide", a: 20, b: 4, expected: 5, err: nil},
		{name: "division by zero", operation: "divide", a: 10, b: 0, expected: 0, err: ErrDivisionByZero},
		{name: "division by zero negative", operation: "divide", a: -10, b: 0, expected: 0, err: ErrDivisionByZero},
		{name: "power valid", operation: "power", a: 2, b: 3, expected: 8, err: nil},
		{name: "power negative exponent", operation: "power", a: 2, b: -3, expected: 0.125, err: nil},
		{name: "power zero base negative exponent", operation: "power", a: 0, b: -1, expected: 0, err: ErrUndefinedResult},
		{name: "power invalid negative fractional", operation: "power", a: -4, b: 0.5, expected: 0, err: ErrUndefinedResult},
		{name: "sqrt valid", operation: "sqrt", a: 16, b: 0, expected: 4, err: nil},
		{name: "sqrt zero", operation: "sqrt", a: 0, b: 0, expected: 0, err: nil},
		{name: "sqrt negative", operation: "sqrt", a: -9, b: 0, expected: 0, err: ErrNegativeSqrt},
		{name: "percentage unary", operation: "percentage", a: 50, b: 0, expected: 0.5, err: nil},
		{name: "percentage zero", operation: "percentage", a: 0, b: 0, expected: 0, err: nil},
		{name: "negative zero normalization", operation: "multiply", a: -0.0, b: 5, expected: 0, err: nil},
		{name: "unsupported operation", operation: "modulo", a: 10, b: 2, expected: 0, err: ErrInvalidOperator},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.Calculate(tc.operation, tc.a, tc.b)

			if tc.err != nil {
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error %v, got %v", tc.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(result-tc.expected) > 1e-9 {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestIsUnary(t *testing.T) {
	svc := NewCalculatorService()

	if !svc.IsUnary("sqrt") {
		t.Errorf("expected sqrt to be unary")
	}
	if !svc.IsUnary("percentage") {
		t.Errorf("expected percentage to be unary")
	}
	if svc.IsUnary("add") {
		t.Errorf("expected add not to be unary")
	}
	if svc.IsUnary("subtract") {
		t.Errorf("expected subtract not to be unary")
	}
}
