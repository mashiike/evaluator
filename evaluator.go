package evaluator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
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
	case *ast.BasicLit:
		return parseBasicLit(str, expr)
	default:
		return nil, fmt.Errorf("can not parse `%s` ast type `%T` not implemented", str, expr)
	}
}

type lockupVariableEvaluator string

func (e lockupVariableEvaluator) Eval(vars Variables) (interface{}, error) {
	if v, ok := vars[string(e)]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("variable %s is not givend", e)
}

func getSubExpr(str string, expr ast.Expr) string {
	return extructStrPos(str, expr.Pos(), expr.End())
}

func extructStrPos(str string, pos, end token.Pos) string {
	return str[pos-1 : end-1]
}

func parseBinaryExpr(str string, expr *ast.BinaryExpr) (Evaluator, error) {
	xStr := getSubExpr(str, expr.X)
	xEvaluator, err := parseExpr(str, expr.X)
	if err != nil {
		return nil, fmt.Errorf("parse BinaryExpr.X `%s` %w", xStr, err)
	}
	yStr := getSubExpr(str, expr.Y)
	yEvaluator, err := parseExpr(str, expr.Y)
	if err != nil {
		return nil, fmt.Errorf("parse BinaryExpr.Y `%s` %w", yStr, err)
	}
	op := strings.TrimSpace(extructStrPos(str, expr.OpPos, expr.Y.Pos()))
	if cf, ok := getComparativeFunc(expr.Op); ok {
		if xLogical, ok := xEvaluator.(*comparativeEvaluator); ok {
			f, _ := getLogicalFunc(token.AND)
			return &logicalEvaluator{
				x: xEvaluator,
				y: &comparativeEvaluator{
					x:  xLogical.y,
					y:  yEvaluator,
					f:  cf,
					op: op,
				},
				f:  f,
				op: "&&",
			}, nil
		}
		return &comparativeEvaluator{
			x:  xEvaluator,
			y:  yEvaluator,
			f:  cf,
			op: op,
		}, nil
	}

	if lf, ok := getLogicalFunc(expr.Op); ok {
		return &logicalEvaluator{
			x:  xEvaluator,
			y:  yEvaluator,
			f:  lf,
			op: op,
		}, err
	}

	return nil, fmt.Errorf("parse `%s` invalid operator", op)
}

type comparativeEvaluator struct {
	x  Evaluator
	y  Evaluator
	f  comparativeFunc
	op string
}

func (e *comparativeEvaluator) Eval(vars Variables) (interface{}, error) {
	v1, err := e.x.Eval(vars)
	if err != nil {
		return nil, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	v2, err := e.y.Eval(vars)
	if err != nil {
		return nil, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	ret, err := e.f(v1, v2)
	if err != nil {
		return nil, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	return ret, nil
}

func (e *comparativeEvaluator) String() string {
	return fmt.Sprintf("%s %s %s", e.x, e.op, e.y)
}

func parseBasicLit(str string, expr *ast.BasicLit) (Evaluator, error) {

	switch expr.Kind {
	case token.INT, token.FLOAT:
		v, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			return nil, err
		}
		return &realNumericLiteralEvaluator{
			str:   strings.TrimSpace(getSubExpr(str, expr)),
			value: v,
		}, nil
	case token.STRING:
		return &stringLiteralEvaluator{
			str: expr.Value,
		}, nil
	default:
		return nil, fmt.Errorf("unknown literal `%s`", expr.Kind)
	}
}

type realNumericLiteralEvaluator struct {
	value float64
	str   string
}

func (e *realNumericLiteralEvaluator) Eval(vars Variables) (interface{}, error) {
	return e.value, nil
}
func (e *realNumericLiteralEvaluator) String() string {
	return e.str
}

type stringLiteralEvaluator struct {
	str string
}

func (e *stringLiteralEvaluator) Eval(vars Variables) (interface{}, error) {
	return e.str, nil
}
func (e *stringLiteralEvaluator) String() string {
	return e.str
}

type logicalEvaluator struct {
	x  Evaluator
	y  Evaluator
	f  logicalFunc
	op string
}

func (e *logicalEvaluator) Eval(vars Variables) (interface{}, error) {
	v1, err := e.x.Eval(vars)
	if err != nil {
		return nil, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	v2, err := e.y.Eval(vars)
	if err != nil {
		return nil, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	b1, b2, ok := isBothBools(v1, v2)
	if !ok {
		return nil, errors.New("is not both bool")
	}
	return e.f(b1, b2), nil
}

func (e *logicalEvaluator) String() string {
	return fmt.Sprintf("(%s) %s (%s)", e.x, e.op, e.y)
}
