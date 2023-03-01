package usecase

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
)

type CreatorUsecase struct {
	repo creator.CreatorRepo
}

func NewCreatorUsecase(repo creator.CreatorRepo) *CreatorUsecase {
	return &CreatorUsecase{repo: repo}
}

func (uc *CreatorUsecase) GetPage(details models.AccessDetails, creatorInfo models.Creator) (models.CreatorPage, error) {
	userId := details.Id
	creatorId := creatorInfo.Id

	creatorPage, err := uc.repo.GetPage(userId, creatorId)
	if err != nil {
		return models.CreatorPage{}, err
	}

	return creatorPage, nil
}
