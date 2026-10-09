package calculator

import (
	"errors"
	"math"
	"strconv"
)

var (
	ErrDivisionByZero  = errors.New("division by zero is undefined")
	ErrNegativeSqrt    = errors.New("square root of negative number is undefined")
	ErrUndefinedResult = errors.New("operation resulted in undefined or infinite value")
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) normalize(val float64) (float64, error) {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, ErrUndefinedResult
	}
	rounded, err := strconv.ParseFloat(strconv.FormatFloat(val, 'g', 15, 64), 64)
	if err != nil {
		return 0, ErrUndefinedResult
	}
	if rounded == 0 {
		return 0, nil
	}
	return rounded, nil
}

func (s *Service) Add(a, b float64) (float64, error) {
	return s.normalize(a + b)
}

func (s *Service) Subtract(a, b float64) (float64, error) {
	return s.normalize(a - b)
}

func (s *Service) Multiply(a, b float64) (float64, error) {
	return s.normalize(a * b)
}

func (s *Service) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return s.normalize(a / b)
}

func (s *Service) Power(a, b float64) (float64, error) {
	return s.normalize(math.Pow(a, b))
}

func (s *Service) Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return s.normalize(math.Sqrt(a))
}

func (s *Service) Percentage(a float64) (float64, error) {
	return s.normalize(a / 100)
}
