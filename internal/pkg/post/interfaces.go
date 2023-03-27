package post

//go:generate mockgen -source=interfaces.go -destination=./mocks/post_mock.go -package=mock

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

type PostUsecase interface {
	CreatePost(postData models.PostCreationData) error
	DeletePost(postID uuid.UUID) error
	IsPostOwner(userId uuid.UUID, postId uuid.UUID) (bool, error)
	AddLike(userID uuid.UUID, postID uuid.UUID) (models.Like, error)
	RemoveLike(userID uuid.UUID, postID uuid.UUID) (models.Like, error)
}
type PostRepo interface {
	CreatePost(postData models.PostCreationData) error
	DeletePost(postID uuid.UUID) error
	IsPostOwner(userId uuid.UUID, postId uuid.UUID) (bool, error)
	AddLike(userID uuid.UUID, postID uuid.UUID) (models.Like, error)
	RemoveLike(userID uuid.UUID, postID uuid.UUID) (models.Like, error)
}
