package evaluator

import (
	"errors"
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
			ret, err := lssComparativeFunc(v1, v2)
			if err != nil || ret {
				return ret, err
			}
			return equalComparativeFunc(v1, v2)
		}, true
	case token.GEQ: // >=
		return func(v1, v2 interface{}) (bool, error) {
			ret, err := gtrComparativeFunc(v1, v2)
			if err != nil || ret {
				return ret, err
			}
			return equalComparativeFunc(v1, v2)
		}, true
	default:
		return nil, false
	}
}

func equalComparativeFunc(v1, v2 interface{}) (bool, error) {
	if b1, b2, ok := isBothBools(v1, v2); ok {
		return b1 == b2, nil
	}
	if s1, s2, ok := isBothStrings(v1, v2); ok {
		return s1 == s2, nil
	}
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		return n1 == n2, nil
	}
	return false, fmt.Errorf("v1[%v]::%T and v2[%v]::%T can not `==` comparatable", v1, v1, v2, v2)
}

func lssComparativeFunc(v1, v2 interface{}) (bool, error) {
	if s1, s2, ok := isBothStrings(v1, v2); ok {
		return s1 < s2, nil
	}
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		return n1 < n2, nil
	}
	return false, fmt.Errorf("v1[%v]::%T and v2[%v]::%T can not `<` comparatable", v1, v1, v2, v2)
}

func gtrComparativeFunc(v1, v2 interface{}) (bool, error) {
	if s1, s2, ok := isBothStrings(v1, v2); ok {
		return s1 > s2, nil
	}
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		return n1 > n2, nil
	}
	return false, fmt.Errorf("v1[%v]::%T and v2[%v]::%T can not `>` comparatable", v1, v1, v2, v2)
}

type logicalFunc func(bool, bool) bool

func getLogicalFunc(op token.Token) (logicalFunc, bool) {
	switch op {
	case token.LAND: // &&
		return func(b1, b2 bool) bool { return b1 && b2 }, true
	case token.LOR: // &&
		return func(b1, b2 bool) bool { return b1 || b2 }, true
	default:
		return nil, false
	}
}

type computableFunc func(interface{}, interface{}) (interface{}, error)

func getComputableFunc(op token.Token) (computableFunc, bool) {
	switch op {
	case token.ADD: // +
		return addComputableFunc, true
	case token.SUB: // -
		return subComputableFunc, true
	case token.MUL: // *
		return mulComputableFunc, true
	case token.QUO: // /
		return quoComputableFunc, true
	default:
		return nil, false
	}
}

func addComputableFunc(v1, v2 interface{}) (interface{}, error) {
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		return n1 + n2, nil
	}
	return false, fmt.Errorf("v1[%v]::%T and v2[%v]::%T can not `+` comparatable", v1, v1, v2, v2)
}

func subComputableFunc(v1, v2 interface{}) (interface{}, error) {
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		return n1 - n2, nil
	}
	return false, fmt.Errorf("v1[%v]::%T and v2[%v]::%T can not `-` comparatable", v1, v1, v2, v2)
}

func mulComputableFunc(v1, v2 interface{}) (interface{}, error) {
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		return n1 * n2, nil
	}
	return false, fmt.Errorf("v1[%v]::%T and v2[%v]::%T can not `*` comparatable", v1, v1, v2, v2)
}

func quoComputableFunc(v1, v2 interface{}) (interface{}, error) {
	if n1, n2, ok := isBothRealNumbers(v1, v2); ok {
		if n2 == 0 {
			return nil, errors.New("divided by 0")
		}
		return n1 / n2, nil
	}
	return false, fmt.Errorf("v1[%v]::%T and v2[%v]::%T can not `*` comparatable", v1, v1, v2, v2)
}
