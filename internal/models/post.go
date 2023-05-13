package models

// easyjson -all ./internal/models/post.go

import (
	"fmt"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/google/uuid"
	"html"
	"time"
)

type Post struct {
	Id            uuid.UUID      `json:"id"`
	Creator       uuid.UUID      `json:"creator"`
	CreatorPhoto  uuid.UUID      `json:"creator_photo,omitempty"`
	CreatorName   string         `json:"creator_name,omitempty"`
	Creation      time.Time      `json:"creation_date"`
	LikesCount    int64          `json:"likes_count"`
	CommentsCount int64          `json:"comments_count"`
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

func (post *Post) PostToModel(postInfo *generatedCreator.Post) error {
	postID, err := uuid.Parse(postInfo.Id)
	if err != nil {
		return err
	}
	creatorID, err := uuid.Parse(postInfo.CreatorID)
	if err != nil {
		return err
	}
	creatorPhoto, err := uuid.Parse(postInfo.CreatorPhoto)
	if err != nil {
		return err
	}

	reg, err := time.Parse("2006-01-02 15:04:05 -0700 -0700", postInfo.Creation)
	if err != nil {
		fmt.Println("date")
		return err
	}

	post.Id = postID
	post.Creator = creatorID
	post.CreatorPhoto = creatorPhoto
	post.CreatorName = postInfo.CreatorName
	post.Creation = reg
	post.LikesCount = postInfo.LikesCount
	post.CommentsCount = postInfo.CommentsCount
	post.Title = postInfo.Title
	post.Text = postInfo.Text
	post.IsAvailable = postInfo.IsAvailable
	post.IsLiked = postInfo.IsLiked

	for _, sub := range postInfo.Subscriptions {
		var subscription Subscription
		err = subscription.ProtoSubscriptionToModel(sub)
		if err != nil {
			return err
		}
		post.Subscriptions = append(post.Subscriptions, subscription)
	}

	for _, attach := range postInfo.PostAttachments {
		var attachment Attachment
		err = attachment.AttachToModel(attach)
		if err != nil {
			return err
		}
		post.Attachments = append(post.Attachments, attachment)
	}
	return nil
}
