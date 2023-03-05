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
	ConflictData  = errors.New("ConflictData")
)
