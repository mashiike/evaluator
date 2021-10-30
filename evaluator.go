package evaluator

import "errors"

// Evaluator is a variable evaluator created based on one expression
type Evaluator interface {
	// Eval performs an evaluation by giving a set of variables. The result of the evaluation is a boolean value
	Eval(Variables) (bool, error)
}

// Variables are a group of variables given to the evaluator
type Variables map[string]interface{}

// New parses the expression to create an evaluator
func New(expr string) (Evaluator, error) {
	return nil, errors.New("not implemented yet")
}
