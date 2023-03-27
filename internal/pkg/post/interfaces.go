package post

//go:generate mockgen -source=interfaces.go -destination=./mocks/post_mock.go -package=mock

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

type PostUsecase interface {
	CreatePost(ctx context.Context, postData models.PostCreationData) error
	DeletePost(ctx context.Context, postID uuid.UUID) error
	IsPostOwner(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error)
	AddLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error)
	RemoveLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error)
}
type PostRepo interface {
	CreatePost(ctx context.Context, postData models.PostCreationData) error
	DeletePost(ctx context.Context, postID uuid.UUID) error
	IsPostOwner(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error)
	AddLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error)
	RemoveLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error)
}
