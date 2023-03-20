package usecase

import (
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
	"audio/mp4":  "mp3", //TODO бля внатуре ?)_)
	"audio/mpeg": "mp3",
}

func NewAttachmentUsecase(repo attachment.AttachmentRepo) *AttachmentUsecase {
	return &AttachmentUsecase{repo: repo}
}

func (u *AttachmentUsecase) CreateAttachs(postID uuid.UUID, attachments ...models.AttachmentData) ([]uuid.UUID, error) {
	resultIds := make([]uuid.UUID, len(attachments))
	for i, attach := range attachments {
		fmt.Println(attach.Type)
		attachmentType, ok := types[attach.Type]
		if !ok {
			return nil, models.WrongData
		}
		resultIds[i] = uuid.New()

		f, err := os.Create(fmt.Sprintf("/home/ubuntu/frontend/2023_1_4from5/public/%s.%s", resultIds[i].String(), attachmentType))
		if err != nil {
			fmt.Println(err)
			return nil, models.InternalError
		}
		defer f.Close()
		if w, err := io.Copy(f, attach.Data); err != nil {
			return nil, models.InternalError
		} else {
			fmt.Println(w)
		}

		if err = u.repo.CreateAttach(postID, resultIds[i], attach.Type); err != nil {
			return nil, models.InternalError
		}
	}
	//TODO: всё на транзакции
	return resultIds, nil
}
