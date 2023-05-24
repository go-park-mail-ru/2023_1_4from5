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
	EditComment        = `UPDATE "comment" SET comment_text = $1 WHERE comment_id = $2;`
	IncCommentsCount   = `UPDATE "post" SET comments_count = post.comments_count + 1 WHERE post_id = $1 RETURNING comments_count;`
	IncLikesCount      = `UPDATE "comment" SET likes_count = comment.likes_count + 1 WHERE comment_id = $1 AND post_id = $2 RETURNING likes_count;`
	DecLikesCount      = `UPDATE "comment" SET likes_count = comment.likes_count - 1 WHERE comment_id = $1 RETURNING likes_count;`
	DecCommentsCount   = `UPDATE "post" SET comments_count = post.comments_count - 1 WHERE post_id = $1;`
	GetUserId          = `SELECT user_id FROM "comment" WHERE comment_id = $1;`
	DeleteCommentLikes = `DELETE FROM "like_comment" WHERE comment_id = $1;`
	DeleteComment      = `DELETE FROM "comment" WHERE comment_id = $1;`
	IsLiked            = `SELECT comment_id FROM "like_comment" WHERE comment_id = $1 AND user_id = $2;`
	AddLike            = `INSERT INTO "like_comment"(comment_id, user_id) VALUES($1, $2);`
	DeleteLike         = `DELETE FROM "like_comment" WHERE comment_id = $1;`
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

func (r *CommentRepo) EditComment(ctx context.Context, commentInfo models.Comment) error {
	row := r.db.QueryRowContext(ctx, EditComment, commentInfo.Text, commentInfo.CommentID)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (r *CommentRepo) DeleteComment(ctx context.Context, commentInfo models.Comment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	_, err = tx.ExecContext(ctx, DeleteCommentLikes, commentInfo.CommentID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	_, err = tx.ExecContext(ctx, DecCommentsCount, commentInfo.PostID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	_, err = tx.ExecContext(ctx, DeleteComment, commentInfo.CommentID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	return nil
}

func (r *CommentRepo) AddLike(ctx context.Context, commentInfo models.Comment) (int64, error) {
	var commentId = uuid.UUID{}
	row := r.db.QueryRowContext(ctx, IsLiked, commentInfo.CommentID, commentInfo.UserID)
	if err := row.Scan(&commentId); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return 0, models.InternalError
	} else if err == nil {
		return 0, models.WrongData
	}

	isCommentOwner, err := r.IsCommentOwner(ctx, commentInfo)
	if err != nil {
		return 0, err
	}
	if isCommentOwner {
		return 0, models.WrongData
	}

	var likesCountTmp int64

	row = r.db.QueryRowContext(ctx, IncLikesCount, commentInfo.CommentID, commentInfo.PostID)

	if err := row.Scan(&likesCountTmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return 0, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return 0, models.WrongData
	}

	row = r.db.QueryRowContext(ctx, AddLike, commentInfo.CommentID, commentInfo.UserID)

	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return 0, models.InternalError
	}

	return likesCountTmp, nil

}

func (r *CommentRepo) RemoveLike(ctx context.Context, commentInfo models.Comment) (int64, error) {
	var commentId = uuid.UUID{}
	row := r.db.QueryRowContext(ctx, IsLiked, commentInfo.CommentID, commentInfo.UserID)
	if err := row.Scan(&commentId); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return 0, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return 0, models.WrongData
	}

	row = r.db.QueryRowContext(ctx, DecLikesCount, commentInfo.CommentID)

	var likesCountTmp int64

	if err := row.Scan(&likesCountTmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return 0, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return 0, models.WrongData
	}

	row = r.db.QueryRowContext(ctx, DeleteLike, commentInfo.CommentID)

	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return 0, models.InternalError
	}

	return likesCountTmp, nil

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
