package middleware

import "github.com/go-park-mail-ru/2023_1_4from5/internal/models"

func UserIsValid(user models.User) bool {
	return user.Login != "" && user.PasswordHash != ""
}

func LoginUserIsValid(user models.LoginUser) bool {
	return user.Login != "" && user.PasswordHash != ""
}
