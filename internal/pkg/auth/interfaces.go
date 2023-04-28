package auth

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/auth_mock.go -package=mock

type AuthUsecase interface {
	SignIn(ctx context.Context, user models.LoginUser) (string, error)
	SignUp(ctx context.Context, user models.User) (string, error)
	IncUserVersion(ctx context.Context, details models.AccessDetails) (int64, error)
	CheckUser(ctx context.Context, user models.User) (models.User, error)
	EncryptPwd(ctx context.Context, pwd string) string
	CheckUserVersion(ctx context.Context, details models.AccessDetails) (int64, error)
}

type AuthRepo interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	CheckUser(ctx context.Context, user models.User) (models.User, error)
	IncUserVersion(ctx context.Context, userId uuid.UUID) (int64, error)
	CheckUserVersion(ctx context.Context, details models.AccessDetails) (int64, error)
}

type TokenGenerator interface {
	GetJWTToken(ctx context.Context, user models.User) (string, error)
}

type Encrypter interface {
	EncryptPswd(ctx context.Context, pswd string) string
}
