package usecase

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"net/http"
)

type AuthUsecase struct {
	repo      auth.AuthRepo
	tokenator auth.TokenGenerator
	encrypter auth.Encrypter
}

func NewAuthUsecase(repo auth.AuthRepo, tokenator auth.TokenGenerator) *AuthUsecase {
	return &AuthUsecase{repo: repo, tokenator: tokenator}
}

func (u *AuthUsecase) SignIn(user models.LoginUser) (string, int) {
	if user.Login == "" || user.PasswordHash == "" {
		return "", http.StatusBadRequest
	}

	user.PasswordHash = u.encrypter.EncryptPswd(user.PasswordHash)
	DBUser, status := u.repo.CheckUser(models.User{Login: user.Login, PasswordHash: user.PasswordHash})

	token := u.tokenator.GetToken(DBUser)
	if status == http.StatusOK && token != "" {
		return token, status
	}
	return "", status
}

func (u *AuthUsecase) SignUp(user models.User) (string, int) {
	user.PasswordHash = u.encrypter.EncryptPswd(user.PasswordHash)
	_, st := u.repo.CheckUser(user)
	if st == http.StatusOK || st == http.StatusUnauthorized {
		return "", http.StatusConflict
	}

	NewUser, status := u.repo.CreateUser(user)
	token := u.tokenator.GetToken(NewUser)
	if status == http.StatusOK && token != "" {
		return token, status
	}
	return "", status
}

func (u *AuthUsecase) CheckUser(user models.User) (models.User, int) {
	return u.repo.CheckUser(user)
}
