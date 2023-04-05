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

func (repo *AttachmentRepo) CreateAttach(ctx context.Context, postID uuid.UUID, attachID uuid.UUID, attachmentType string) error {
	row := repo.db.QueryRow(InsertAttach, attachID, postID, attachmentType)

	if err := row.Err(); err != nil {
		repo.logger.Error(err)
		return models.InternalError
	}

	return nil
}

func (r *AttachmentRepo) DeleteAttach(ctx context.Context, attachID, postID uuid.UUID) error {
	var attachIDtmp uuid.UUID
	row := r.db.QueryRow(DeleteAttach, attachID, postID)
	if err := row.Scan(&attachIDtmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.WrongData
	}
	return nil
}

func (repo *AttachmentRepo) DeleteAttachByID(ctx context.Context, attachID uuid.UUID) error {
	row := repo.db.QueryRow(DeleteAttachByID, attachID)

	if err := row.Err(); err != nil {
		repo.logger.Error(err)
		return models.InternalError
	}

	return nil
}

func (repo *AttachmentRepo) DeleteAttachesByPostID(ctx context.Context, postID uuid.UUID) ([]models.AttachmentData, error) {
	resultAttachs := make([]models.AttachmentData, models.MaxFiles)
	rows, err := repo.db.Query(DeleteAttachByPostID, postID)
	if err != nil {
		repo.logger.Error(err)
		return nil, models.InternalError
	}
	defer rows.Close()
	var i uint16
	for rows.Next() {
		if err := rows.Scan(&resultAttachs[i].Id, &resultAttachs[i].Type); err != nil {
			return nil, models.InternalError
		}
		i++
	}
	return resultAttachs, nil
}
