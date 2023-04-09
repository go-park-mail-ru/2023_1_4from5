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
	"path/filepath"
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
	"audio/mp4":  "mp3",
	"audio/mpeg": "mp3",
}

func GetFileExtension(key string) (string, bool) {
	val, ok := types[key]
	return val, ok
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
		return err
	}
	return u.DeleteAttaches(ctx, attachs...)
}

func (u *AttachmentUsecase) CreateAttaches(ctx context.Context, attachments ...models.AttachmentData) error {
	for i, attach := range attachments {
		attachmentType, ok := GetFileExtension(attach.Type)
		if !ok {
			if err := u.DeleteAttaches(ctx, attachments[:i]...); err != nil {
				u.logger.Error(err)
				return models.InternalError
			}
			return models.Unsupported
		}
		f, err := os.Create(fmt.Sprintf("%s.%s", filepath.Join(models.FolderPath, attach.Id.String()), attachmentType))
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

func (u *AttachmentUsecase) DeleteAttaches(ctx context.Context, attachments ...models.AttachmentData) error {
	for _, file := range attachments {
		if err := deleteAttach(file); err != nil {
			u.logger.Error(err)
			return models.InternalError
		}
	}
	return nil
}

func (u *AttachmentUsecase) DeleteAttach(ctx context.Context, attachID, postID uuid.UUID) error {
	return u.repo.DeleteAttach(ctx, attachID, postID)
}

func deleteAttach(attach models.AttachmentData) error {
	val, ok := GetFileExtension(attach.Type)
	if !ok {
		return models.WrongData
	}
	filename := filepath.Join(models.FolderPath, attach.Id.String()) + "." + val

	if err := os.Remove(filename); err != nil {
		fmt.Println(err)
		return models.InternalError
	}
	return nil
}
