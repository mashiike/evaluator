package evaluator

import (
	"fmt"
	"regexp"
	"strings"
)

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
	case "string_contains": // string_contains(string,string)
		if len(argEvaluators) != 2 {
			return nil, newNumOfArgumentsMismatchError(funcName, 2, len(argEvaluators))
		}
		return stringContainsCallFunc, nil
	case "regexp_match": // regexp_match(string,string)
		if len(argEvaluators) != 2 {
			return nil, newNumOfArgumentsMismatchError(funcName, 2, len(argEvaluators))
		}
		return regexMatchCallFunc, nil
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

func stringContainsCallFunc(args ...interface{}) (interface{}, error) {
	s1, s2, ok := isBothStrings(args[0], args[1])
	if !ok {
		return nil, fmt.Errorf("string_contains(v1[%v]::%T,v2[%v]::%T) can not eval", args[0], args[0], args[1], args[1])
	}
	return strings.Contains(s1, s2), nil
}

type regexpCacheEntry struct {
	pattern string
	reg     *regexp.Regexp
}

var regexpCache = make([]regexpCacheEntry, 0, 10)

func regexMatchCallFunc(args ...interface{}) (interface{}, error) {
	s1, s2, ok := isBothStrings(args[0], args[1])
	if !ok {
		return nil, fmt.Errorf("regex_match(v1[%v]::%T,v2[%v]::%T) can not eval", args[0], args[0], args[1], args[1])
	}
	var reg *regexp.Regexp
	for _, entry := range regexpCache {
		if entry.pattern == s2 {
			reg = entry.reg
			break
		}
	}
	if reg == nil {
		var err error
		reg, err = regexp.Compile(s2)
		if err != nil {
			return nil, fmt.Errorf("regex_match(v1[%v]::%T,v2[%v]::%T) pattern can not compile: %w", args[0], args[0], args[1], args[1], err)
		}
		regexpCache = append(regexpCache, regexpCacheEntry{
			pattern: s2,
			reg:     reg,
		})
		if len(regexpCache) > 100 {
			regexpCache = regexpCache[1:]
		}
	}
	return reg.Match([]byte(s1)), nil
}
