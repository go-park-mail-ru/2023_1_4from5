package auth

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/auth_mock.go -package=mock

type AuthUsecase interface {
	SignIn(user models.LoginUser) (string, int)
	SignUp(user models.User) (string, int)
}

type AuthRepo interface {
	CreateUser(user models.User) (models.User, int)
}
