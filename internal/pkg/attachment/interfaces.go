package attachment

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/attachment_mock.go -package=mock

type AttachmentUsecase interface {
	CreateAttaches(ctx context.Context, attachments ...models.AttachmentData) error
	DeleteAttaches(ctx context.Context, attachments ...models.AttachmentData) error
	DeleteAttachesByPostID(ctx context.Context, postID uuid.UUID) error
	DeleteAttach(ctx context.Context, attachID, postID uuid.UUID) error
}

type AttachmentRepo interface {
	CreateAttach(ctx context.Context, postID uuid.UUID, attachID uuid.UUID, attachmentType string) error
	DeleteAttachesByPostID(ctx context.Context, postID uuid.UUID) ([]models.AttachmentData, error)
	DeleteAttach(ctx context.Context, attachID, postID uuid.UUID) error
}
