package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	comment "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/comment"
	"go.uber.org/zap"
)

type CommentUsecase struct {
	repo   comment.CommentRepo
	logger *zap.SugaredLogger
}

func NewCommentUsecase(repo comment.CommentRepo, logger *zap.SugaredLogger) *CommentUsecase {
	return &CommentUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *CommentUsecase) CreateComment(ctx context.Context, commentInfo models.Comment) error {
	return uc.repo.CreateComment(ctx, commentInfo)
}

func (uc *CommentUsecase) DeleteComment(ctx context.Context, commentInfo models.Comment) error {
	return uc.repo.DeleteComment(ctx, commentInfo)
}

func (uc *CommentUsecase) EditComment(ctx context.Context, commentInfo models.Comment) error {
	return uc.repo.EditComment(ctx, commentInfo)
}

func (uc *CommentUsecase) AddLike(ctx context.Context, commentInfo models.Comment) (int64, error) {
	return uc.repo.AddLike(ctx, commentInfo)
}

func (uc *CommentUsecase) RemoveLike(ctx context.Context, commentInfo models.Comment) (int64, error) {
	return uc.repo.RemoveLike(ctx, commentInfo)
}

func (uc *CommentUsecase) IsCommentOwner(ctx context.Context, commentInfo models.Comment) (bool, error) {
	return uc.repo.IsCommentOwner(ctx, commentInfo)
}
