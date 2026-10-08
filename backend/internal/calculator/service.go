package calculator

import (
	"errors"
	"math"
)

var (
	ErrDivisionByZero  = errors.New("division by zero is undefined")
	ErrNegativeSqrt    = errors.New("square root of negative number is undefined")
	ErrInvalidOperator = errors.New("unsupported operation")
	ErrUndefinedResult = errors.New("operation resulted in undefined or infinite value")
)

type CalculatorService struct{}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (s *CalculatorService) IsUnary(operation string) bool {
	return operation == "sqrt" || operation == "percentage"
}

func (s *CalculatorService) Calculate(operation string, a float64, b float64) (float64, error) {
	var result float64

	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		result = a / b
	case "power":
		result = math.Pow(a, b)
	case "sqrt":
		if a < 0 {
			return 0, ErrNegativeSqrt
		}
		result = math.Sqrt(a)
	case "percentage":
		result = a / 100
	default:
		return 0, ErrInvalidOperator
	}

	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrUndefinedResult
	}

	if result == 0 {
		result = 0
	}

	return result, nil
}
