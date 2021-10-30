package evaluator

import (
	"errors"
	"fmt"
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
	return parseExpr(expr, astExpr)
}

func parseExpr(str string, expr ast.Expr) (Evaluator, error) {
	switch expr := expr.(type) {
	case *ast.Ident:
		return lockupVariableEvaluator(expr.Name), nil
	case *ast.BinaryExpr:
		return parseBinaryExpr(str, expr)
	default:
		return nil, fmt.Errorf("can not parse `%s` ast type `%T` not implemented", str, expr)
	}
}

type lockupVariableEvaluator string

func (e lockupVariableEvaluator) Eval(vars Variables) (interface{}, error) {
	if v, ok := vars[string(e)]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("Variable %s is not givend", e)
}

func parseBinaryExpr(str string, expr *ast.BinaryExpr) (Evaluator, error) {
	return nil, errors.New("not implemented yet")
}
