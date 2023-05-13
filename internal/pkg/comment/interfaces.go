package subscription

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/post_mock.go -package=mock

type CommentUsecase interface {
	CreateComment(ctx context.Context, commentData models.Comment) error
	//DeleteComment(ctx context.Context, commentID uuid.UUID) error
	//EditComment(ctx context.Context, subscriptionInfo models.Subscription) error
	//AddLike(ctx context.Context, subscriptionInfo models.Subscription) error
	//RemoveLike(ctx context.Context, subscriptionInfo models.Subscription) error
}
type CommentRepo interface {
	CreateComment(ctx context.Context, commentData models.Comment) error
	//DeleteComment(ctx context.Context, subscriptionID, creatorID uuid.UUID) error
	//EditComment(ctx context.Context, subscriptionInfo models.Subscription) error
	//AddLike(ctx context.Context, subscriptionInfo models.Subscription) error
	//RemoveLike(ctx context.Context, subscriptionInfo models.Subscription) error
}
