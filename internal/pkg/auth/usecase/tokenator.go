package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

type Tokenator struct {
}

func NewTokenator() *Tokenator {
	return &Tokenator{}
}

func (t *Tokenator) GetToken(user models.User) (string, error) {
	tokenModel := models.Token{
		Login: user.Login,
		Id:    user.Id.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
		},
	}
	SecretKey, flag := os.LookupEnv("SECRET")
	if !flag {
		return "", errors.New("NoSecretKey")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenModel)

	jwtCookie, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", errors.New("NoSecretKey")
	}
	return jwtCookie, nil
}
