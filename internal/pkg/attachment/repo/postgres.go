package repo

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

const (
	InsertAttach         = `INSERT INTO "attachment"(attachment_id, post_id, attachment_type) VALUES ($1,$2,$3)`
	DeleteAttachByID     = `DELETE FROM "attachment" WHERE attachment_id = $1`
	DeleteAttachByPostID = `DELETE FROM "attachment" WHERE post_id = $1 RETURNING attachment_id`
)

type AttachmentRepo struct {
	db *sql.DB
}

func NewAttachmentRepo(db *sql.DB) *AttachmentRepo {
	return &AttachmentRepo{db: db}
}

func (repo *AttachmentRepo) CreateAttach(postID uuid.UUID, attachID uuid.UUID, attachmentType string) error {
	row := repo.db.QueryRow(InsertAttach, attachID, postID, attachmentType)

	if err := row.Err(); err != nil {
		return models.InternalError
	}

	return nil
}

func (repo *AttachmentRepo) DeleteAttachByID(attachID uuid.UUID) error {
	row := repo.db.QueryRow(DeleteAttachByID, attachID)

	if err := row.Err(); err != nil {
		return models.InternalError
	}

	return nil
}

func (repo *AttachmentRepo) DeleteAttachByPostID(postID uuid.UUID) ([]uuid.UUID, error) {
	resultIds := make([]uuid.UUID, 0)
	rows, err := repo.db.Query(DeleteAttachByPostID, postID)
	if err != nil {
		return nil, models.InternalError
	}
	defer rows.Close()
	for rows.Next() {
		var attachmentID uuid.UUID
		if err := rows.Scan(&attachmentID); err != nil {
			return nil, models.InternalError
		}
		resultIds = append(resultIds, attachmentID)
	}
	return resultIds, nil
}
