package models

import "errors"

var (
	WrongPassword                 = errors.New("WrongPassword")
	NotFound                      = errors.New("NotFound")
	InternalError                 = errors.New("InternalError")
	ExpiredToken                  = errors.New("ExpiredToken")
	NoToken                       = errors.New("NoToken")
	NoAuthData                    = errors.New("NoAuthData")
	InvalidToken                  = errors.New("InvalidToken")
	WrongData                     = errors.New("WrongData")
	Unauthorized                  = errors.New("Unauthorized")
	Forbbiden                     = errors.New("Forbidden")
	Unsupported                   = errors.New("Unsupported")
	WrongLoginLength              = errors.New("wrong login length")
	WrongLoginSymbols             = errors.New("wrong login contains black symbols")
	WrongNameLength               = errors.New("wrong username length")
	WrongNameSymbols              = errors.New("wrong username contains black symbols")
	WrongNameHasNoLetter          = errors.New("wrong username has no letter")
	WrongCreatorNameHasNoLetter   = errors.New("wrong creator name has no letter")
	WrongCreatorNameLength        = errors.New("wrong creator name length")
	WrongCreatorDescriptionLength = errors.New("wrong creator description length")
	WrongCreatorNameSymbols       = errors.New("wrong creator name contains black symbols")
	WrongPasswordLength           = errors.New("wrong password length")
	WrongPasswordSymbols          = errors.New("wrong password contains black symbols")
	PasswordHasNoNumber           = errors.New("wrong password has no number")
)
