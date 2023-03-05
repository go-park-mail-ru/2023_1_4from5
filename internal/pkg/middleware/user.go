package middleware

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"unicode"
)

func isValid(s string) bool {
	var (
		hasMinLen = false
		hasNumber = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}
	return hasMinLen && hasNumber
}

func UserIsValid(user models.User) bool {
	return len(user.Login) > 7 && len(user.Login) < 20 && isValid(user.PasswordHash)
}
