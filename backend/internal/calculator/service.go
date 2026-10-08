package calculator

import (
	"errors"
	"math"
)

var (
	ErrDivisionByZero  = errors.New("division by zero is undefined")
	ErrNegativeSqrt    = errors.New("square root of negative number is undefined")
	ErrUndefinedResult = errors.New("operation resulted in undefined or infinite value")
)

type CalculatorService struct{}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (s *CalculatorService) normalize(val float64) (float64, error) {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, ErrUndefinedResult
	}
	if val == 0 {
		return 0, nil
	}
	return val, nil
}

func (s *CalculatorService) Add(a, b float64) (float64, error) {
	return s.normalize(a + b)
}

func (s *CalculatorService) Subtract(a, b float64) (float64, error) {
	return s.normalize(a - b)
}

func (s *CalculatorService) Multiply(a, b float64) (float64, error) {
	return s.normalize(a * b)
}

func (s *CalculatorService) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return s.normalize(a / b)
}

func (s *CalculatorService) Power(a, b float64) (float64, error) {
	return s.normalize(math.Pow(a, b))
}

func (s *CalculatorService) Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return s.normalize(math.Sqrt(a))
}

func (s *CalculatorService) Percentage(a float64) (float64, error) {
	return s.normalize(a / 100)
}
