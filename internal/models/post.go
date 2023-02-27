package models

// easyjson -all ./internal/models/post.go

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	Id       uuid.UUID `json:"id"`
	Creator  uuid.UUID `json:"creator"`
	Creation time.Time `json:"creation_date"`
	Title    string    `json:"title"`
	Text     string    `json:"text"`
}
