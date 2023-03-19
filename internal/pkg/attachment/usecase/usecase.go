package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	"github.com/google/uuid"
	"os"
)

type AttachmentUsecase struct {
	repo attachment.AttachmentRepo
}

func NewAttachmentUsecase(repo attachment.AttachmentRepo) *AttachmentUsecase {
	return &AttachmentUsecase{repo: repo}
}

func (u *AttachmentUsecase) CreateAttach(postID uuid.UUID, attachments ...models.AttachmentData) error {
	for _, attach := range attachments {
		attachId, err := u.repo.CreateAttach(postID, attach)
		if err != nil {
			return models.InternalError
		}

		file, err := os.Create(fmt.Sprintf("/home/ubuntu/frontend/2023_1_4from5/public/%s.%s", attachId.String(), attach.Type))
		if err != nil {
			return models.InternalError
		}
		defer file.Close()

		_, err = file.Write(attach.Data)
		if err != nil {
			return models.InternalError
		}
	}
	//TODO: всё на транзакции
	return nil
}
