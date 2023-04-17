package subscription

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/post_mock.go -package=mock

type SubscriptionUsecase interface {
	CreateSubscription(ctx context.Context, subscriptionInfo models.Subscription) error
}
type SubscriptionRepo interface {
	CreateSubscription(ctx context.Context, subscriptionInfo models.Subscription) error
}
