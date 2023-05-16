package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	CreateComment      = `INSERT INTO "comment"(comment_id, post_id, user_id, comment_text) VALUES ($1, $2, $3, $4);`
	IncCommentsCount   = `UPDATE "post" SET comments_count = post.comments_count + 1 WHERE post_id = $1 RETURNING comments_count;`
	DecCommentsCount   = `UPDATE "post" SET comments_count = post.comments_count - 1 WHERE post_id = $1;`
	GetUserId          = `SELECT user_id FROM "comment" WHERE comment_id = $1;`
	DeleteCommentLikes = `DELETE FROM "like_comment" WHERE comment_id = $1;`
	DeleteComment      = `DELETE FROM "comment" WHERE comment_id = $1;`
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

func (r *CommentRepo) DeleteComment(ctx context.Context, commentInfo models.Comment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	row, err := tx.QueryContext(ctx, DeleteCommentLikes, commentInfo.CommentID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}
	if err = row.Close(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	row, err = tx.QueryContext(ctx, DecCommentsCount, commentInfo.PostID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}
	if err = row.Close(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	row, err = tx.QueryContext(ctx, DeleteComment, commentInfo.CommentID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}
	if err = row.Close(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	return nil
}

func (r *CommentRepo) IsCommentOwner(ctx context.Context, commentInfo models.Comment) (bool, error) {
	row := r.db.QueryRowContext(ctx, GetUserId, commentInfo.CommentID)

	userIdTmp := uuid.UUID{}

	if err := row.Scan(&userIdTmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return false, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return false, models.WrongData
	}
	if userIdTmp != commentInfo.UserID {
		return false, nil
	}
	return true, nil
}
