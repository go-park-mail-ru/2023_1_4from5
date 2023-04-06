package models

// easyjson -all ./internal/models/subscription.go

import (
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

type Subscription struct {
	Id          uuid.UUID `json:"id"`
	Creator     uuid.UUID `json:"creator"`
	MonthConst  int       `json:"month_const"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (subscription *Subscription) Sanitize() {
	sanitizer := bluemonday.StrictPolicy()
	subscription.Title = sanitizer.Sanitize(subscription.Title)
	subscription.Description = sanitizer.Sanitize(subscription.Description)
}
