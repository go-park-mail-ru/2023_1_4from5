package models

import (
	"github.com/google/uuid"
	"io"
)

const (
	MaxFileSize = 5 << 20
	MaxFormSize = 20 << 20
	MaxFiles    = 10
	FolderPath  = "images/"
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
