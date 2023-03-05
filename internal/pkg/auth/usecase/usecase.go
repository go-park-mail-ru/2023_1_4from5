package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
)

type AuthUsecase struct {
	repo      auth.AuthRepo
	tokenator auth.TokenGenerator
	encrypter auth.Encrypter
}

func NewAuthUsecase(repo auth.AuthRepo, tokenator auth.TokenGenerator, encrypter auth.Encrypter) *AuthUsecase {
	return &AuthUsecase{repo: repo, tokenator: tokenator, encrypter: encrypter}
}

func (u *AuthUsecase) SignIn(user models.LoginUser) (string, error) {
	user.PasswordHash = u.encrypter.EncryptPswd(user.PasswordHash)
	DBUser, status := u.repo.CheckUser(models.User{Login: user.Login, PasswordHash: user.PasswordHash})
	token := u.tokenator.GetToken(DBUser)
	if status == nil && token != "" {
		return token, nil
	}
	return "", models.NotFound
}

func (u *AuthUsecase) SignUp(user models.User) (string, error) {
	user.PasswordHash = u.encrypter.EncryptPswd(user.PasswordHash)

	_, err := u.repo.CheckUser(user)
	if err == nil || errors.Is(err, models.WrongPassword) {
		return "", models.ConflictData
	}
	NewUser, err := u.repo.CreateUser(user)
	token := u.tokenator.GetToken(NewUser)
	if err == nil && token != "" {
		return token, nil
	}
	return "", models.InternalError
}
