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
	GetFeed(ctx context.Context, userID uuid.UUID) ([]models.Post, error)
	UpdateProfilePhoto(ctx context.Context, creatorId uuid.UUID) (uuid.UUID, error)
	DeleteProfilePhoto(ctx context.Context, creatorId uuid.UUID) error
	UpdateCoverPhoto(ctx context.Context, creatorId uuid.UUID) (uuid.UUID, error)
	DeleteCoverPhoto(ctx context.Context, creatorId uuid.UUID) error
	Statistics(ctx context.Context, statsInput models.StatisticsDates) (models.Statistics, error)
	StatisticsFirstDate(ctx context.Context, creatorID uuid.UUID) (string, error)
	CreatorNotificationInfo(ctx context.Context, creatorID uuid.UUID) (models.NotificationCreatorInfo, error)
	GetCreatorBalance(ctx context.Context, creatorID uuid.UUID) (float32, error)
	UpdateBalance(ctx context.Context, transfer models.CreatorTransfer) (float32, error)
}

type CreatorRepo interface {
	GetCreatorSubs(ctx context.Context, creatorID uuid.UUID) ([]models.Subscription, error)
	GetPage(ctx context.Context, userID, creatorID uuid.UUID) (models.CreatorPage, error)
	CreateAim(ctx context.Context, aimInfo models.Aim) error
	GetAllCreators(ctx context.Context) ([]models.Creator, error)
	FindCreators(ctx context.Context, keyword string) ([]models.Creator, error)
	UpdateCreatorData(ctx context.Context, updateData models.UpdateCreatorInfo) error
	CheckIfCreator(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetFeed(ctx context.Context, userID uuid.UUID) ([]models.Post, error)
	UpdateProfilePhoto(ctx context.Context, creatorId, path uuid.UUID) error
	DeleteProfilePhoto(ctx context.Context, creatorId uuid.UUID) error
	UpdateCoverPhoto(ctx context.Context, creatorId, path uuid.UUID) error
	DeleteCoverPhoto(ctx context.Context, creatorId uuid.UUID) error
	Statistics(ctx context.Context, statsInput models.StatisticsDates) (models.Statistics, error)
	StatisticsFirstDate(ctx context.Context, creatorID uuid.UUID) (string, error)
	CreatorNotificationInfo(ctx context.Context, creatorID uuid.UUID) (models.NotificationCreatorInfo, error)
	GetCreatorBalance(ctx context.Context, creatorID uuid.UUID) (float32, error)
	UpdateBalance(ctx context.Context, transfer models.CreatorTransfer) (float32, error)
}
