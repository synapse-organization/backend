package errors

import "errors"

var (
	ErrorUserNotFound    = StringError{Msg: "کاربر یافت نشد"}
	ErrEmailInvalid      = StringError{Msg: "ایمیل نامعتبر است"}
	ErrPhoneInvalid      = StringError{Msg: "شماره تلفن نامعتبر است"}
	ErrFirstNameInvalid  = StringError{Msg: "نام نامعتبر است"}
	ErrLastNameInvalid   = StringError{Msg: "نام خانوادگی نامعتبر است"}
	ErrPasswordIncorrect = StringError{Msg: "رمز عبور اشتباه است"}
	ErrPasswordNotMatch  = StringError{Msg: "رمز عبور ها مطابقت ندارند"}
	ErrEmailExists       = StringError{Msg: "این ایمیل وجود دارد"}
	ErrStartTimeInvalid  = StringError{Msg: "زمان شروع نامعتبر است"}
	ErrEndTimeInvalid    = StringError{Msg: "زمان پایان نامعتبر است"}
	ErrImageInvalid      = StringError{Msg: "تصویر نامعتبر است"}
	ErrForbidden         = StringError{Msg: "دسترسی غیر مجاز"}
	ErrUnableToGetUser   = StringError{Msg: "خطایی در گرفتن اطلاعات کاربر رخ داده است"}
	ErrBadRequest        = StringError{Msg: "درخواست نامعتبر است"}
	ErrDidntLogin        = StringError{Msg: "شما وارد نشده اید"}
	ErrInternalError     = StringError{Msg: "خطای داخلی"}
	ErrNotEnoughBalance  = StringError{Msg: "موجودی کافی نیست"}
	ErrCapacityInvalid   = StringError{Msg: "ظرفیت وارد شده نامعتبر است"}
)

type StringError struct {
	Msg string
}

func (e *StringError) Error() error {
	return errors.New(e.Msg)
}
