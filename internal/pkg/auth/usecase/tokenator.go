package usecase

import (
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

func (t *Tokenator) GetToken(user models.User) string {
	tokenModel := models.Token{
		Login: user.Login,
		Id:    user.Id.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
		},
	}
	//храним в переменной окружения SECRET
	SecretKey, flag := os.LookupEnv("SECRET")
	if !flag {
		return "no secret key"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenModel)

	jwtCookie, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "no secret key"
	}
	return jwtCookie
}
