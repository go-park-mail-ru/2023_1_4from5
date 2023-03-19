package attachment

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/attachment_mock.go -package=mock

type AttachmentUsecase interface {
	CreateAttach(postID uuid.UUID, attachments ...models.AttachmentData) error
}

type AttachmentRepo interface {
	CreateAttach(postID uuid.UUID, attachment models.AttachmentData) (uuid.UUID, error)
}
