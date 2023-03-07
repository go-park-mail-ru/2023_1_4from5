package models

import (
	"github.com/google/uuid"
	"time"
)

// easyjson -all ./internal/models/user.go

type User struct {
	Id           uuid.UUID `json:"id"`
	Login        string    `json:"login"    example:"Hacker2003"`
	Name         string    `json:"name" example:"Danila Polyakov"`
	ProfilePhoto string    `json:"profile_photo"`
	PasswordHash string    `json:"password_hash" example:"1cbedcfebd7efb060916156dafe1dc4b7007db6b7e2312aeb5eed4a43f54e8f767e7d823b54119771723f87aa0bb05df34806fc598cd889042e4da9a609571c3"`
	Registration time.Time `json:"registration"`
	UserVersion  int       `json:"user_version"`
}

type LoginUser struct {
	Login        string `json:"login"    example:"Hacker2003"`
	PasswordHash string `json:"password_hash" example:"1cbedcfebd7efb060916156dafe1dc4b7007db6b7e2312aeb5eed4a43f54e8f767e7d823b54119771723f87aa0bb05df34806fc598cd889042e4da9a609571c3"`
}

type UserProfile struct {
	Login        string    `json:"login"    example:"Hacker2003"`
	Name         string    `json:"name" example:"Danila Polyakov"`
	ProfilePhoto string    `json:"profile_photo"`
	Registration time.Time `json:"registration"`
}

type UserHomePage struct {
	Name         string    `json:"name" example:"Danila Polyakov"`
	ProfilePhoto string    `json:"profile_photo"`
	Posts        []Post    `json:"posts"`
	IsCreator    bool      `json:"is_creator"`
	CreatorId    uuid.UUID `json:"creator_id"`
}
