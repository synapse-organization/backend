package errors

import "errors"

var (
	ErrPasswordInvalid   = StringError{Msg: "password is invalid"}
	ErrorUserNotFound    = StringError{Msg: "user not found"}
	ErrEmailInvalid      = StringError{Msg: "email is invalid"}
	ErrPhoneInvalid      = StringError{Msg: "phone is invalid"}
	ErrFirstNameInvalid  = StringError{Msg: "first_name is invalid"}
	ErrLastNameInvalid   = StringError{Msg: "last_name is invalid"}
	ErrPasswordIncorrect = StringError{Msg: "password is incorrect"}
	ErrEmailExists       = StringError{Msg: "email already exists"}
	ErrStartTimeInvalid  = StringError{Msg: "start time is invalid"}
	ErrEndTimeInvalid    = StringError{Msg: "end time is invalid"}
	ErrImageInvalid		= StringError{Msg: "image doesn't exist"}
)

type StringError struct {
	Msg string
}

func (e *StringError) Error() error {
	return errors.New(e.Msg)
}
