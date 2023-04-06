package models

// easyjson -all ./internal/models/subscription.go

import (
	"github.com/google/uuid"
	"html"
)

type Subscription struct {
	Id          uuid.UUID `json:"id"`
	Creator     uuid.UUID `json:"creator"`
	MonthConst  int       `json:"month_const"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (subscription *Subscription) Sanitize() {
	subscription.Title = html.EscapeString(subscription.Title)
	subscription.Description = html.EscapeString(subscription.Description)
}
