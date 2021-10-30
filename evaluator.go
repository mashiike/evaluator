package evaluator

import (
	"errors"
	"go/ast"
	"go/parser"
)

// Evaluator is a variable evaluator created based on one expression
type Evaluator interface {
	// Eval performs an evaluation by giving a set of variables.
	Eval(Variables) (interface{}, error)
}

// Variables are a group of variables given to the evaluator
type Variables map[string]interface{}

// New parses the expression to create an evaluator
func New(expr string) (Evaluator, error) {
	astExpr, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}
	return parseExpr(astExpr)
}

func parseExpr(expr ast.Expr) (Evaluator, error) {
	return nil, errors.New("not implemented yet")
}
