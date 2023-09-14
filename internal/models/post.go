package models

// easyjson -all ./internal/models/post.go

import (
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/google/uuid"
	"html"
	"time"
	"unicode"
)

const (
	postTitleMaxLength = 80
	postTextMaxLength  = 4000
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

type PostWithComments struct {
	Post     Post      `json:"post"`
	Comments []Comment `json:"comments"`
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

func (postCreationData PostCreationData) IsValid() error {
	if len(postCreationData.Title) == 0 || len([]rune(postCreationData.Title)) > postTitleMaxLength {
		return WrongPostTitleLength
	}
	if len([]rune(postCreationData.Text)) > postTextMaxLength {
		return WrongPostTextLength
	}
	for _, c := range postCreationData.Title {
		if !unicode.IsLetter(c) && !(c >= 32 && c <= 126) && c != 10 && c != 13 {
			return WrongPostTitleSymbols
		}
	}
	for _, c := range postCreationData.Text {
		if !unicode.IsLetter(c) && !(c >= 32 && c <= 126) && c != 10 && c != 13 {
			return WrongPostTextSymbols
		}
	}
	return nil
}

func (post *Post) Sanitize() {
	post.Title = html.EscapeString(post.Title)
	post.Text = html.EscapeString(post.Text)
	for i := range post.Subscriptions {
		post.Subscriptions[i].Sanitize()
	}
}

func (postWithComments *PostWithComments) Sanitize() {
	postWithComments.Post.Sanitize()
	for i := range postWithComments.Comments {
		postWithComments.Comments[i].Sanitize()
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

	reg, err := time.Parse(time.RFC3339, postInfo.Creation)
	if err != nil {
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

func (postWithComments *PostWithComments) PostWithCommentsToModel(postInfo *generatedCreator.PostWithComments) error {
	err := postWithComments.Post.PostToModel(postInfo.Post)
	if err != nil {
		return err
	}
	for _, com := range postInfo.Comments {
		var comment Comment
		err = comment.CommentToModel(com)
		if err != nil {
			return err
		}
		postWithComments.Comments = append(postWithComments.Comments, comment)
	}

	return nil
}
