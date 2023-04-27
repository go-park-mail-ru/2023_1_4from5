package attachment

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/attachment_mock.go -package=mock

type AttachmentUsecase interface {
	DeleteAttachmentsFiles(ctx context.Context, attachments ...models.AttachmentData) error
	DeleteAttachmentsByPostID(ctx context.Context, postID uuid.UUID) error
	DeleteAttachment(ctx context.Context, postID uuid.UUID, attach models.AttachmentData) error
	AddAttach(ctx context.Context, postID uuid.UUID, attachment models.AttachmentData) error
	GetFileExtension(key string) (string, bool)
}

type AttachmentRepo interface {
	CreateAttachment(ctx context.Context, postID uuid.UUID, attachmentID uuid.UUID, attachmentType string) error
	DeleteAttachmentsByPostID(ctx context.Context, postID uuid.UUID) ([]models.AttachmentData, error)
	DeleteAttachment(ctx context.Context, attachmentID, postID uuid.UUID) error
}
