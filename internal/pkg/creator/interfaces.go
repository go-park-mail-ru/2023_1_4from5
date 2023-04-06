package creator

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/creator_mock.go -package=mock

type CreatorUsecase interface {
	GetPage(ctx context.Context, details *models.AccessDetails, creatorUUID string) (models.CreatorPage, error)
	CreateAim(ctx context.Context, aimInfo models.Aim) error
}

type CreatorRepo interface {
	GetPage(ctx context.Context, userId uuid.UUID, creatorId uuid.UUID) (models.CreatorPage, error)
	CreateAim(ctx context.Context, aimInfo models.Aim) error
}
