package auth

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/auth_mock.go -package=mock

type AuthUsecase interface {
	SignIn(user models.LoginUser) (string, int)
	SignUp(user models.User) (string, int)
	Logout(details models.AccessDetails) (int, error)
	CheckUserVersion(details models.AccessDetails) (int, error)
}

type AuthRepo interface {
	CreateUser(user models.User) (models.User, error)
	CheckUser(user models.User) (models.User, error)
	IncUserVersion(userId uuid.UUID) (int, error)
	CheckUserVersion(details models.AccessDetails) (int, error)
}

type TokenGenerator interface {
	GetToken(user models.User) string
}

type Encryptor interface {
	EncryptPswd(pwd string) string
}
