package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/google/uuid"
)

type PostUsecase struct {
	repo post.PostRepo
}

func NewPostUsecase(repo post.PostRepo) *PostUsecase {
	return &PostUsecase{repo: repo}
}

func (u *PostUsecase) CreatePost(ctx context.Context, postData models.PostCreationData) error {
	if err := u.repo.CreatePost(ctx, postData); err != nil {
		return models.InternalError
	}
	return nil
}
func (u *PostUsecase) DeletePost(ctx context.Context, postID uuid.UUID) error {
	if err := u.repo.DeletePost(ctx, postID); err != nil {
		return models.InternalError
	}
	return nil
}
func (u *PostUsecase) IsPostOwner(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
	return u.repo.IsPostOwner(ctx, userId, postId)
}

func (u *PostUsecase) AddLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
	like, err := u.repo.AddLike(ctx, userID, postID)
	if err != nil {
		return models.Like{}, err
	}
	return like, nil
}

func (u *PostUsecase) RemoveLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
	like, err := u.repo.RemoveLike(ctx, userID, postID)
	if err != nil {
		return models.Like{}, err
	}
	return like, nil
}
