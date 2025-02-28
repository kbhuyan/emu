package emu

import "fmt"

type emuError struct {
	msg string
}

func (e *emuError) Error() string {
	return e.msg
}

func NewError(msg string) *emuError {
	return &emuError{
		msg: msg,
	}
}

func (e *emuError) Errorf(format string, args ...interface{}) error {
	e.msg = fmt.Sprintf(format, args...)
	return e
}
