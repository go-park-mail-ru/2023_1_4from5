package models

// easyjson -all ./internal/models/subscription.go

import (
	"github.com/google/uuid"
)

type Subscription struct {
	Id          uuid.UUID `json:"id"`
	Creator     uuid.UUID `json:"creator"`
	MonthConst  int       `json:"month_const"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
