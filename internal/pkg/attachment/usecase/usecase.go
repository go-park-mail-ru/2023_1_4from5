package usecase

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"os"
)

type AttachmentUsecase struct {
	repo   attachment.AttachmentRepo
	logger *zap.SugaredLogger
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

func NewAttachmentUsecase(repo attachment.AttachmentRepo, logger *zap.SugaredLogger) *AttachmentUsecase {
	return &AttachmentUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *AttachmentUsecase) DeleteAttachesByPostID(ctx context.Context, postID uuid.UUID) error {
	attachs, err := u.repo.DeleteAttachesByPostID(ctx, postID)
	if err != nil {
		u.logger.Error(err)
		return err
	}
	err = u.DeleteAttaches(ctx, attachs...)
	return err
}

func (u *AttachmentUsecase) CreateAttaches(ctx context.Context, attachments ...models.AttachmentData) error {
	for i, attach := range attachments {
		attachmentType, ok := types[attach.Type]
		if !ok {
			if err := u.DeleteAttaches(ctx, attachments[:i]...); err != nil {
				u.logger.Error(err)
				return models.InternalError
			}
			return models.WrongData
		}
		f, err := os.Create(fmt.Sprintf("%s%s.%s", models.FolderPath, attach.Id.String(), attachmentType))
		if err != nil {
			if err := u.DeleteAttaches(ctx, attachments[:i]...); err != nil {
				u.logger.Error(err)
				return models.InternalError
			}
			return models.InternalError
		}

		defer f.Close()

		if _, err := io.Copy(f, attach.Data); err != nil {
			if err := u.DeleteAttaches(ctx, attachments[:i]...); err != nil {
				u.logger.Error(err)
				return models.InternalError
			}
			u.logger.Error(err)
			return models.InternalError
		}

	}
	return nil
}

//	func (u *AttachmentUsecase) DeleteAttachsByPostID(postID uuid.UUID) error {
//		attachIDs, err := u.repo.DeleteAttachByPostID(postID)
//		if err != nil {
//			return models.InternalError
//		}
//		for _, attachID := range attachIDs {
//			if err := deleteByFileName(attachID.String()); err != nil {
//				return err
//			}
//		}
//		return nil
//	}

func (u *AttachmentUsecase) DeleteAttaches(ctx context.Context, attachments ...models.AttachmentData) error {
	for _, file := range attachments {
		if err := deleteAttach(file); err != nil {
			u.logger.Error(err)
			return models.InternalError
		}
	}
	return nil
}

func deleteAttach(attach models.AttachmentData) error {
	filename := models.FolderPath + attach.Id.String() + "." + types[attach.Type]

	if err := os.Remove(filename); err != nil {
		return models.InternalError
	}
	return nil
}
