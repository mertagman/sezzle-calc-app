package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestServiceBinaryOperations(t *testing.T) {
	svc := New()

	tests := []struct {
		name      string
		op        func(float64, float64) (float64, error)
		a         float64
		b         float64
		want      float64
		wantErr   error
		tolerance float64
	}{
		{name: "Add positive", op: svc.Add, a: 5, b: 3, want: 8},
		{name: "Add negative", op: svc.Add, a: -5, b: -3, want: -8},
		{name: "Add floating precision", op: svc.Add, a: 0.1, b: 0.2, want: 0.3},
		{name: "Subtract positive", op: svc.Subtract, a: 10, b: 4, want: 6},
		{name: "Multiply positive", op: svc.Multiply, a: 6, b: 7, want: 42},
		{name: "Multiply negative zero", op: svc.Multiply, a: -0.0, b: 5, want: 0},
		{name: "Multiply precision", op: svc.Multiply, a: 1.1, b: 3, want: 3.3},
		{name: "Multiply overflow", op: svc.Multiply, a: math.MaxFloat64, b: 2, wantErr: ErrUndefinedResult},
		{name: "Divide valid", op: svc.Divide, a: 20, b: 4, want: 5},
		{name: "Divide by zero", op: svc.Divide, a: 10, b: 0, wantErr: ErrDivisionByZero},
		{name: "Power valid", op: svc.Power, a: 2, b: 3, want: 8},
		{name: "Power negative exp", op: svc.Power, a: 2, b: -3, want: 0.125, tolerance: 1e-9},
		{name: "Power zero negative exp", op: svc.Power, a: 0, b: -1, wantErr: ErrUndefinedResult},
		{name: "Power negative fractional", op: svc.Power, a: -4, b: 0.5, wantErr: ErrUndefinedResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.op(tt.a, tt.b)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.tolerance > 0 {
				if math.Abs(got-tt.want) > tt.tolerance {
					t.Errorf("expected %v, got %v", tt.want, got)
				}
				return
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestServiceUnaryOperations(t *testing.T) {
	svc := New()

	tests := []struct {
		name    string
		op      func(float64) (float64, error)
		a       float64
		want    float64
		wantErr error
	}{
		{name: "Sqrt positive", op: svc.Sqrt, a: 16, want: 4},
		{name: "Sqrt zero", op: svc.Sqrt, a: 0, want: 0},
		{name: "Sqrt negative", op: svc.Sqrt, a: -9, wantErr: ErrNegativeSqrt},
		{name: "Percentage positive", op: svc.Percentage, a: 50, want: 0.5},
		{name: "Percentage zero", op: svc.Percentage, a: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.op(tt.a)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
