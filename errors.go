package evaluator

import (
	"errors"
	"fmt"
)

//reserved error
var (
	ErrDivideByZero     = errors.New("divide by 0")
	ErrVariableNotFound = errors.New("variable not found")
)

//NumOfArgumentsMismatchError is an error that occurs when the number of arguments of the called function is different.
type NumOfArgumentsMismatchError struct {
	FunctionName string
	Expected     int
	Given        int
}

func (e *NumOfArgumentsMismatchError) Error() string {
	return fmt.Sprintf("%s() func is expected %d arg, but given %d args", e.FunctionName, e.Expected, e.Given)
}

func newNumOfArgumentsMismatchError(functionName string, expected, given int) *NumOfArgumentsMismatchError {
	return &NumOfArgumentsMismatchError{
		FunctionName: functionName,
		Expected:     expected,
		Given:        given,
	}
}

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
