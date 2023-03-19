package repo

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

const (
	InsertAttach = `INSERT INTO "attachment"(attachment_id, post_id, attachment_type) VALUES ($1,$2,$3)`
)

type AttachmentRepo struct {
	db *sql.DB
}

func NewAttachmentRepo(db *sql.DB) *AttachmentRepo {
	return &AttachmentRepo{db: db}
}

func (repo *AttachmentRepo) CreateAttach(postID uuid.UUID, attachment models.AttachmentData) error {
	attachment.Id = uuid.New()

	row := repo.db.QueryRow(InsertAttach, attachment.Id, postID, attachment.Type)

	if err := row.Err(); err != nil {
		return models.InternalError
	}

	return nil
}
