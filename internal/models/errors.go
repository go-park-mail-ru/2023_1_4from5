package models

import "errors"

var (
	WrongPassword = errors.New("WrongPassword")
	NotFound      = errors.New("NotFound")
	InternalError = errors.New("InternalError")
	ExpiredToken  = errors.New("ExpiredToken")
	NoToken       = errors.New("NoToken")
	NoAuthData    = errors.New("NoAuthData")
	InvalidToken  = errors.New("InvalidToken")
	WrongData     = errors.New("WrongData")
	Unauthorized  = errors.New("Unauthorized")
	Forbbiden     = errors.New("Forbidden")
	Unsupported   = errors.New("Unsupported")
)
