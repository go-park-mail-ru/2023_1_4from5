package models

import "errors"

var (
	WrongPassword = errors.New("WrongPassword")
	NotFound      = errors.New("NotFound")
	InternalError = errors.New("InternalError")
)
