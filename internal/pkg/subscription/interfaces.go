package subscription

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/post_mock.go -package=mock

type SubscriptionUsecase interface {
	CreateSubscription(ctx context.Context, subscriptionInfo models.Subscription) error
	DeleteSubscription(ctx context.Context, subscriptionID, creatorID uuid.UUID) error
	EditSubscription(ctx context.Context, subscriptionInfo models.Subscription) error
}
type SubscriptionRepo interface {
	CreateSubscription(ctx context.Context, subscriptionInfo models.Subscription) error
	DeleteSubscription(ctx context.Context, subscriptionID, creatorID uuid.UUID) error
	EditSubscription(ctx context.Context, subscriptionInfo models.Subscription) error
}
