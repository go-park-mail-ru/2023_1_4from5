package attachment

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/attachment_mock.go -package=mock

type AttachmentUsecase interface {
	CreateAttaches(attachments ...models.AttachmentData) error
	DeleteAttaches(attachments ...models.AttachmentData) error
}

type AttachmentRepo interface {
	CreateAttach(postID uuid.UUID, attachID uuid.UUID, attachmentType string) error
	//DeleteAttachsByPostID(postID uuid.UUID) error
}
