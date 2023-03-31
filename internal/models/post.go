package models

// easyjson -all ./internal/models/post.go

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	Id            uuid.UUID      `json:"id"`
	Creator       uuid.UUID      `json:"creator"`
	Creation      time.Time      `json:"creation_date"`
	Title         string         `json:"title"`
	Text          string         `json:"text"`
	IsAvailable   bool           `json:"is_available"`
	IsLiked       bool           `json:"is_liked"`
	Attachments   []Attachment   `json:"attachments"`
	Subscriptions []Subscription `json:"subscriptions"`
}

//easyjson:skip
type PostCreationData struct {
	Id                     uuid.UUID
	Creator                uuid.UUID
	Title                  string
	Text                   string
	Attachments            []AttachmentData
	AvailableSubscriptions []uuid.UUID
}

func (post PostCreationData) IsValid() bool {
	return len(post.Text) != 0 || len(post.Title) != 0 || post.Attachments != nil
}
