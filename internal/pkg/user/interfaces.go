package user

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/user_mock.go -package=mock

type UserUsecase interface {
	GetProfile(details models.AccessDetails) (models.UserProfile, error)
	GetHomePage(details models.AccessDetails) (models.UserHomePage, error)
}

type UserRepo interface {
	GetUserProfile(id uuid.UUID) (models.UserProfile, error)
	GetHomePage(id uuid.UUID) (models.UserHomePage, error)
}
