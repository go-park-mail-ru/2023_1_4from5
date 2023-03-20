package models

import (
	"github.com/google/uuid"
	"mime/multipart"
)

// easyjson -all ./internal/models/attachment.go

type Attachment struct {
	Id   uuid.UUID
	Type string
}

type AttachmentData struct {
	Id   uuid.UUID
	Data multipart.File
	Type string
}
