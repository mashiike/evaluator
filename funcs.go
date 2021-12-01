package evaluator

import "fmt"

type callFunc func(...interface{}) (interface{}, error)

func getCallFunc(funcName string, argEvaluators []Evaluator) (callFunc, error) {
	switch funcName {
	case "rate": //rate(number, number)
		if len(argEvaluators) != 2 {
			return nil, newNumOfArgumentsMismatchError(funcName, 2, len(argEvaluators))
		}
		return rateCallFunc, nil
	case "coalesce": //coalesce(any, any, ...)
		return coalesceCallFunc, nil
	case "as_numeric": // as_numeric(any)
		if len(argEvaluators) != 1 {
			return nil, newNumOfArgumentsMismatchError(funcName, 1, len(argEvaluators))
		}
		return asNumericCallFunc, nil
	case "as_string": // as_string(any)
		if len(argEvaluators) != 1 {
			return nil, newNumOfArgumentsMismatchError(funcName, 1, len(argEvaluators))
		}
		return asStringCallFunc, nil
	case "if": // if(bool, any, any)
		if len(argEvaluators) != 3 {
			return nil, newNumOfArgumentsMismatchError(funcName, 3, len(argEvaluators))
		}
		return ifCallFunc, nil
	default:
		return nil, fmt.Errorf("%s() func is not found", funcName)
	}
}

func rateCallFunc(args ...interface{}) (interface{}, error) {
	n1, n2, ok := isBothRealNumbers(args[0], args[1])
	if !ok {
		return nil, fmt.Errorf("rate(v1[%v]::%T,v2[%v]::%T) can not eval", args[0], args[0], args[1], args[1])
	}
	if n2 == 0.0 {
		return nil, nil
	}
	return n1 / n2, nil
}

func coalesceCallFunc(args ...interface{}) (interface{}, error) {
	for _, arg := range args {
		if arg != nil {
			return arg, nil
		}
	}
	return nil, nil
}

func asStringCallFunc(args ...interface{}) (interface{}, error) {
	if v, ok := asString(args[0]); ok {
		return v, nil
	}
	return nil, nil
}

func asNumericCallFunc(args ...interface{}) (interface{}, error) {
	if v, ok := asNumber(args[0]); ok {
		return v, nil
	}
	return nil, nil
}

func ifCallFunc(args ...interface{}) (interface{}, error) {
	cond, ok := asBool(args[0])
	if !ok {
		return nil, fmt.Errorf("if(v1[%v]::%T,v2[%v]::%T,v3[%v]::%T) can not eval", args[0], args[0], args[1], args[1], args[2], args[2])
	}
	if cond {
		return args[1], nil
	}
	return args[2], nil
}
