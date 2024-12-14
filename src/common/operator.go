package common

import "errors"

var opToPrecedence = map[string]int{
	"+": 10,
	"-": 10,
	"*": 20,
  "(": 0,
  ")": 0,
}

type Operator struct {
	Value      string
	Precedence int
}

func NewOperator(op string) (*Operator, error) {
	if precedence, ok := opToPrecedence[op]; ok {
		return &Operator{
			Value:      op,
			Precedence: precedence,
		}, nil
	} else {
		return nil, errors.New("Invalid operator")
	}
}

