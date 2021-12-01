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

	//Strict sets the expression evaluation to run in strict mode. The default is false.
	//For example, the behavior changes when the variable is not set.
	//  Referenced as nil if Strict is set to false.
	//  If set to true, the expression evaluation will be Error.
	Strict(bool)

	// AsComparator attempts to convert to Comparator
	AsComparator() (Comparator, bool)

	fmt.Stringer
}

// Comparator is a special Evaluator whose evaluation expression is a comparison expression.
type Comparator interface {
	// Compare performs an comparison by giving a set of variables.
	Compare(Variables) (bool, error)

	fmt.Stringer
}

// Variables are a group of variables given to the evaluator
type Variables map[string]interface{}

// New parses the expression to create an evaluator
func New(expr string) (Evaluator, error) {
	expr = prepare(expr)
	astExpr, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}
	return parseExpr(expr, astExpr)
}

func prepare(expr string) string {
	//replace if( => __if(
	return strings.ReplaceAll(expr, "if(", "__if(")
}

func parseExpr(str string, expr ast.Expr) (Evaluator, error) {
	switch expr := expr.(type) {
	case *ast.Ident:
		return parseIdent(str, expr)
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
	case *ast.CallExpr:
		return parseCallExpr(str, expr)
	case *ast.UnaryExpr:
		return parseUnaryExpr(str, expr)
	default:
		return nil, fmt.Errorf("can not parse `%s` ast type `%T` not implemented", str, expr)
	}
}

func parseIdent(str string, expr *ast.Ident) (Evaluator, error) {
	if expr.Name == "nil" {
		return nilEvaluator{}, nil
	}
	return newLockupVariableEvaluator(expr.Name), nil
}

type nilEvaluator struct{}

func (e nilEvaluator) Eval(vars Variables) (interface{}, error) {
	return nil, nil
}

func (e nilEvaluator) Strict(bool) {}

func (e nilEvaluator) AsComparator() (Comparator, bool) {
	return nil, false
}

func (e nilEvaluator) String() string {
	return "nil"
}

type lockupVariableEvaluator struct {
	strict bool
	name   string
}

func newLockupVariableEvaluator(name string) *lockupVariableEvaluator {
	return &lockupVariableEvaluator{
		name:   name,
		strict: false,
	}
}

func (e *lockupVariableEvaluator) Eval(vars Variables) (interface{}, error) {
	if v, ok := vars[string(e.name)]; ok {
		return v, nil
	}
	if e.strict {
		return nil, fmt.Errorf("%s %w", e, ErrVariableNotFound)
	}
	return nil, nil
}

func (e *lockupVariableEvaluator) Strict(v bool) {
	e.strict = v
}

func (e *lockupVariableEvaluator) AsComparator() (Comparator, bool) {
	return nil, false
}

func (e *lockupVariableEvaluator) String() string {
	return e.name
}

func getSubExpr(str string, expr ast.Expr) string {
	return extractStrPos(str, expr.Pos(), expr.End())
}

func extractStrPos(str string, pos, end token.Pos) string {
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
	op := strings.TrimSpace(extractStrPos(str, expr.OpPos, expr.Y.Pos()))
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

func (e *comparativeEvaluator) Strict(v bool) {
	e.x.Strict(v)
	e.y.Strict(v)
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
		return newRealNumericLiteralEvaluator(v, strings.TrimSpace(getSubExpr(str, expr))), nil
	case token.STRING:
		return newStringLiteralEvaluator(strings.Trim(expr.Value, "`\"")), nil
	case token.CHAR:
		return newStringLiteralEvaluator(strings.Trim(expr.Value, "'")), nil
	default:
		return nil, fmt.Errorf("unknown literal `%s`", expr.Kind)
	}
}

type realNumericLiteralEvaluator struct {
	value float64
	str   string
}

func newRealNumericLiteralEvaluator(value float64, str string) *realNumericLiteralEvaluator {
	return &realNumericLiteralEvaluator{
		str:   str,
		value: value,
	}
}

func (e *realNumericLiteralEvaluator) Eval(vars Variables) (interface{}, error) {
	return e.value, nil
}

func (e *realNumericLiteralEvaluator) Strict(bool) {}

func (e *realNumericLiteralEvaluator) AsComparator() (Comparator, bool) {
	return nil, false
}

func (e *realNumericLiteralEvaluator) String() string {
	return e.str
}

type stringLiteralEvaluator struct {
	str string
}

func newStringLiteralEvaluator(str string) *stringLiteralEvaluator {
	return &stringLiteralEvaluator{
		str: str,
	}
}

func (e *stringLiteralEvaluator) Eval(vars Variables) (interface{}, error) {
	return e.str, nil
}

func (e *stringLiteralEvaluator) Strict(bool) {}

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

func (e *logicalEvaluator) Strict(v bool) {
	e.x.Strict(v)
	e.y.Strict(v)
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

func (e *computableEvaluator) Strict(v bool) {
	e.x.Strict(v)
	e.y.Strict(v)
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

func (e *parenEvaluator) Strict(v bool) {
	e.x.Strict(v)
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

func parseCallExpr(str string, expr *ast.CallExpr) (Evaluator, error) {
	argEvaluators := make([]Evaluator, 0, len(expr.Args))
	for i, arg := range expr.Args {
		argStr := getSubExpr(str, arg)
		argEvaluator, err := parseExpr(str, arg)
		if err != nil {
			return nil, fmt.Errorf("parse CallExpr.Args[%d] `%s` %w", i, argStr, err)
		}
		argEvaluators = append(argEvaluators, argEvaluator)
	}
	var funcName string
	funStr := getSubExpr(str, expr.Fun)
	switch fun := expr.Fun.(type) {
	case *ast.Ident:
		funcName = fun.Name
	default:
		return nil, fmt.Errorf("parse CallExpr.Fun `%s` unexpected type %T", funStr, fun)
	}
	f, err := getCallFunc(funcName, argEvaluators)
	if err != nil {
		return nil, err
	}
	return &callEvaluator{
		args:     argEvaluators,
		f:        f,
		funcName: funcName,
	}, nil
}

type callEvaluator struct {
	args     []Evaluator
	f        callFunc
	funcName string
}

func (e *callEvaluator) Eval(vars Variables) (interface{}, error) {
	args := make([]interface{}, 0, len(e.args))
	for i, a := range e.args {
		arg, err := a.Eval(vars)
		if err != nil {
			return nil, fmt.Errorf("Eval(`%s`) Args[%d] %w", e, i, err)
		}
		args = append(args, arg)
	}
	return e.f(args...)
}

func (e *callEvaluator) Strict(v bool) {
	for _, arg := range e.args {
		arg.Strict(v)
	}
}

func (e *callEvaluator) AsComparator() (Comparator, bool) {
	return nil, false
}

func (e *callEvaluator) String() string {
	var builder strings.Builder
	builder.WriteString(e.funcName)
	builder.WriteRune('(')
	for i, arg := range e.args {
		builder.WriteString(arg.String())
		if i+1 < len(e.args) {
			builder.WriteString(", ")
		}
	}
	builder.WriteRune(')')
	return builder.String()
}

func parseUnaryExpr(str string, expr *ast.UnaryExpr) (Evaluator, error) {
	xStr := getSubExpr(str, expr.X)
	xEvaluator, err := parseExpr(str, expr.X)
	if err != nil {
		return nil, fmt.Errorf("parse BinaryExpr.X `%s` %w", xStr, err)
	}
	op := strings.TrimSpace(extractStrPos(str, expr.OpPos, expr.End()))
	if f, ok := getUnaryFunc(expr.Op); ok {
		return &unaryEvaluator{
			x:  xEvaluator,
			f:  f,
			op: op,
		}, nil
	}
	return nil, fmt.Errorf("parse `%s` invalid operator", op)
}

type unaryEvaluator struct {
	x  Evaluator
	f  unaryFunc
	op string
}

func (e *unaryEvaluator) Eval(vars Variables) (interface{}, error) {
	v, err := e.x.Eval(vars)
	if err != nil {
		return false, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	ret, err := e.f(v)
	if err != nil {
		return false, fmt.Errorf("Eval(`%s`) %w", e, err)
	}
	return ret, nil
}

func (e *unaryEvaluator) Strict(v bool) {
	e.x.Strict(true)
}

func (e *unaryEvaluator) AsComparator() (Comparator, bool) {
	return nil, false
}

func (e *unaryEvaluator) String() string {
	return fmt.Sprintf("%s(%s)", e.op, e.x)
}
