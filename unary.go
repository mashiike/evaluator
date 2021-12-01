package evaluator

import (
	"fmt"
	"go/token"
)

type unaryFunc func(interface{}) (interface{}, error)

func getUnaryFunc(op token.Token) (unaryFunc, bool) {
	switch op {
	case token.NOT: // !
		return notUnaryFunc, true
	default:
		return nil, false
	}
}

func notUnaryFunc(v interface{}) (interface{}, error) {
	b, ok := asBool(v)
	if !ok {
		return nil, fmt.Errorf("v[%v]::%T can not `!` operation", v, v)
	}
	return !b, nil
}
