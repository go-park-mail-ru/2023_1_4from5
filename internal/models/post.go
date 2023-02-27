package models

// easyjson -all ./internal/models/post.go

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	Id           uuid.UUID `json:"id"`
	CreatorId    uuid.UUID `json:"creator_id"`
	CreationDate time.Time `json:"creation_date"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
}
