package models

import "github.com/google/uuid"

// easyjson -all ./internal/models/creator.go

type Creator struct {
	Id             uuid.UUID `json:"creator_id"`
	UserId         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	CoverPhoto     string    `json:"cover_photo"`
	FollowersCount int       `json:"followers_count"`
	Description    string    `json:"description"`
	PostsCount     int       `json:"posts_count"`
}

type CreatorPage struct {
	CreatorInfo Creator `json:"creator_info"`
	IsMyPage    bool    `json:"is_my_page"`
	Posts       []Post  `json:"posts"`
	//Subscriptions [] Subscription `json:"subscriptions"`
}
