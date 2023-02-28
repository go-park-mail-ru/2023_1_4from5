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
