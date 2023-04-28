package creator

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/creator_mock.go -package=mock

type CreatorUsecase interface {
	GetPage(ctx context.Context, userID, creatorID uuid.UUID) (models.CreatorPage, error)
	CreateAim(ctx context.Context, aimInfo models.Aim) error
	GetAllCreators(ctx context.Context) ([]models.Creator, error)
	FindCreators(ctx context.Context, keyword string) ([]models.Creator, error)
	UpdateCreatorData(ctx context.Context, updateData models.UpdateCreatorInfo) error
	CheckIfCreator(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

type CreatorRepo interface {
	GetCreatorSubs(ctx context.Context, creatorID uuid.UUID) ([]models.Subscription, error)
	GetPage(ctx context.Context, userID, creatorID uuid.UUID) (models.CreatorPage, error)
	CreateAim(ctx context.Context, aimInfo models.Aim) error
	GetAllCreators(ctx context.Context) ([]models.Creator, error)
	FindCreators(ctx context.Context, keyword string) ([]models.Creator, error)
	UpdateCreatorData(ctx context.Context, updateData models.UpdateCreatorInfo) error
	CheckIfCreator(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}
