package models

// easyjson -all ./internal/models/post.go

import (
	"github.com/google/uuid"
	"html"
	"time"
)

type Post struct {
	Id            uuid.UUID      `json:"id"`
	Creator       uuid.UUID      `json:"creator"`
	Creation      time.Time      `json:"creation_date"`
	LikesCount    int            `json:"likes_count"`
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
	Id                     uuid.UUID   `json:"-"`
	Title                  string      `json:"title"`
	Text                   string      `json:"text"`
	AvailableSubscriptions []uuid.UUID `json:"available_subscriptions"`
}

func (post PostCreationData) IsValid() bool {
	return len(post.Text) != 0 || len(post.Title) != 0 || post.Attachments != nil
}

func (post *Post) Sanitize() {
	post.Title = html.EscapeString(post.Title)
	post.Text = html.EscapeString(post.Text)
	for i := range post.Subscriptions {
		post.Subscriptions[i].Sanitize()
	}
}
