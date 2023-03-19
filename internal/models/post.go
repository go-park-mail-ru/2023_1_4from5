package models

// easyjson -all ./internal/models/post.go

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	Id          uuid.UUID   `json:"id"`
	Creator     uuid.UUID   `json:"creator"`
	Creation    time.Time   `json:"creation_date"`
	Title       string      `json:"title"`
	Text        string      `json:"text"`
	IsAvailable bool        `json:"is_available"`
	Attachments []uuid.UUID `json:"attachments,omitempty"`
}

type PostCreationData struct {
	Id                     uuid.UUID        `json:"id"`
	Creator                uuid.UUID        `json:"creator"`
	Title                  string           `json:"title"`
	Text                   string           `json:"text"`
	Attachments            []AttachmentData `json:"attachments,omitempty"`
	AvailableSubscriptions uuid.UUID        `json:"available_subscriptions"`
}

func (post PostCreationData) IsValid() bool {
	return len(post.Text) != 0 || len(post.Title) != 0 || post.Attachments != nil
}
