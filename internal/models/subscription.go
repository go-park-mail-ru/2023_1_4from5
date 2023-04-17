package models

// easyjson -all ./internal/models/subscription.go

import (
	"github.com/google/uuid"
	"html"
)

type Subscription struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Creator     uuid.UUID `json:"creator,omitempty"`
	MonthCost   int       `json:"month_cost"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
}

func (subscription *Subscription) Sanitize() {
	subscription.Title = html.EscapeString(subscription.Title)
	subscription.Description = html.EscapeString(subscription.Description)
}

func (subscription *Subscription) IsValid() bool {
	return 0 < len(subscription.Title) && len(subscription.Title) < 41 && len(subscription.Description) < 201
}
