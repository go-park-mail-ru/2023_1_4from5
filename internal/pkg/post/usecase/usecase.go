package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostUsecase struct {
	repo   post.PostRepo
	logger *zap.SugaredLogger
}

func NewPostUsecase(repo post.PostRepo, logger *zap.SugaredLogger) *PostUsecase {
	return &PostUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *PostUsecase) CreatePost(ctx context.Context, postData models.PostCreationData) error {
	return u.repo.CreatePost(ctx, postData)
}
func (u *PostUsecase) IsCreator(ctx context.Context, userID, creatorID uuid.UUID) (bool, error) {
	return u.repo.IsCreator(ctx, userID, creatorID)
}

func (u *PostUsecase) GetPost(ctx context.Context, postID, userID uuid.UUID) (models.Post, error) {
	var isAvailable bool
	err := u.repo.IsPostAvailable(ctx, userID, postID)
	if err == models.InternalError {
		return models.Post{}, err
	}
	if err == nil {
		isAvailable = true
	}
	post, err := u.repo.GetPost(ctx, postID, userID)
	if err != nil {
		return models.Post{}, models.InternalError
	}

	if !isAvailable {
		post.Attachments = nil
		post.Text = ""
	}

	post.IsAvailable = isAvailable
	return post, nil
}

func (u *PostUsecase) DeletePost(ctx context.Context, postID uuid.UUID) error {
	return u.repo.DeletePost(ctx, postID)
}
func (u *PostUsecase) IsPostOwner(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
	return u.repo.IsPostOwner(ctx, userId, postId)
}

func (u *PostUsecase) AddLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
	return u.repo.AddLike(ctx, userID, postID)
}

func (u *PostUsecase) RemoveLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
	return u.repo.RemoveLike(ctx, userID, postID)
}
