package models

import (
	"github.com/google/uuid"
	"html"
	"time"
	"unicode"
)

// easyjson -all ./internal/models/user.go

type User struct {
	Id           uuid.UUID `json:"id"`
	Login        string    `json:"login"    example:"Hacker2003"`
	Name         string    `json:"name" example:"Danila Polyakov"`
	ProfilePhoto uuid.UUID `json:"profile_photo"`
	PasswordHash string    `json:"password_hash" example:"1cbedcfebd7efb060916156dafe1dc4b7007db6b7e2312aeb5eed4a43f54e8f767e7d823b54119771723f87aa0bb05df34806fc598cd889042e4da9a609571c3"`
	Registration time.Time `json:"registration"`
	UserVersion  int       `json:"user_version"`
}

func (user User) UserLoginIsValid() bool {
	return len(user.Login) >= 7 && len(user.Login) < 40
}

func (user User) UserPasswordIsValid() bool {
	var (
		hasMinLen = false
		hasNumber = false
	)
	if len(user.PasswordHash) >= 7 {
		hasMinLen = true
	}
	for _, char := range user.PasswordHash {
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}
	return hasMinLen && hasNumber

}

func (user User) UserNameIsValid() bool {
	return len(user.Name) >= 7 && len(user.Name) < 40
}

func (user User) UserAuthIsValid() bool {
	return user.UserLoginIsValid() && user.UserPasswordIsValid()
}

func (user User) UserIsValid() bool {
	return user.UserLoginIsValid() && user.UserPasswordIsValid() && user.UserNameIsValid()
}

type LoginUser struct {
	Login        string `json:"login"    example:"Hacker2003"`
	PasswordHash string `json:"password_hash" example:"1cbedcfebd7efb060916156dafe1dc4b7007db6b7e2312aeb5eed4a43f54e8f767e7d823b54119771723f87aa0bb05df34806fc598cd889042e4da9a609571c3"`
}

type UserProfile struct {
	Login        string    `json:"login"    example:"Hacker2003"`
	Name         string    `json:"name" example:"Danila Polyakov"`
	ProfilePhoto uuid.UUID `json:"profile_photo"`
	Registration time.Time `json:"registration"`
}

type UserHomePage struct {
	Name         string    `json:"name" example:"Danila Polyakov"`
	ProfilePhoto uuid.UUID `json:"profile_photo"`
	Posts        []Post    `json:"posts"`
	IsCreator    bool      `json:"is_creator"`
	CreatorId    uuid.UUID `json:"creator_id"`
}

type UpdatePasswordInfo struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}
type UpdateProfileInfo struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Donate struct {
	CreatorID  uuid.UUID `json:"creator_id"`
	MoneyCount int       `json:"money_count"`
}

func (user *User) Sanitize() {
	user.Login = html.EscapeString(user.Login)
	user.Name = html.EscapeString(user.Name)
}

func (userProfile *UserProfile) Sanitize() {
	userProfile.Login = html.EscapeString(userProfile.Login)
	userProfile.Name = html.EscapeString(userProfile.Name)
}

func (userHomePage *UserHomePage) Sanitize() {
	userHomePage.Name = html.EscapeString(userHomePage.Name)
	for i := range userHomePage.Posts {
		userHomePage.Posts[i].Sanitize()
	}
}
