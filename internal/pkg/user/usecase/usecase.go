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

func (uc *UserUsecase) CheckIfCreator(ctx context.Context, userId uuid.UUID) (uuid.UUID, bool, error) {
	return uc.repo.CheckIfCreator(ctx, userId)
}

func (uc *UserUsecase) Follow(ctx context.Context, userId, creatorId uuid.UUID) error {
	if isFollowing, err := uc.repo.CheckIfFollow(ctx, userId, creatorId); err == models.InternalError {
		return err
	} else if isFollowing == true {
		return models.WrongData
	}
	return uc.repo.Follow(ctx, userId, creatorId)
}

func (uc *UserUsecase) Unfollow(ctx context.Context, userId, creatorId uuid.UUID) error {
	if isFollowing, err := uc.repo.CheckIfFollow(ctx, userId, creatorId); err == models.InternalError {
		return err
	} else if isFollowing == false {
		return models.WrongData
	}
	return uc.repo.Unfollow(ctx, userId, creatorId)
}

func (uc *UserUsecase) Subscribe(ctx context.Context, subscription models.SubscriptionDetails) error {
	if isFollowing, err := uc.repo.CheckIfFollow(ctx, subscription.UserID, subscription.CreatorId); err == models.InternalError {
		return err
	} else if isFollowing == false {
		err = uc.repo.Follow(ctx, subscription.UserID, subscription.CreatorId)
		if err == models.InternalError {
			return err
		}
	}
	return uc.repo.Subscribe(ctx, subscription)
}

func (uc *UserUsecase) GetProfile(ctx context.Context, userId uuid.UUID) (models.UserProfile, error) {
	return uc.repo.GetUserProfile(ctx, userId)
}

func (uc *UserUsecase) UpdatePhoto(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	path := uuid.New()
	err := uc.repo.UpdateProfilePhoto(ctx, userId, path)
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

func (uc *UserUsecase) Donate(ctx context.Context, donateInfo models.Donate, userID uuid.UUID) (int64, error) {
	return uc.repo.Donate(ctx, donateInfo, userID)
}

func (uc *UserUsecase) BecomeCreator(ctx context.Context, creatorInfo models.BecameCreatorInfo, userId uuid.UUID) (uuid.UUID, error) {
	return uc.repo.BecomeCreator(ctx, creatorInfo, userId)
}

func (uc *UserUsecase) UserSubscriptions(ctx context.Context, userId uuid.UUID) ([]models.Subscription, error) {
	return uc.repo.UserSubscriptions(ctx, userId)
}

func (uc *UserUsecase) DeletePhoto(ctx context.Context, userId uuid.UUID) error {
	return uc.repo.DeletePhoto(ctx, userId)
}
