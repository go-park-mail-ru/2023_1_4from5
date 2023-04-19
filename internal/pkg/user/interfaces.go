package user

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/user_mock.go -package=mock

type UserUsecase interface {
	GetProfile(ctx context.Context, details models.AccessDetails) (models.UserProfile, error)
	GetHomePage(ctx context.Context, details models.AccessDetails) (models.UserHomePage, error)
	UpdatePhoto(ctx context.Context, details models.AccessDetails) (uuid.UUID, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	UpdateProfileInfo(ctx context.Context, profileInfo models.UpdateProfileInfo, id uuid.UUID) error
	Donate(ctx context.Context, donateInfo models.Donate, userID uuid.UUID) (int, error)
	CheckIfCreator(ctx context.Context, userId uuid.UUID) (uuid.UUID, bool, error)
	BecomeCreator(ctx context.Context, creatorInfo models.BecameCreatorInfo, userId uuid.UUID) (uuid.UUID, error)
	Follow(ctx context.Context, userId, creatorId uuid.UUID) error
	Subscribe(ctx context.Context, subscription models.SubscriptionDetails) error
}

type UserRepo interface {
	GetUserProfile(ctx context.Context, id uuid.UUID) (models.UserProfile, error)
	GetHomePage(ctx context.Context, id uuid.UUID) (models.UserHomePage, error)
	UpdateProfilePhoto(ctx context.Context, userID uuid.UUID, path uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	UpdateProfileInfo(ctx context.Context, profileInfo models.UpdateProfileInfo, id uuid.UUID) error
	Donate(ctx context.Context, donateInfo models.Donate, userID uuid.UUID) (int, error)
	CheckIfCreator(ctx context.Context, userId uuid.UUID) (uuid.UUID, bool, error)
	BecomeCreator(ctx context.Context, creatorInfo models.BecameCreatorInfo, userId uuid.UUID) (uuid.UUID, error)
	Follow(ctx context.Context, userId, creatorId uuid.UUID) error
	CheckIfFollow(ctx context.Context, userId, creatorId uuid.UUID) (bool, error)
	Subscribe(ctx context.Context, subscription models.SubscriptionDetails) error
}
