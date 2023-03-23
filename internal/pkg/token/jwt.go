package token

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"net/http"
	"os"
)

func ExtractJWTTokenFromCookie(r *http.Request) string {
	tokenCookie, err := r.Cookie("SSID")
	if err != nil {
		return ""
	}
	token := tokenCookie.Value
	return token
}

func VerifyJWTToken(r *http.Request) (*models.Token, error) {
	tokenStr := ExtractJWTTokenFromCookie(r)
	if tokenStr == "" {
		return nil, models.NoToken
	}
	token, err := jwt.ParseWithClaims(tokenStr, &models.Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
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

func ExtractJWTTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := VerifyJWTToken(r)
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
