package subscription

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/post_mock.go -package=mock

type CommentUsecase interface {
	CreateComment(ctx context.Context, commentData models.Comment) error
	DeleteComment(ctx context.Context, commentData models.Comment) error
	EditComment(ctx context.Context, commentInfo models.Comment) error
	IsCommentOwner(ctx context.Context, commentInfo models.Comment) (bool, error)
	AddLike(ctx context.Context, commentInfo models.Comment) (int64, error)
	RemoveLike(ctx context.Context, commentInfo models.Comment) (int64, error)
}
type CommentRepo interface {
	CreateComment(ctx context.Context, commentData models.Comment) error
	DeleteComment(ctx context.Context, commentData models.Comment) error
	EditComment(ctx context.Context, commentInfo models.Comment) error
	IsCommentOwner(ctx context.Context, commentInfo models.Comment) (bool, error)
	AddLike(ctx context.Context, commentInfo models.Comment) (int64, error)
	RemoveLike(ctx context.Context, commentInfo models.Comment) (int64, error)
}
