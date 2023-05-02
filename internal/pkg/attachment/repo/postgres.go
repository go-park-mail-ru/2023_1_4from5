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
	InsertAttach         = `INSERT INTO "attachment"(attachment_id, post_id, attachment_type) VALUES ($1,$2,$3)`
	DeleteAttachByID     = `DELETE FROM "attachment" WHERE attachment_id = $1`
	DeleteAttachByPostID = `DELETE FROM "attachment" WHERE post_id = $1 RETURNING attachment_id, attachment_type`
	DeleteAttach         = `DELETE FROM "attachment" WHERE attachment_id = $1 AND post_id = $2 RETURNING attachment_id`
)

type AttachmentRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewAttachmentRepo(db *sql.DB, logger *zap.SugaredLogger) *AttachmentRepo {
	return &AttachmentRepo{
		db:     db,
		logger: logger,
	}
}

func (repo *AttachmentRepo) CreateAttachment(ctx context.Context, postID uuid.UUID, attachmentID uuid.UUID, attachmentType string) error {
	row := repo.db.QueryRowContext(ctx, InsertAttach, attachmentID, postID, attachmentType)

	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		repo.logger.Error(err)
		return models.InternalError
	}

	return nil
}

func (r *AttachmentRepo) DeleteAttachment(ctx context.Context, attachmentID, postID uuid.UUID) error {
	var attachmentIDtmp uuid.UUID
	row := r.db.QueryRowContext(ctx, DeleteAttach, attachmentID, postID)
	if err := row.Scan(&attachmentIDtmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.WrongData
	}
	return nil
}

func (repo *AttachmentRepo) DeleteAttachmentByID(ctx context.Context, attachID uuid.UUID) error {
	row := repo.db.QueryRowContext(ctx, DeleteAttachByID, attachID)

	if err := row.Err(); err != nil {
		repo.logger.Error(err)
		return models.InternalError
	}

	return nil
}

func (repo *AttachmentRepo) DeleteAttachmentsByPostID(ctx context.Context, postID uuid.UUID) ([]models.Attachment, error) {
	resultAttachs := make([]models.Attachment, 0)
	rows, err := repo.db.Query(DeleteAttachByPostID, postID)
	if err != nil {
		repo.logger.Error(err)
		return nil, models.InternalError
	}
	defer rows.Close()
	tmp := models.Attachment{}
	for rows.Next() {
		if err := rows.Scan(&tmp.Id, &tmp.Type); err != nil {
			repo.logger.Error(err)
			return nil, models.InternalError
		}
		resultAttachs = append(resultAttachs, tmp)
	}
	return resultAttachs, nil
}
