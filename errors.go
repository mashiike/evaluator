package evaluator

import "errors"

//reserved error
var (
	ErrDivideByZero     = errors.New("divide by 0")
	ErrVariableNotFound = errors.New("variable not found")
)

//IsDivideByZero check error DivideByZero
func IsDivideByZero(err error) bool {
	return equalError(err, ErrDivideByZero)
}

//IsVariableNotFound check error VariableNotFound
func IsVariableNotFound(err error) bool {
	return equalError(err, ErrVariableNotFound)
}

func equalError(err, other error) bool {
	if err == other {
		return true
	}
	if err := errors.Unwrap(err); err != nil {
		return equalError(err, other)
	}
	return false
}
