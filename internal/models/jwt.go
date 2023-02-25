package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// easyjson:skip
type Token struct {
	Login string
	Id    string
	jwt.StandardClaims
}

type TokenView struct {
	Token string `json:"token"`
}

type AccessDetails struct {
	Login string
	Id    uuid.UUID
}
