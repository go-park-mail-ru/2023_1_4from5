package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthUsecase struct {
	repo      auth.AuthRepo
	tokenator auth.TokenGenerator
	encrypter auth.Encrypter
	logger    *zap.SugaredLogger
}

func NewAuthUsecase(repo auth.AuthRepo, tokenator auth.TokenGenerator, encrypter auth.Encrypter, logger *zap.SugaredLogger) *AuthUsecase {
	return &AuthUsecase{
		repo:      repo,
		tokenator: tokenator,
		encrypter: encrypter,
		logger:    logger,
	}
}

func (u *AuthUsecase) EncryptPwd(ctx context.Context, pwd string) string {
	return u.encrypter.EncryptPswd(ctx, pwd)
}

func (u *AuthUsecase) IncUserVersion(ctx context.Context, details models.AccessDetails) (int64, error) {
	return u.repo.IncUserVersion(ctx, details.Id)
}

func (u *AuthUsecase) SignIn(ctx context.Context, user models.LoginUser) (string, error) {
	user.PasswordHash = u.encrypter.EncryptPswd(ctx, user.PasswordHash)
	dbUser, dbErr := u.repo.CheckUser(ctx, models.User{Login: user.Login, PasswordHash: user.PasswordHash})
	token, err := u.tokenator.GetJWTToken(ctx, dbUser)
	if dbErr == nil && err == nil {
		return token, nil
	}
	return "", dbErr
}

func (u *AuthUsecase) CheckUser(ctx context.Context, user models.User) (models.User, error) {
	user.PasswordHash = u.encrypter.EncryptPswd(ctx, user.PasswordHash)
	return u.repo.CheckUser(ctx, user)
}

func (u *AuthUsecase) SignUp(ctx context.Context, user models.User) (string, error) {
	user.PasswordHash = u.encrypter.EncryptPswd(ctx, user.PasswordHash)
	_, err := u.repo.CheckUser(ctx, user)
	if err == nil || errors.Is(err, models.WrongPassword) {
		return "", models.WrongData
	}
	user.Id = uuid.New()
	newUser, dbErr := u.repo.CreateUser(ctx, user)
	token, tokenErr := u.tokenator.GetJWTToken(ctx, newUser)
	if dbErr == nil && tokenErr == nil {
		return token, nil
	}
	return "", models.InternalError
}

func (u *AuthUsecase) CheckUserVersion(ctx context.Context, details models.AccessDetails) (int64, error) {
	return u.repo.CheckUserVersion(ctx, details)
}
