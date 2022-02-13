package utils

import "fmt"

type ErrorF struct {
	CodeError int64  // error occurred after reading Offset bytes
	Message   string // description of error
}

func (err *ErrorF) Error() string {
	return fmt.Sprintf("%s, code Error: %d", err.Message, err.CodeError)
}

func NewError(code int64, msg string) *ErrorF {
	/*e := new(ErrorF)
	e.codeError = code
	e.message = msg*/
	return &ErrorF{
		CodeError: code,
		Message:   msg,
	}
}
