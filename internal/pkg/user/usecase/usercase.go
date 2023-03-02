package usecase

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
)

type UserUsecase struct {
	repo user.UserRepo
}

func NewUserUsecase(repo user.UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) GetProfile(details models.AccessDetails) (models.UserProfile, error) {
	userId := details.Id
	userProfile, err := uc.repo.GetUserProfile(userId)
	if err != nil {
		return models.UserProfile{}, err
	}

	return userProfile, nil
}

func (uc *UserUsecase) GetHomePage(details models.AccessDetails) (models.UserHomePage, error) {
	userId := details.Id
	homePage, err := uc.repo.GetHomePage(userId)
	if err != nil {
		return models.UserHomePage{}, err
	}

	return homePage, nil
}
