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

func (u *PostUsecase) GetPost(ctx context.Context, postID, userID uuid.UUID) (models.PostWithComments, error) {
	var postWithComments models.PostWithComments
	var isAvailable bool
	err := u.repo.IsPostAvailable(ctx, userID, postID)
	if err == models.InternalError {
		return models.PostWithComments{}, err
	}
	if err == nil {
		isAvailable = true
	}
	postWithComments.Post, err = u.repo.GetPost(ctx, postID, userID)
	if err != nil {
		return models.PostWithComments{}, err
	}
	postWithComments.Comments, err = u.repo.GetComments(ctx, postID, userID)
	if err != nil {
		return models.PostWithComments{}, err
	}
	if !isAvailable {
		postWithComments.Post.Attachments = nil
		postWithComments.Post.Text = ""
		postWithComments.Comments = nil
	}

	postWithComments.Post.IsAvailable = isAvailable
	return postWithComments, nil
}

func (u *PostUsecase) DeletePost(ctx context.Context, postID uuid.UUID) error {
	return u.repo.DeletePost(ctx, postID)
}
func (u *PostUsecase) IsPostOwner(ctx context.Context, userId, postId uuid.UUID) (bool, error) {
	return u.repo.IsPostOwner(ctx, userId, postId)
}

func (u *PostUsecase) AddLike(ctx context.Context, userID, postID uuid.UUID) (models.Like, error) {
	return u.repo.AddLike(ctx, userID, postID)
}

func (u *PostUsecase) RemoveLike(ctx context.Context, userID, postID uuid.UUID) (models.Like, error) {
	return u.repo.RemoveLike(ctx, userID, postID)
}
func (u *PostUsecase) EditPost(ctx context.Context, postData models.PostEditData) error {
	return u.repo.EditPost(ctx, postData)
}

func (u *PostUsecase) IsPostAvailable(ctx context.Context, postID, userID uuid.UUID) error {
	return u.repo.IsPostAvailable(ctx, postID, userID)
}
