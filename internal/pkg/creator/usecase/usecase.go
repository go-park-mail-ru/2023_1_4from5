package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreatorUsecase struct {
	repo   creator.CreatorRepo
	logger *zap.SugaredLogger
}

func NewCreatorUsecase(repo creator.CreatorRepo, logger *zap.SugaredLogger) *CreatorUsecase {
	return &CreatorUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *CreatorUsecase) GetPage(ctx context.Context, userID, creatorID uuid.UUID) (models.CreatorPage, error) {
	return uc.repo.GetPage(ctx, userID, creatorID)
}

func (uc *CreatorUsecase) CreatorNotificationInfo(ctx context.Context, creatorID uuid.UUID) (models.NotificationCreatorInfo, error) {
	return uc.repo.CreatorNotificationInfo(ctx, creatorID)
}

func (uc *CreatorUsecase) CreateAim(ctx context.Context, aimInfo models.Aim) error {
	return uc.repo.CreateAim(ctx, aimInfo)
}

func (uc *CreatorUsecase) UpdateCreatorData(ctx context.Context, updateData models.UpdateCreatorInfo) error {
	return uc.repo.UpdateCreatorData(ctx, updateData)
}

func (uc *CreatorUsecase) CheckIfCreator(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	return uc.repo.CheckIfCreator(ctx, userID)
}

func (uc *CreatorUsecase) FindCreators(ctx context.Context, keyword string) ([]models.Creator, error) {
	return uc.repo.FindCreators(ctx, keyword)
}

func (uc *CreatorUsecase) GetAllCreators(ctx context.Context) ([]models.Creator, error) {
	return uc.repo.GetAllCreators(ctx)
}

func (uc *CreatorUsecase) GetFeed(ctx context.Context, userID uuid.UUID) ([]models.Post, error) {
	return uc.repo.GetFeed(ctx, userID)
}

func (uc *CreatorUsecase) UpdateProfilePhoto(ctx context.Context, creatorId uuid.UUID) (uuid.UUID, error) {
	path := uuid.New()
	err := uc.repo.UpdateProfilePhoto(ctx, creatorId, path)
	if err != nil {
		return uuid.Nil, models.InternalError
	}
	return path, nil
}

func (uc *CreatorUsecase) UpdateCoverPhoto(ctx context.Context, creatorId uuid.UUID) (uuid.UUID, error) {
	path := uuid.New()
	err := uc.repo.UpdateCoverPhoto(ctx, creatorId, path)
	if err != nil {
		return uuid.Nil, models.InternalError
	}
	return path, nil
}

func (uc *CreatorUsecase) StatisticsFirstDate(ctx context.Context, creatorID uuid.UUID) (string, error) {
	return uc.repo.StatisticsFirstDate(ctx, creatorID)
}

func (uc *CreatorUsecase) DeleteCoverPhoto(ctx context.Context, creatorId uuid.UUID) error {
	return uc.repo.DeleteCoverPhoto(ctx, creatorId)
}

func (uc *CreatorUsecase) DeleteProfilePhoto(ctx context.Context, creatorId uuid.UUID) error {
	return uc.repo.DeleteProfilePhoto(ctx, creatorId)
}
func (uc *CreatorUsecase) Statistics(ctx context.Context, statsInput models.StatisticsDates) (models.Statistics, error) {
	return uc.repo.Statistics(ctx, statsInput)
}

func (uc *CreatorUsecase) GetCreatorBalance(ctx context.Context, creatorID uuid.UUID) (float32, error) {
	return uc.repo.GetCreatorBalance(ctx, creatorID)
}

func (uc *CreatorUsecase) UpdateBalance(ctx context.Context, transfer models.CreatorTransfer) (float32, error) {
	return uc.repo.UpdateBalance(ctx, transfer)
}
