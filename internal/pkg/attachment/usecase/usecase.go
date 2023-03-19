package usecase

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	"github.com/google/uuid"
)

type AttachmentUsecase struct {
	repo attachment.AttachmentRepo
}

func NewAttachmentUsecase(repo attachment.AttachmentRepo) *AttachmentUsecase {
	return &AttachmentUsecase{repo: repo}
}

func (u *AttachmentUsecase) CreateAttach(postID uuid.UUID, attachments ...models.AttachmentData) error {
	for _, attach := range attachments {
		if err := u.repo.CreateAttach(postID, attach); err != nil {
			return models.InternalError
		}
	}
	//TODO: всё на транзакции
	return nil
}
