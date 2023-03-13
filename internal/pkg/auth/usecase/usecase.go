package usecase

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"net/http"
)

type AuthUsecase struct {
	repo      auth.AuthRepo
	tokenator auth.TokenGenerator
	encryptor auth.Encryptor
}

func NewAuthUsecase(repo auth.AuthRepo, tokenator auth.TokenGenerator, encrypter auth.Encryptor) *AuthUsecase {
	return &AuthUsecase{repo: repo, tokenator: tokenator, encryptor: encrypter}
}

func (u *AuthUsecase) SignIn(user models.LoginUser) (string, int) {
	if user.Login == "" || user.PasswordHash == "" {
		return "", http.StatusBadRequest
	}
	user.PasswordHash = u.encryptor.EncryptPswd(user.PasswordHash)
	DBUser, status := u.repo.CheckUser(models.User{Login: user.Login, PasswordHash: user.PasswordHash})
	fmt.Println("Usecase  ", status)
	token := u.tokenator.GetToken(DBUser)
	if status == nil && token != "" {
		return token, http.StatusOK
	}
	return "", http.StatusUnauthorized
}

func (u *AuthUsecase) SignUp(user models.User) (string, int) {
	user.PasswordHash = u.encryptor.EncryptPswd(user.PasswordHash)

	_, err := u.CheckUser(user)
	if err == nil || errors.Is(err, models.WrongPassword) {
		return "", http.StatusConflict
	}
	NewUser, err := u.repo.CreateUser(user)
	token := u.tokenator.GetToken(NewUser)
	if err == nil && token != "" {
		return token, http.StatusOK
	}
	return "", http.StatusInternalServerError
}

func (u *AuthUsecase) CheckUser(user models.User) (models.User, error) {
	return u.repo.CheckUser(user)
}

func (u *AuthUsecase) Logout(details models.AccessDetails) (int, error) {
	return u.repo.IncUserVersion(details.Id)
}

func (u *AuthUsecase) CheckUserVersion(details models.AccessDetails) (int, error) {
	return u.repo.CheckUserVersion(details)
}
