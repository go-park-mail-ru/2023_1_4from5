package models

import "github.com/google/uuid"

// easyjson -all ./internal/models/creator.go

type Attachment struct {
	Id   uuid.UUID
	Type string
}

type AttachmentData struct {
	Id   uuid.UUID `json:"id"`
	Data []byte    `json:"data"`
	Type string    `json:"type"`
}
