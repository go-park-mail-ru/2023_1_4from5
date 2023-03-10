package auth

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/auth_mock.go -package=mock

type AuthUsecase interface {
	SignIn(user models.LoginUser) (string, error)
	SignUp(user models.User) (string, error)
}

type AuthRepo interface {
	CreateUser(user models.User) (models.User, error)
	CheckUser(user models.User) (models.User, error)
}

type TokenGenerator interface {
	GetToken(user models.User) (string, error)
}

type Encrypter interface {
	EncryptPswd(pswd string) string
}
