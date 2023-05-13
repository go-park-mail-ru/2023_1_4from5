package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"go.uber.org/zap"
)

const (
	CreateComment    = `INSERT INTO "comment"(comment_id, post_id, user_id, comment_text) VALUES ($1, $2, $3, $4);`
	IncCommentsCount = `UPDATE "post" SET comments_count = post.comments_count + 1 WHERE post_id = $1 RETURNING comments_count;`
)

type CommentRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewCommentRepo(db *sql.DB, logger *zap.SugaredLogger) *CommentRepo {
	return &CommentRepo{
		db:     db,
		logger: logger,
	}
}

func (r *CommentRepo) CreateComment(ctx context.Context, commentInfo models.Comment) error {
	row := r.db.QueryRowContext(ctx, CreateComment, commentInfo.CommentID, commentInfo.PostID, commentInfo.UserID, commentInfo.Text)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.InternalError
	}

	row = r.db.QueryRowContext(ctx, IncCommentsCount, commentInfo.PostID)
	var commCountTmp int
	if err := row.Scan(&commCountTmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.WrongData
	}

	return nil
}
