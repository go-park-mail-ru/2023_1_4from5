package middleware

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"unicode"
)

func isValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func UserIsValid(user models.User) bool {
	return len(user.Login) > 7 && len(user.Login) < 20 && isValid(user.PasswordHash)
}

func LoginUserIsValid(user models.LoginUser) bool {
	return len(user.Login) > 7 && len(user.Login) < 20 && isValid(user.PasswordHash)
}

//TODO: userVersionCheck here
//Can middleware work with repo+usecase?
// Another endpoint for it?
// don't know how to realize
