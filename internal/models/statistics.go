package models

import (
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/google/uuid"
	"time"
)

// easyjson -all ./internal/models/statistics.go

type Statistics struct {
	CreatorId              uuid.UUID `json:"creator_id"`
	PostsPerMonth          int64     `json:"posts_per_month"`
	SubscriptionsBought    int64     `json:"subscriptions_bought"`
	DonationsCount         int64     `json:"donations_count"`
	MoneyFromDonations     float64   `json:"money_from_donations"`
	MoneyFromSubscriptions float64   `json:"money_from_subscriptions"`
	NewFollowers           int64     `json:"new_followers"`
	LikesCount             int64     `json:"likes_count"`
	CommentsCount          int64     `json:"comments_count"`
}

type StatisticsDates struct {
	CreatorId   uuid.UUID `json:"creator_id,omitempty"`
	FirstMonth  time.Time `json:"first_month"`
	SecondMonth time.Time `json:"second_month"`
}

func (statistics *Statistics) StatToModel(statInfo *generatedCreator.Stat) error {
	creatorId, err := uuid.Parse(statInfo.CreatorId)
	if err != nil {
		return err
	}

	statistics.CreatorId = creatorId
	statistics.PostsPerMonth = statInfo.PostsPerMonth
	statistics.SubscriptionsBought = statInfo.SubscriptionsBought
	statistics.DonationsCount = statInfo.DonationsCount
	statistics.MoneyFromDonations = statInfo.MoneyFromDonations
	statistics.MoneyFromSubscriptions = statInfo.MoneyFromSubscriptions
	statistics.NewFollowers = statInfo.NewFollowers
	statistics.LikesCount = statInfo.LikesCount
	statistics.CommentsCount = statInfo.CommentsCount
	return nil
}
