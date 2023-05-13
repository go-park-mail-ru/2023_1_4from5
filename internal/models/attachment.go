package models

import (
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/google/uuid"
	"io"
)

const (
	MaxFileSize = 10 << 20
	MaxFormSize = 50 << 20
	MaxFiles    = 10
	FolderPath  = "/images/"
	//FolderPath = "./images"
)

// easyjson -all ./internal/models/attachment.go

type Attachment struct {
	Id   uuid.UUID `json:"id"`
	Type string    `json:"type"`
}

//easyjson:skip
type AttachmentData struct {
	Id   uuid.UUID
	Data io.Reader
	Type string
}

func (attachment *Attachment) AttachToModel(attach *generatedCreator.Attachment) error {
	attachID, err := uuid.Parse(attach.ID)
	if err != nil {
		return err
	}

	attachment.Id = attachID
	attachment.Type = attach.Type
	return nil
}
