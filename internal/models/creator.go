package models

import "github.com/google/uuid"

// easyjson -all ./internal/models/creator.go

type Creator struct {
	Id             uuid.UUID `json:"id"`
	User           uuid.UUID `json:"user"`
	CoverPhoto     string    `json:"cover_photo"`
	FollowersCount int       `json:"followers_count"`
	Description    string    `json:"description"`
	PostsCount     int       `json:"posts_count"`
}
