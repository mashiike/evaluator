package evaluator

import "errors"

//resuered error
var (
	ErrDivideByZero     = errors.New("divide by 0")
	ErrVariableNotFound = errors.New("variable not found")
)

//IsDivideByZero check error DivideByZero
func IsDivideByZero(err error) bool {
	return equalErorr(err, ErrDivideByZero)
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
