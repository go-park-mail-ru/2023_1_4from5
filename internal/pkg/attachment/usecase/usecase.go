package usecase

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	"github.com/google/uuid"
	"io"
	"os"
)

type AttachmentUsecase struct {
	repo attachment.AttachmentRepo
}

var types = map[string]string{
	"image/jpeg": "jpeg",
	"image/png":  "png",
	"image/webp": "webp",
	"video/mpeg": "mpeg",
	"video/mp4":  "mp4",
	"audio/mp4":  "mp3", //TODO реально?))
	"audio/mpeg": "mp3",
}

func NewAttachmentUsecase(repo attachment.AttachmentRepo) *AttachmentUsecase {
	return &AttachmentUsecase{repo: repo}
}

func (u *AttachmentUsecase) CreateAttachs(postID uuid.UUID, attachments ...models.AttachmentData) ([]uuid.UUID, error) {
	resultIds := make([]uuid.UUID, len(attachments))
	for i, attach := range attachments {
		attachmentType, ok := types[attach.Type]
		if !ok {
			return nil, models.WrongData
		}
		resultIds[i] = uuid.New()

		f, err := os.Create(fmt.Sprintf("/images/%s.%s", resultIds[i].String(), attachmentType))
		if err != nil {
			fmt.Println(err)
			return nil, models.InternalError
		}
		defer f.Close()
		if _, err := io.Copy(f, attach.Data); err != nil {
			return nil, models.InternalError
		}

		if err = u.repo.CreateAttach(postID, resultIds[i], attach.Type); err != nil {
			return nil, models.InternalError
		}
	}
	//TODO: всё на транзакции
	return resultIds, nil
}

func (u *AttachmentUsecase) DeleteAttachsByID(attachmentIDs ...uuid.UUID) error {
	for _, attachId := range attachmentIDs {
		if _, err := os.Stat(fmt.Sprintf("/images/%s", attachId)); errors.Is(err, os.ErrNotExist) {
			return models.WrongData
		}

		if err := u.repo.DeleteAttachByID(attachId); err != nil {
			return models.InternalError
		}

		if err := os.Remove(fmt.Sprintf("/images/%s", attachId)); err != nil {
			return models.WrongData
		}
	}

	return nil
}

func (u *AttachmentUsecase) DeleteAttachsByPostID(postID uuid.UUID) error {
	attachIDs, err := u.repo.DeleteAttachByPostID(postID)
	if err != nil {
		return models.InternalError
	}
	for _, attachID := range attachIDs {
		if err := os.Remove(fmt.Sprintf("/images/%s", attachID)); err != nil {
			return models.WrongData
		}
	}
	return nil
}
