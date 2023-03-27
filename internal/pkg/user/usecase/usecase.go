package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
	"github.com/google/uuid"
)

type UserUsecase struct {
	repo user.UserRepo
}

func NewUserUsecase(repo user.UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) GetProfile(ctx context.Context, details models.AccessDetails) (models.UserProfile, error) {
	userId := details.Id
	userProfile, err := uc.repo.GetUserProfile(ctx, userId)
	if err != nil {
		return models.UserProfile{}, err
	}

	return userProfile, nil
}

func (uc *UserUsecase) GetHomePage(ctx context.Context, details models.AccessDetails) (models.UserHomePage, error) {
	userId := details.Id
	homePage, err := uc.repo.GetHomePage(ctx, userId)
	if err != nil {
		return models.UserHomePage{}, err
	}

	return homePage, nil
}

func (uc *UserUsecase) UpdatePhoto(ctx context.Context, details models.AccessDetails) (uuid.UUID, error) {
	path := uuid.New()
	err := uc.repo.UpdateProfilePhoto(ctx, details.Id, path)
	if err != nil {
		return uuid.Nil, models.InternalError
	}
	return path, nil
}
