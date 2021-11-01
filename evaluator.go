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

	// AsComparator attempts to convert to Comparator
	AsComparator() (Comparator, bool)
}

// Comparator is a special Evalutor whose evaluation expression is a comparison expression.
type Comparator interface {
	// Compare performs an comparison by giving a set of variables.
	Compare(Variables) (bool, error)
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
	case *ast.ParenExpr:
		x, err := parseExpr(str, expr.X)
		if err != nil {
			return nil, err
		}
		return &parenEvaluator{
			x: x,
		}, nil
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

func (e lockupVariableEvaluator) AsComparator() (Comparator, bool) {
	return nil, false
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
			f, _ := getLogicalFunc(token.LAND)
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

	if f, ok := getComputableFunc(expr.Op); ok {
		xLiteral, xok := xEvaluator.(*realNumericLiteralEvaluator)
		yLiteral, yok := yEvaluator.(*realNumericLiteralEvaluator)
		if xok && yok {
			value, err := f(xLiteral.value, yLiteral.value)
			if err == nil {
				value, ok := isRealNumber(value)
				if ok {
					return &realNumericLiteralEvaluator{
						value: value,
						str:   fmt.Sprintf("%f", value),
					}, nil
				}
			}
		}
		return &computableEvaluator{
			x:  xEvaluator,
			y:  yEvaluator,
			f:  f,
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
	ret, err := e.Compare(vars)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *comparativeEvaluator) Compare(vars Variables) (bool, error) {
	v1, err := e.x.Eval(vars)
	if err != nil {
		return false, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	v2, err := e.y.Eval(vars)
	if err != nil {
		return false, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	ret, err := e.f(v1, v2)
	if err != nil {
		return false, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	return ret, nil
}

func (e *comparativeEvaluator) AsComparator() (Comparator, bool) {
	return e, true
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

func (e *realNumericLiteralEvaluator) AsComparator() (Comparator, bool) {
	return nil, false
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

func (e *stringLiteralEvaluator) AsComparator() (Comparator, bool) {
	return nil, false
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
	ret, err := e.Compare(vars)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *logicalEvaluator) Compare(vars Variables) (bool, error) {
	v1, err := e.x.Eval(vars)
	if err != nil {
		return false, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	v2, err := e.y.Eval(vars)
	if err != nil {
		return false, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	b1, b2, ok := isBothBools(v1, v2)
	if !ok {
		return false, errors.New("is not both bool")
	}
	return e.f(b1, b2), nil
}

func (e *logicalEvaluator) AsComparator() (Comparator, bool) {
	return e, true
}

func (e *logicalEvaluator) String() string {
	return fmt.Sprintf("(%s) %s (%s)", e.x, e.op, e.y)
}

type computableEvaluator struct {
	x  Evaluator
	y  Evaluator
	f  computableFunc
	op string
}

func (e *computableEvaluator) Eval(vars Variables) (interface{}, error) {
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

func (e *computableEvaluator) AsComparator() (Comparator, bool) {
	return nil, false
}

func (e *computableEvaluator) String() string {
	return fmt.Sprintf("%s %s %s", e.x, e.op, e.y)
}

type parenEvaluator struct {
	x Evaluator
}

func (e *parenEvaluator) Eval(vars Variables) (interface{}, error) {
	return e.x.Eval(vars)
}

func (e *parenEvaluator) AsComparator() (Comparator, bool) {
	if x, ok := e.x.AsComparator(); ok {
		return x, true
	}
	return nil, false
}

func (e *parenEvaluator) String() string {
	return fmt.Sprintf("(%s)", e.x)
}
