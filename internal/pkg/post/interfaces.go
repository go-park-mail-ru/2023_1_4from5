package post

//go:generate mockgen -source=interfaces.go -destination=./mocks/post_mock.go -package=mock

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

type PostUsecase interface {
	CreatePost(ctx context.Context, postData models.PostCreationData) error
	GetPost(ctx context.Context, postID, userID uuid.UUID) (models.PostWithComments, error)
	DeletePost(ctx context.Context, postID uuid.UUID) error
	IsPostOwner(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error)
	AddLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error)
	RemoveLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error)
	IsCreator(ctx context.Context, userID uuid.UUID, creatorID uuid.UUID) (bool, error)
	EditPost(ctx context.Context, postData models.PostEditData) error
	IsPostAvailable(ctx context.Context, postID, userID uuid.UUID) error
}
type PostRepo interface {
	CreatePost(ctx context.Context, postData models.PostCreationData) error
	GetPost(ctx context.Context, postID, userID uuid.UUID) (models.Post, error)
	DeletePost(ctx context.Context, postID uuid.UUID) error
	GetSubsByID(ctx context.Context, subsIDs ...uuid.UUID) ([]models.Subscription, error)
	IsPostOwner(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error)
	AddLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error)
	RemoveLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error)
	IsCreator(ctx context.Context, userID uuid.UUID, creatorID uuid.UUID) (bool, error)
	IsPostAvailable(ctx context.Context, userID, postID uuid.UUID) error
	EditPost(ctx context.Context, postData models.PostEditData) error
	GetComments(ctx context.Context, postID, userID uuid.UUID) ([]models.Comment, error)
}
