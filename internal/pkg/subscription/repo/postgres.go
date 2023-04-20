package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	CreateSubscription        = `INSERT INTO "subscription"(subscription_id,creator_id, month_cost, title, description) VALUES ($1, $2, $3, $4, $5);`
	DeleteSubsccriptionsPosts = `DELETE FROM "post_subscription" WHERE subscription_id = $1;`
	DeleteUsersSubscriptions  = `DELETE FROM "user_subscription" WHERE subscription_id = $1;`
	DeleteSubscription        = `DELETE FROM "subscription" WHERE subscription_id = $1 AND creator_id = $2 RETURNING subscription_id;`
	EditSubscription          = `UPDATE "subscription" SET month_cost = $1, title = $2, description = $3 WHERE subscription_id = $4;`
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
func (r *SubscriptionRepo) DeleteSubscription(ctx context.Context, subscriptionID, creatorID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error(err)
		return models.InternalError
	}
	row, err := tx.QueryContext(ctx, DeleteSubsccriptionsPosts, subscriptionID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}
	if err = row.Close(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	row, err = tx.QueryContext(ctx, DeleteUsersSubscriptions, subscriptionID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}
	if err = row.Close(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	row, err = tx.QueryContext(ctx, DeleteSubscription, subscriptionID, creatorID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}
	if !row.Next() {
		_ = tx.Rollback()
		return models.Forbbiden
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	return nil
}

func (r *SubscriptionRepo) EditSubscription(ctx context.Context, subscriptionNewInfo models.Subscription) error {
	row := r.db.QueryRowContext(ctx, EditSubscription, subscriptionNewInfo.MonthCost, subscriptionNewInfo.Title, subscriptionNewInfo.Description, subscriptionNewInfo.Id)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}
