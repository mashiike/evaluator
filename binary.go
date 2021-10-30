package evaluator

import (
	"fmt"
	"go/token"
)

type comparativeFunc func(interface{}, interface{}) (bool, error)

func getComparativeFunc(op token.Token) (comparativeFunc, bool) {
	switch op {
	case token.EQL, token.ASSIGN: // == or =
		return equalComparativeFunc, true
	case token.LSS: // <
		return lssComparativeFunc, true
	case token.GTR: // >
		return gtrComparativeFunc, true
	case token.NEQ: // >
		return func(v1, v2 interface{}) (bool, error) {
			ret, err := equalComparativeFunc(v1, v2)
			return !ret, err
		}, true
	case token.LEQ: // <=
		return func(v1, v2 interface{}) (bool, error) {
			ret, err := equalComparativeFunc(v1, v2)
			if err != nil || ret {
				return ret, err
			}
			return lssComparativeFunc(v1, v2)
		}, true
	case token.GEQ: // >=
		return func(v1, v2 interface{}) (bool, error) {
			ret, err := equalComparativeFunc(v1, v2)
			if err != nil || ret {
				return ret, err
			}
			return gtrComparativeFunc(v1, v2)
		}, true
	default:
		return nil, false
	}
}

func isBothStrings(v1, v2 interface{}) (s1, s2 string, ok bool) {
	s1, ok = v1.(string)
	if !ok {
		return
	}
	s2, ok = v2.(string)
	return
}

func isBothRealNumbers(v1, v2 interface{}) (n1, n2 float64, ok bool) {
	n1, ok = isRealNumber(v1)
	if !ok {
		return
	}
	n2, ok = isRealNumber(v2)
	return
}

func isRealNumber(v interface{}) (float64, bool) {
	switch v := v.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

func equalComparativeFunc(v1, v2 interface{}) (bool, error) {
	if s1, s2, ok := isBothStrings(v1, v2); ok {
		return s1 == s2, nil
	}
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		return n1 == n2, nil
	}
	return false, fmt.Errorf("v1[%v] and v2[%v] can not `==` comparatable", v1, v2)
}

func lssComparativeFunc(v1, v2 interface{}) (bool, error) {
	if s1, s2, ok := isBothStrings(v1, v2); ok {
		return s1 < s2, nil
	}
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		return n1 < n2, nil
	}
	return false, fmt.Errorf("v1[%v] and v2[%v] can not `<` comparatable", v1, v2)
}

func gtrComparativeFunc(v1, v2 interface{}) (bool, error) {
	if s1, s2, ok := isBothStrings(v1, v2); ok {
		return s1 > s2, nil
	}
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		return n1 > n2, nil
	}
	return false, fmt.Errorf("v1[%v] and v2[%v] can not `>` comparatable", v1, v2)
}
