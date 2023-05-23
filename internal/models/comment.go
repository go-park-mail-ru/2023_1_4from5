package models

import (
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/google/uuid"
	"html"
	"time"
)

// easyjson -all ./internal/models/comment.go

type Comment struct {
	CommentID  uuid.UUID `json:"comment_id,omitempty"`
	UserID     uuid.UUID `json:"user_id,omitempty"`
	Username   string    `json:"username,omitempty"`
	UserPhoto  uuid.UUID `json:"user_photo,omitempty"`
	PostID     uuid.UUID `json:"post_id,omitempty"`
	Text       string    `json:"text,omitempty"`
	Creation   time.Time `json:"creation,omitempty"`
	LikesCount int64     `json:"likes_count,omitempty"`
	IsLiked    bool      `json:"is_liked"`
	IsOwner    bool      `json:"is_owner"`
}

func (comment *Comment) IsValid() bool {
	return 0 < len(comment.Text) && len(comment.Text) < 401
}

func (comment *Comment) Sanitize() {
	comment.Text = html.EscapeString(comment.Text)
	comment.Username = html.EscapeString(comment.Username)
}

func (comment *Comment) CommentToModel(com *generatedCreator.Comment) error {
	commentID, err := uuid.Parse(com.Id)
	if err != nil {
		return err
	}
	postID, err := uuid.Parse(com.PostID)
	if err != nil {
		return err
	}
	userPhoto, err := uuid.Parse(com.UserPhoto)
	if err != nil {
		return err
	}

	userID, err := uuid.Parse(com.UserId)
	if err != nil {
		return err
	}

	creation, err := time.Parse(time.RFC3339, com.Creation)
	if err != nil {
		return err
	}

	comment.CommentID = commentID
	comment.UserPhoto = userPhoto
	comment.Creation = creation
	comment.Text = com.Text
	comment.LikesCount = com.LikesCount
	comment.PostID = postID
	comment.UserID = userID
	comment.Username = com.Username
	comment.IsLiked = com.IsLiked
	comment.IsOwner = com.IsOwner
	return nil
}
