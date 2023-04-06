package models

// easyjson -all ./internal/models/post.go

import (
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
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

type PostEditData struct {
	Id                     uuid.UUID
	Title                  string      `json:"title"`
	Text                   string      `json:"text"`
	AvailableSubscriptions []uuid.UUID `json:"available_subscriptions"`
}

func (post PostCreationData) IsValid() bool {
	return len(post.Text) != 0 || len(post.Title) != 0 || post.Attachments != nil
}

func (post *Post) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()
	post.Title = sanitizer.Sanitize(post.Title)
	post.Text = sanitizer.Sanitize(post.Text)
	for i, _ := range post.Subscriptions {
		post.Subscriptions[i].Sanitize()
	}
}
