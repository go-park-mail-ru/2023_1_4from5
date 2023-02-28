package usecase

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
)

type CreatorUsecase struct {
	repo creator.CreatorRepo
}

func NewUserUsecase(repo creator.CreatorRepo) *CreatorUsecase {
	return &CreatorUsecase{repo: repo}
}

//func (uc *CreatorUsecase) GetPage(details models.AccessDetails) (models.UserProfile, error) {

//}
