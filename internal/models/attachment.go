package models

import "github.com/google/uuid"

// easyjson -all ./internal/models/attachment.go

type Attachment struct {
	Id   uuid.UUID `json:"attachment_id"`
	Path string    `json:"attachment_path"`
}
