package evaluator

import "fmt"

type callFunc func(...interface{}) (interface{}, error)

func getCallFunc(funcName string, argEvaluators []Evaluator) (callFunc, error) {
	switch funcName {
	case "rate":
		if len(argEvaluators) != 2 {
			return nil, fmt.Errorf("rate() func is expected 2 args, but given %d args", len(argEvaluators))
		}
		return rateCallFunc, nil
	case "coalesce":
		return coalesceCallFunc, nil
	case "as_numeric":
		if len(argEvaluators) != 1 {
			return nil, fmt.Errorf("as_numeric() func is expected 1 arg, but given %d args", len(argEvaluators))
		}
		return asNumericCallFunc, nil
	case "as_string":
		if len(argEvaluators) != 1 {
			return nil, fmt.Errorf("as_string() func is expected 1 arg, but given %d args", len(argEvaluators))
		}
		return asStringCallFunc, nil
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
