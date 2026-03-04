package errorsx

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInternal        = errors.New("internal server error")
)
