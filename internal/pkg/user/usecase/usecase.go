package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserUsecase struct {
	repo   user.UserRepo
	logger *zap.SugaredLogger
}

func NewUserUsecase(repo user.UserRepo, logger *zap.SugaredLogger) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *UserUsecase) CheckIfCreator(ctx context.Context, userId uuid.UUID) (bool, error) {
	return uc.repo.CheckIfCreator(ctx, userId)
}

func (uc *UserUsecase) GetProfile(ctx context.Context, details models.AccessDetails) (models.UserProfile, error) {
	userId := details.Id
	return uc.repo.GetUserProfile(ctx, userId)
}

func (uc *UserUsecase) GetHomePage(ctx context.Context, details models.AccessDetails) (models.UserHomePage, error) {
	userId := details.Id
	return uc.repo.GetHomePage(ctx, userId)
}

func (uc *UserUsecase) UpdatePhoto(ctx context.Context, details models.AccessDetails) (uuid.UUID, error) {
	path := uuid.New()
	err := uc.repo.UpdateProfilePhoto(ctx, details.Id, path)
	if err != nil {
		return uuid.Nil, models.InternalError
	}
	return path, nil
}

func (uc *UserUsecase) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	return uc.repo.UpdatePassword(ctx, id, password)
}

func (uc *UserUsecase) UpdateProfileInfo(ctx context.Context, profileInfo models.UpdateProfileInfo, id uuid.UUID) error {
	return uc.repo.UpdateProfileInfo(ctx, profileInfo, id)
}

func (uc *UserUsecase) Donate(ctx context.Context, donateInfo models.Donate, userID uuid.UUID) (int, error) {
	return uc.repo.Donate(ctx, donateInfo, userID)
}

func (uc *UserUsecase) BecomeCreator(ctx context.Context, creatorInfo models.BecameCreatorInfo, userId uuid.UUID) (uuid.UUID, error) {
	return uc.repo.BecomeCreator(ctx, creatorInfo, userId)
}
