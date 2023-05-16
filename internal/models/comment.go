package models

import (
	"github.com/google/uuid"
	"html"
	"time"
)

// easyjson -all ./internal/models/comment.go

type Comment struct {
	CommentID  uuid.UUID `json:"comment_id,omitempty"`
	UserID     uuid.UUID `json:"user_id,omitempty"`
	UserPhoto  uuid.UUID `json:"user_photo,omitempty"`
	PostID     uuid.UUID `json:"post_id"`
	Text       string    `json:"text,omitempty"`
	Creation   time.Time `json:"creation,omitempty"`
	LikesCount int64     `json:"likes_count,omitempty"`
}

func (commentData *Comment) IsValid() bool {
	return 0 < len(commentData.Text) && len(commentData.Text) < 401
}

func (comment *Comment) Sanitize() {
	comment.Text = html.EscapeString(comment.Text)
}
