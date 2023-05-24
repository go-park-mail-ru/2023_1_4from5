package user

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/user_mock.go -package=mock

type UserUsecase interface {
	GetProfile(ctx context.Context, userId uuid.UUID) (models.UserProfile, error)
	UpdatePhoto(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
	DeletePhoto(ctx context.Context, userId uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	UpdateProfileInfo(ctx context.Context, profileInfo models.UpdateProfileInfo, id uuid.UUID) error
	Donate(ctx context.Context, donateInfo models.Donate) (float32, error)
	CheckIfCreator(ctx context.Context, userId uuid.UUID) (uuid.UUID, bool, error)
	BecomeCreator(ctx context.Context, creatorInfo models.BecameCreatorInfo, userId uuid.UUID) (uuid.UUID, error)
	Follow(ctx context.Context, userId, creatorId uuid.UUID) error
	Subscribe(ctx context.Context, paymentInfo uuid.UUID, money float32) (models.NotificationSubInfo, error)
	Unfollow(ctx context.Context, userId, creatorId uuid.UUID) error
	UserSubscriptions(ctx context.Context, userId uuid.UUID) ([]models.Subscription, error)
	UserFollows(ctx context.Context, userId uuid.UUID) ([]models.Follow, error)
	AddPaymentInfo(ctx context.Context, subscription models.SubscriptionDetails) error
}

type UserRepo interface {
	GetUserProfile(ctx context.Context, id uuid.UUID) (models.UserProfile, error)
	UpdateProfilePhoto(ctx context.Context, userID uuid.UUID, path uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	UpdateProfileInfo(ctx context.Context, profileInfo models.UpdateProfileInfo, id uuid.UUID) error
	Donate(ctx context.Context, donateInfo models.Donate) (float32, error)
	CheckIfCreator(ctx context.Context, userId uuid.UUID) (uuid.UUID, bool, error)
	BecomeCreator(ctx context.Context, creatorInfo models.BecameCreatorInfo, userId uuid.UUID) (uuid.UUID, error)
	Follow(ctx context.Context, userId, creatorId uuid.UUID) error
	CheckIfFollow(ctx context.Context, userId, creatorId uuid.UUID) (bool, error)
	Subscribe(ctx context.Context, subscription models.SubscriptionDetails) (models.NotificationSubInfo, error)
	Unfollow(ctx context.Context, userId, creatorId uuid.UUID) error
	UserSubscriptions(ctx context.Context, userId uuid.UUID) ([]models.Subscription, error)
	DeletePhoto(ctx context.Context, userId uuid.UUID) error
	UserFollows(ctx context.Context, userId uuid.UUID) ([]models.Follow, error)
	AddPaymentInfo(ctx context.Context, subscription models.SubscriptionDetails) error
	CheckPaymentInfo(ctx context.Context, paymentInfo uuid.UUID) (models.SubscriptionDetails, error)
	UpdatePaymentInfo(ctx context.Context, money float32, paymentInfo uuid.UUID) error
	GetCreatorID(ctx context.Context, subscriptionID uuid.UUID) (uuid.UUID, error)
}
