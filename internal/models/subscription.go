package models

// easyjson -all ./internal/models/subscription.go

import (
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	"github.com/google/uuid"
	"html"
)

type Subscription struct {
	Id           uuid.UUID `json:"id,omitempty"`
	Creator      uuid.UUID `json:"creator,omitempty"`
	CreatorName  string    `json:"creator_name,omitempty"`
	CreatorPhoto uuid.UUID `json:"creator_photo,omitempty"`
	MonthCost    int64     `json:"month_cost"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
}

type Follow struct {
	Creator      uuid.UUID `json:"creator,omitempty"`
	CreatorName  string    `json:"creator_name,omitempty"`
	CreatorPhoto uuid.UUID `json:"creator_photo,omitempty"`
	Description  string    `json:"description,omitempty"`
}

type SubscriptionDetails struct {
	CreatorId  uuid.UUID `json:"creator_id"`
	Id         uuid.UUID `json:"id,omitempty"`
	UserID     uuid.UUID `json:"user_id,omitempty"`
	MonthCount int64     `json:"month_count"`
	Money      int64     `json:"money"`
}

func (subscription *Subscription) Sanitize() {
	subscription.Title = html.EscapeString(subscription.Title)
	subscription.Description = html.EscapeString(subscription.Description)
}

func (follow *Follow) Sanitize() {
	follow.CreatorName = html.EscapeString(follow.CreatorName)
	follow.Description = html.EscapeString(follow.Description)
}

func (subscription *Subscription) IsValid() bool {
	return 0 < len(subscription.Title) && len(subscription.Title) < 41 && len(subscription.Description) < 201
}

func (subscription *Subscription) ProtoSubscriptionToModel(sub *generatedCommon.Subscription) error {
	subID, err := uuid.Parse(sub.Id)
	if err != nil {
		return err
	}
	creatorID, err := uuid.Parse(sub.Creator)
	if err != nil {
		return err
	}
	creatorPhoto, err := uuid.Parse(sub.CreatorPhoto)
	if err != nil {
		return err
	}
	subscription.Id = subID
	subscription.Creator = creatorID
	subscription.CreatorName = sub.CreatorName
	subscription.CreatorPhoto = creatorPhoto
	subscription.MonthCost = sub.MonthCost
	subscription.Title = sub.Title
	subscription.Description = sub.Description
	return nil
}
