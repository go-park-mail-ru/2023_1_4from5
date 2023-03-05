package middleware

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strings"
	"time"
)

type Extractor func(r *http.Request) string

/*
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
}*/

func ExtractTokenFromCookie(r *http.Request) string {
	tokenCookie, err := r.Cookie("SSID")
	if err != nil {
		return ""
	}
	token := tokenCookie.Value
	strArr := strings.Split(token, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return strArr[0]
}

func VerifyToken(r *http.Request, extractor Extractor) (*models.Token, error) {
	tokenStr := extractor(r)
	if tokenStr == "" {
		return nil, models.NoToken
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
	return nil, models.NoAuthData
}

func ExtractTokenMetadata(r *http.Request, extractor Extractor) (*models.AccessDetails, error) {
	token, err := VerifyToken(r, extractor)
	if err != nil {
		return nil, err
	}
	exp := token.ExpiresAt
	now := time.Now().Unix()
	fmt.Println("now", now)
	fmt.Println("got", exp)
	if exp < now {
		return nil, models.ExpiredToken
	}
	uid, err := uuid.Parse(token.Id)
	if err != nil {
		return nil, err
	}
	data := &models.AccessDetails{Login: token.Login, Id: uid}
	if data.Login == "" || data.Id.String() == "" {
		return nil, models.InvalidToken
	}

	return data, err
}
