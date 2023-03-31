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
}

type UserRepo interface {
	GetUserProfile(ctx context.Context, id uuid.UUID) (models.UserProfile, error)
	GetHomePage(ctx context.Context, id uuid.UUID) (models.UserHomePage, error)
	UpdateProfilePhoto(ctx context.Context, userID uuid.UUID, path uuid.UUID) error
}
