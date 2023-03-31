package usecase

import (
	"context"
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

func (t *Tokenator) GetJWTToken(ctx context.Context, user models.User) (string, error) {
	tokenModel := models.Token{
		Login:       user.Login,
		Id:          user.Id.String(),
		UserVersion: user.UserVersion,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
		},
	}
	secretKey, flag := os.LookupEnv("TOKEN_SECRET")
	if !flag {
		return "", errors.New("NoSecretKey")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenModel)

	jwtCookie, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.New("NoSecretKey")
	}
	return jwtCookie, nil
}
