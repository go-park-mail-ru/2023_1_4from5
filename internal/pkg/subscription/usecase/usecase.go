package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/subscription"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionUsecase struct {
	repo   subscription.SubscriptionRepo
	logger *zap.SugaredLogger
}

func NewSubscriptionUsecase(repo subscription.SubscriptionRepo, logger *zap.SugaredLogger) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *SubscriptionUsecase) CreateSubscription(ctx context.Context, subscriptionInfo models.Subscription) error {
	return uc.repo.CreateSubscription(ctx, subscriptionInfo)
}

func (uc *SubscriptionUsecase) DeleteSubscription(ctx context.Context, subscriptionID, creatorID uuid.UUID) error {
	return uc.repo.DeleteSubscription(ctx, subscriptionID, creatorID)
}

func (uc *SubscriptionUsecase) EditSubscription(ctx context.Context, subscriptionNewInfo models.Subscription) error {
	return uc.repo.EditSubscription(ctx, subscriptionNewInfo)
}
