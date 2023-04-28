package models

import "github.com/google/uuid"

// easyjson -all ./internal/models/like.go

type Like struct {
	LikesCount int64     `json:"likes_count"`
	PostID     uuid.UUID `json:"post_id"`
}
