package usecase

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (u *AttachmentUsecase) GetFileExtension(key string) (string, bool) {
	val, ok := types[key]
	return val, ok
}

func NewAttachmentUsecase(repo attachment.AttachmentRepo, logger *zap.SugaredLogger) *AttachmentUsecase {
	return &AttachmentUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *AttachmentUsecase) DeleteAttachmentsByPostID(ctx context.Context, postID uuid.UUID) error {
	attachs, err := u.repo.DeleteAttachmentsByPostID(ctx, postID)
	if err != nil {
		return err
	}
	return u.DeleteAttachmentsFiles(ctx, attachs...)
}

func (u *AttachmentUsecase) AddAttach(ctx context.Context, postID uuid.UUID, attachment models.AttachmentData) error {
	return u.repo.CreateAttachment(ctx, postID, attachment.Id, attachment.Type)
}

func (u *AttachmentUsecase) DeleteAttachmentsFiles(ctx context.Context, attachments ...models.AttachmentData) error {
	for _, file := range attachments {
		if err := u.DeleteAttachmentFile(ctx, file); err != nil {
			u.logger.Error(err)
			return models.InternalError
		}
	}
	return nil
}

func (u *AttachmentUsecase) DeleteAttachment(ctx context.Context, postID uuid.UUID, attach models.AttachmentData) error {
	if err := u.DeleteAttachmentsFiles(ctx, attach); err != nil {
		return err
	}
	return u.repo.DeleteAttachment(ctx, attach.Id, postID)
}

func (u *AttachmentUsecase) DeleteAttachmentFile(ctx context.Context, attachment models.AttachmentData) error {
	val, ok := u.GetFileExtension(attachment.Type)
	if !ok {
		return models.WrongData
	}
	filename := filepath.Join(models.FolderPath, attachment.Id.String()) + "." + val
	fmt.Println(filename)
	if err := os.Remove(filename); err != nil {
		u.logger.Error(err)
		return models.InternalError
	}
	return nil
}
