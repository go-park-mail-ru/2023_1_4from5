package models

import (
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

// easyjson -all ./internal/models/creator.go

type Creator struct {
	Id             uuid.UUID `json:"creator_id"`
	UserId         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	CoverPhoto     uuid.UUID `json:"cover_photo"`
	FollowersCount int       `json:"followers_count"`
	Description    string    `json:"description"`
	PostsCount     int       `json:"posts_count"`
}

type CreatorPage struct {
	CreatorInfo   Creator        `json:"creator_info"`
	Aim           Aim            `json:"aim"`
	IsMyPage      bool           `json:"is_my_page"`
	Posts         []Post         `json:"posts"`
	Subscriptions []Subscription `json:"subscriptions"`
}

type Aim struct {
	Creator     uuid.UUID `json:"creator_id"`
	Description string    `json:"description"`
	MoneyNeeded int       `json:"money_needed"`
	MoneyGot    int       `json:"money_got"`
}

func (creator *Creator) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()
	creator.Name = sanitizer.Sanitize(creator.Name)
	creator.Description = sanitizer.Sanitize(creator.Description)
}

func (aim *Aim) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()
	aim.Description = sanitizer.Sanitize(aim.Description)
}

func (page *CreatorPage) Sanitize() {
	page.CreatorInfo.Sanitize()
	page.Aim.Sanitize()
	for i, _ := range page.Posts {
		page.Posts[i].Sanitize()
	}
	for i, _ := range page.Subscriptions {
		page.Subscriptions[i].Sanitize()
	}
}
