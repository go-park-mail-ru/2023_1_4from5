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
