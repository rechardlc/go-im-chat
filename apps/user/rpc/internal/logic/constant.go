package logic

import "errors"

var (
	ErrInvalidPhone    = errors.New("invalid phone")
	ErrInvalidPassword = errors.New("invalid password")
)