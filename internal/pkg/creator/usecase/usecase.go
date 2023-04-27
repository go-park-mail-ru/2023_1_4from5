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

func (uc *CreatorUsecase) CreateAim(ctx context.Context, aimInfo models.Aim) error {
	return uc.repo.CreateAim(ctx, aimInfo)
}

func (uc *CreatorUsecase) FindCreators(ctx context.Context, keyword string) ([]models.Creator, error) {
	return uc.repo.FindCreators(ctx, keyword)
}

func (uc *CreatorUsecase) GetAllCreators(ctx context.Context) ([]models.Creator, error) {
	return uc.repo.GetAllCreators(ctx)
}
