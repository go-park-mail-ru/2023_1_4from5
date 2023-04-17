package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"go.uber.org/zap"
)

const (
	CreateSubscription = `INSERT INTO "subscription"(subscription_id,creator_id, month_cost, title, description) VALUES ($1, $2, $3, $4, $5);`
)

type SubscriptionRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewSubscriptionRepo(db *sql.DB, logger *zap.SugaredLogger) *SubscriptionRepo {
	return &SubscriptionRepo{
		db:     db,
		logger: logger,
	}
}

func (r *SubscriptionRepo) CreateSubscription(ctx context.Context, subscriptionInfo models.Subscription) error {
	row := r.db.QueryRowContext(ctx, CreateSubscription, subscriptionInfo.Id, subscriptionInfo.Creator, subscriptionInfo.MonthCost, subscriptionInfo.Title, subscriptionInfo.Description)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}
