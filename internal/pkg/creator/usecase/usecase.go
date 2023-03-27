package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	"github.com/google/uuid"
)

type CreatorUsecase struct {
	repo creator.CreatorRepo
}

func NewCreatorUsecase(repo creator.CreatorRepo) *CreatorUsecase {
	return &CreatorUsecase{repo: repo}
}

func (uc *CreatorUsecase) GetPage(ctx context.Context, details *models.AccessDetails, creatorUUID string) (models.CreatorPage, error) {
	userId := details.Id
	creatorId, err := uuid.Parse(creatorUUID)
	if err != nil {
		return models.CreatorPage{}, models.WrongData
	}
	creatorPage, err := uc.repo.GetPage(ctx, userId, creatorId)
	if err != nil {
		return models.CreatorPage{}, err
	}

	return creatorPage, nil
}
