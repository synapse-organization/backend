package errors

import "errors"

var (
	ErrPasswordInvalid  = StringError{Msg: "password is invalid"}
	ErrorUserNotFound   = StringError{Msg: "user not found"}
	ErrEmailInvalid     = StringError{Msg: "email is invalid"}
	ErrPhoneInvalid     = StringError{Msg: "phone is invalid"}
	ErrFirstNameInvalid = StringError{Msg: "first_name is invalid"}
	ErrLastNameInvalid  = StringError{Msg: "last_name is invalid"}
)

type StringError struct {
	Msg string
}

func (e *StringError) Error() error {
	return errors.New(e.Msg)
}
