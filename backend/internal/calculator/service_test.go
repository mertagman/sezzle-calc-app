package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestServiceOperations(t *testing.T) {
	svc := NewCalculatorService()

	t.Run("Add", func(t *testing.T) {
		res, err := svc.Add(5, 3)
		if err != nil || res != 8 {
			t.Fatalf("expected 8, got %v (err: %v)", res, err)
		}

		res, err = svc.Add(-5, -3)
		if err != nil || res != -8 {
			t.Fatalf("expected -8, got %v (err: %v)", res, err)
		}
	})

	t.Run("Subtract", func(t *testing.T) {
		res, err := svc.Subtract(10, 4)
		if err != nil || res != 6 {
			t.Fatalf("expected 6, got %v (err: %v)", res, err)
		}
	})

	t.Run("Multiply", func(t *testing.T) {
		res, err := svc.Multiply(6, 7)
		if err != nil || res != 42 {
			t.Fatalf("expected 42, got %v (err: %v)", res, err)
		}

		res, err = svc.Multiply(-0.0, 5)
		if err != nil || res != 0 {
			t.Fatalf("expected 0, got %v (err: %v)", res, err)
		}
	})

	t.Run("Divide", func(t *testing.T) {
		res, err := svc.Divide(20, 4)
		if err != nil || res != 5 {
			t.Fatalf("expected 5, got %v (err: %v)", res, err)
		}

		_, err = svc.Divide(10, 0)
		if !errors.Is(err, ErrDivisionByZero) {
			t.Fatalf("expected ErrDivisionByZero, got %v", err)
		}
	})

	t.Run("Power", func(t *testing.T) {
		res, err := svc.Power(2, 3)
		if err != nil || res != 8 {
			t.Fatalf("expected 8, got %v (err: %v)", res, err)
		}

		res, err = svc.Power(2, -3)
		if err != nil || math.Abs(res-0.125) > 1e-9 {
			t.Fatalf("expected 0.125, got %v (err: %v)", res, err)
		}

		_, err = svc.Power(0, -1)
		if !errors.Is(err, ErrUndefinedResult) {
			t.Fatalf("expected ErrUndefinedResult, got %v", err)
		}

		_, err = svc.Power(-4, 0.5)
		if !errors.Is(err, ErrUndefinedResult) {
			t.Fatalf("expected ErrUndefinedResult, got %v", err)
		}
	})

	t.Run("Sqrt", func(t *testing.T) {
		res, err := svc.Sqrt(16)
		if err != nil || res != 4 {
			t.Fatalf("expected 4, got %v (err: %v)", res, err)
		}

		res, err = svc.Sqrt(0)
		if err != nil || res != 0 {
			t.Fatalf("expected 0, got %v (err: %v)", res, err)
		}

		_, err = svc.Sqrt(-9)
		if !errors.Is(err, ErrNegativeSqrt) {
			t.Fatalf("expected ErrNegativeSqrt, got %v", err)
		}
	})

	t.Run("Percentage", func(t *testing.T) {
		res, err := svc.Percentage(50)
		if err != nil || res != 0.5 {
			t.Fatalf("expected 0.5, got %v (err: %v)", res, err)
		}

		res, err = svc.Percentage(0)
		if err != nil || res != 0 {
			t.Fatalf("expected 0, got %v (err: %v)", res, err)
		}
	})
}
