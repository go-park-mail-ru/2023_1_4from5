package token

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

func ExtractCSRFTokenFromCookie(r *http.Request) string {
	tokenCookie, err := r.Cookie("X-CSRF-Token")
	if err != nil {
		return ""
	}
	token := tokenCookie.Value
	return token
}

func VerifyCSRFToken(r *http.Request) (*models.Token, error) {
	tokenStr := ExtractCSRFTokenFromCookie(r)
	if tokenStr == "" {
		return nil, models.NoToken
	}
	token, err := jwt.ParseWithClaims(tokenStr, &models.Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("CSRF_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claim := token.Claims
	claims, ok := claim.(*models.Token)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, models.NoAuthData
}

func ExtractCSRFTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := VerifyCSRFToken(r)
	if err != nil {
		return nil, err
	}

	if err = token.Valid(); err != nil {
		return nil, models.ExpiredToken
	}

	uid, err := uuid.Parse(token.Id)
	if err != nil {
		return nil, err
	}

	data := &models.AccessDetails{Login: token.Login, Id: uid, UserVersion: token.UserVersion}
	if data.Login == "" || data.Id.String() == "" {
		return nil, models.InvalidToken
	}
	return data, err
}

func GetCSRFToken(user models.User) (string, error) {
	tokenModel := models.Token{
		Login:       user.Login,
		Id:          user.Id.String(),
		UserVersion: user.UserVersion,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}
	secretKey, flag := os.LookupEnv("CSRF_SECRET")
	if !flag {
		fmt.Println(flag)
		fmt.Println(secretKey)
		return "", errors.New("NoSecretKey")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenModel)

	jwtCookie, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println(err)
		return "", errors.New("NoSecretKey")
	}
	return jwtCookie, nil
}
