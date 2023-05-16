package models

import (
	"fmt"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/google/uuid"
	"html"
	"time"
)

// easyjson -all ./internal/models/comment.go

type Comment struct {
	CommentID  uuid.UUID `json:"comment_id,omitempty"`
	UserID     uuid.UUID `json:"user_id,omitempty"`
	UserPhoto  uuid.UUID `json:"user_photo,omitempty"`
	PostID     uuid.UUID `json:"post_id,omitempty"`
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

	creation, err := time.Parse("2006-01-02 15:04:05 -0700 -0700", com.Creation)
	if err != nil {
		fmt.Println("date")
		return err
	}

	comment.CommentID = commentID
	comment.UserPhoto = userPhoto
	comment.Creation = creation
	comment.Text = com.Text
	comment.LikesCount = com.LikesCount
	comment.PostID = postID
	comment.UserID = userID
	return nil
}
