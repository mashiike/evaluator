package evaluator

import "errors"

//resuered error
var (
	ErrDevideByZero     = errors.New("devide by 0")
	ErrVariableNotFound = errors.New("variable not found")
)

//IsDevideByZero check error DevideByZero
func IsDevideByZero(err error) bool {
	return equalErorr(err, ErrDevideByZero)
}

//IsVariableNotFound check error VariableNotFound
func IsVariableNotFound(err error) bool {
	return equalErorr(err, ErrVariableNotFound)
}

func equalErorr(err, other error) bool {
	if err == other {
		return true
	}
	if err := errors.Unwrap(err); err != nil {
		return equalErorr(err, other)
	}
	return false
}
