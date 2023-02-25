package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
	"time"
)

type AccessDetails struct {
	Login string
}

type Extracter func(r *http.Request) string

func ExtractToken(r *http.Request) string {
	token := models.TokenView{}
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		return ""
	}

	strArr := strings.Split(token.Token, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return strArr[0]
}

func ExtractTokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	strArr := strings.Split(token, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return strArr[0]
}

func VerifyToken(r *http.Request, extracter Extracter) (*models.Token, error) {
	tokenStr := extracter(r)
	if tokenStr == "" {
		return nil, errors.New("no token")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &models.Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claim := token.Claims
	claims, ok := claim.(*models.Token)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("no auth data")
}

func ExtractTokenMetadata(r *http.Request, extracter Extracter) (*AccessDetails, error) {
	token, err := VerifyToken(r, extracter)
	if err != nil {
		return nil, err
	}
	exp := token.ExpiresAt
	now := time.Now().Unix()
	if exp < now {
		return nil, errors.New("token expired")
	}
	data := &AccessDetails{Login: token.Login}
	if data.Login == "" {
		return nil, errors.New("invalid token")
	}

	return data, err
}
