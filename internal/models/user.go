package models

import (
	"github.com/google/uuid"
	"html"
	"time"
	"unicode"
)

// easyjson -all ./internal/models/user.go
const (
	creatorNameMaxLength        = 40
	userLoginMaxLength          = 40
	userLoginMinLength          = 7
	userPasswordMaxLength       = 20
	userPasswordMinLength       = 7
	userNameMaxLength           = 40
	creatorDescriptionMaxLength = 500
)

type User struct {
	Id           uuid.UUID `json:"id"`
	Login        string    `json:"login"    example:"Hacker2003"`
	Name         string    `json:"name" example:"Danila Polyakov"`
	ProfilePhoto uuid.UUID `json:"profile_photo"`
	PasswordHash string    `json:"password_hash" example:"1cbedcfebd7efb060916156dafe1dc4b7007db6b7e2312aeb5eed4a43f54e8f767e7d823b54119771723f87aa0bb05df34806fc598cd889042e4da9a609571c3"`
	Registration time.Time `json:"registration"`
	UserVersion  int64     `json:"user_version"`
}

func (user User) UserLoginIsValid() error {
	if !(len([]rune(user.Login)) >= userLoginMinLength && len([]rune(user.Login)) < userLoginMaxLength) {
		return WrongLoginLength
	}
	for _, c := range user.Login {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !(c == '.') && !(c == '_') && !(c == '-') {
			return WrongLoginSymbols
		}
	}
	return nil
}

func (user User) UserPasswordIsValid() error {
	if len([]rune(user.PasswordHash)) >= userPasswordMaxLength {
		return WrongPasswordLength
	}
	for _, c := range user.PasswordHash {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsPunct(c) {
			return WrongPasswordSymbols
		}
	}

	var (
		hasMinLen = false
		hasNumber = false
	)
	if len([]rune(user.PasswordHash)) >= userPasswordMinLength {
		hasMinLen = true
	}

	for _, char := range user.PasswordHash {
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}
	if !hasMinLen {
		return WrongPasswordLength
	}
	if !hasNumber {
		return PasswordHasNoNumber
	}
	return nil
}

func (user User) UserNameIsValid() error {
	if len(user.Name) == 0 || len([]rune(user.Name)) > userNameMaxLength {
		return WrongNameLength
	}
	hasLetter := false
	for _, c := range user.Name {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !(c == '.') && !(c == '_') && !(c == '-') && !(c == ' ') {
			return WrongNameSymbols
		}
		if unicode.IsLetter(c) {
			hasLetter = true
		}
	}
	if !hasLetter {
		return WrongNameHasNoLetter
	}
	return nil
}

func (user User) UserAuthIsValid() error {
	if err := user.UserLoginIsValid(); err != nil {
		return err
	}
	return user.UserPasswordIsValid()
}

func (user User) UserIsValid() error {
	if err := user.UserLoginIsValid(); err != nil {
		return err
	}
	if err := user.UserPasswordIsValid(); err != nil {
		return err
	}
	return user.UserNameIsValid()
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
	IsCreator    bool      `json:"is_creator"`
	CreatorId    uuid.UUID `json:"creator_id"`
}

type BecameCreatorInfo struct {
	Name        string `json:"name" example:"Danila Polyakov"`
	Description string `json:"description"`
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
	MoneyCount float32   `json:"money_count"`
}

func (becameCreatorInfo *BecameCreatorInfo) IsValid() error {
	if len(becameCreatorInfo.Name) == 0 || len([]rune(becameCreatorInfo.Name)) > creatorNameMaxLength {
		return WrongCreatorNameLength
	}
	if len([]rune(becameCreatorInfo.Description)) > creatorDescriptionMaxLength {
		return WrongCreatorDescriptionLength
	}
	hasLetter := false
	for _, c := range becameCreatorInfo.Name {
		if !unicode.IsLetter(c) && !(c >= 32 && c <= 126) {
			return WrongCreatorNameSymbols
		}
		if unicode.IsLetter(c) {
			hasLetter = true
		}
	}
	if !hasLetter {
		return WrongCreatorNameHasNoLetter
	}
	return nil
}

func (user *User) Sanitize() {
	user.Login = html.EscapeString(user.Login)
	user.Name = html.EscapeString(user.Name)
}

func (userProfile *UserProfile) Sanitize() {
	userProfile.Login = html.EscapeString(userProfile.Login)
	userProfile.Name = html.EscapeString(userProfile.Name)
}
