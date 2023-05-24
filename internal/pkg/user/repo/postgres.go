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
	UserProfile          = `SELECT login, display_name, profile_photo, registration_date FROM "user" WHERE user_id=$1;`
	UserNamePhoto        = `SELECT display_name, profile_photo FROM "user" WHERE user_id=$1;`
	CheckIfCreator       = `SELECT creator_id FROM "creator" WHERE user_id=$1;`
	UpdateProfilePhoto   = `UPDATE "user" SET profile_photo = $1 WHERE user_id = $2;`
	UpdatePassword       = `UPDATE "user" SET password_hash = $1, user_version = user_version+1 WHERE user_id = $2;`
	UpdateProfileInfo    = `UPDATE "user" SET login = $1, display_name = $2 WHERE user_id = $3;`
	UpdateAuthorAimMoney = `UPDATE "creator" SET money_got = money_got + $1 WHERE creator_id = $2 RETURNING money_got;`
	AddDonate            = `INSERT INTO "donation"(creator_id, money_count) VALUES ($1, $2);`
	BecameCreator        = `INSERT INTO "creator"(creator_id, user_id, name, description) VALUES ($1, $2, $3, $4);`
	Follow               = `INSERT INTO "follow" (user_id, creator_id) VALUES ($1, $2);`
	Unfollow             = `DELETE FROM "follow" WHERE user_id = $1 AND creator_id = $2;`
	CheckIfFollow        = `SELECT user_id FROM "follow" WHERE user_id = $1 AND creator_id = $2;`
	UpdateSubscription   = `UPDATE "user_subscription" SET expire_date = expire_date + $1 * INTERVAL '1 MONTH' WHERE user_id = $2 AND subscription_id = $3 RETURNING user_id;`
	Subscribe            = `INSERT INTO "user_subscription" VALUES ($1, $2, now() + $3 * INTERVAL '1 MONTH');`
	CheckIfSubExists     = `SELECT title, creator_id FROM subscription WHERE subscription_id = $1;`
	AddPaymentInfo       = `INSERT INTO "user_payments" (user_id, subscription_id, payment_timestamp, month_count, payment_info, money) VALUES ($1, $2, now(), $3, $4, 0);`
	CheckPaymentInfo     = `SELECT user_id, subscription_id, month_count FROM "user_payments" WHERE payment_info = $1;`
	UpdatePaymentInfo    = `UPDATE "user_payments" SET money = $1 WHERE payment_info = $2`
	UserSubscriptions    = `SELECT us.subscription_id, c.creator_id, name, profile_photo, month_cost, title, subscription.description FROM "subscription" join user_subscription us on subscription.subscription_id = us.subscription_id join creator c on c.creator_id = subscription.creator_id WHERE us.user_id = $1;`
	DeletePhoto          = `UPDATE "user" SET profile_photo = null WHERE user_id = $1`
	GetCreatorIDFromSub  = `SELECT creator_id FROM subscription WHERE subscription_id = $1 `
	FollowsList          = `SELECT c.creator_id, name, profile_photo, description FROM "follow" join creator c on c.creator_id = follow.creator_id WHERE follow.user_id = $1;`
)

type UserRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewUserRepo(db *sql.DB, logger *zap.SugaredLogger) *UserRepo {
	return &UserRepo{
		db:     db,
		logger: logger,
	}
}

func (ur *UserRepo) GetUserProfile(ctx context.Context, id uuid.UUID) (models.UserProfile, error) {
	var profile models.UserProfile

	row := ur.db.QueryRowContext(ctx, UserProfile, id)
	if err := row.Scan(&profile.Login, &profile.Name, &profile.ProfilePhoto, &profile.Registration); err != nil && !errors.Is(sql.ErrNoRows, err) {
		ur.logger.Error(err)
		return models.UserProfile{}, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.UserProfile{}, models.NotFound
	}
	var err error
	if profile.CreatorId, profile.IsCreator, err = ur.CheckIfCreator(ctx, id); err != nil {
		return models.UserProfile{}, models.InternalError
	}
	return profile, nil
}

func (ur *UserRepo) UserSubscriptions(ctx context.Context, userId uuid.UUID) ([]models.Subscription, error) {
	subs := make([]models.Subscription, 0)
	rows, err := ur.db.QueryContext(ctx, UserSubscriptions, userId)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		ur.logger.Error(err)
		return nil, models.InternalError
	}
	defer rows.Close()
	for rows.Next() {
		var sub models.Subscription
		var descriptionTmp sql.NullString
		err = rows.Scan(&sub.Id, &sub.Creator, &sub.CreatorName,
			&sub.CreatorPhoto, &sub.MonthCost, &sub.Title, &descriptionTmp)
		if err != nil {
			ur.logger.Error(err)
			return nil, models.InternalError
		}
		sub.Description = descriptionTmp.String

		subs = append(subs, sub)
	}
	return subs, nil
}

func (ur *UserRepo) UpdateProfilePhoto(ctx context.Context, userID uuid.UUID, path uuid.UUID) error {
	row := ur.db.QueryRowContext(ctx, UpdateProfilePhoto, path, userID)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (ur *UserRepo) CheckIfFollow(ctx context.Context, userId, creatorId uuid.UUID) (bool, error) {
	row := ur.db.QueryRowContext(ctx, CheckIfFollow, userId, creatorId)
	if err := row.Scan(&userId); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return false, models.InternalError
	} else if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, nil
}

func (ur *UserRepo) Follow(ctx context.Context, userId, creatorId uuid.UUID) error {
	row := ur.db.QueryRowContext(ctx, Follow, userId, creatorId)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (ur *UserRepo) Unfollow(ctx context.Context, userId, creatorId uuid.UUID) error {
	row := ur.db.QueryRowContext(ctx, Unfollow, userId, creatorId)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (ur *UserRepo) AddPaymentInfo(ctx context.Context, subscription models.SubscriptionDetails) error {
	row := ur.db.QueryRowContext(ctx, AddPaymentInfo, subscription.UserID, subscription.Id, subscription.MonthCount, subscription.PaymentInfo)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (ur *UserRepo) GetCreatorID(ctx context.Context, subscriptionID uuid.UUID) (uuid.UUID, error) {
	var creatorID uuid.UUID
	row := ur.db.QueryRowContext(ctx, GetCreatorIDFromSub, subscriptionID)
	if err := row.Scan(&creatorID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return uuid.Nil, models.InternalError
	}
	return creatorID, nil
}

func (ur *UserRepo) CheckPaymentInfo(ctx context.Context, paymentInfo uuid.UUID) (models.SubscriptionDetails, error) {
	var subscription models.SubscriptionDetails
	row := ur.db.QueryRowContext(ctx, CheckPaymentInfo, paymentInfo)
	if err := row.Scan(&subscription.UserID, &subscription.Id, &subscription.MonthCount); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return models.SubscriptionDetails{}, models.InternalError
	} else if errors.Is(err, sql.ErrNoRows) {
		return models.SubscriptionDetails{}, models.NotFound
	}
	return subscription, nil
}

func (ur *UserRepo) UpdatePaymentInfo(ctx context.Context, money float32, paymentInfo uuid.UUID) error {
	row := ur.db.QueryRowContext(ctx, UpdatePaymentInfo, money, paymentInfo)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (ur *UserRepo) Subscribe(ctx context.Context, subscription models.SubscriptionDetails) (models.NotificationSubInfo, error) {
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		ur.logger.Error(err)
		return models.NotificationSubInfo{}, models.InternalError
	}
	// если подписка уже есть обновляем expire date
	row := ur.db.QueryRowContext(ctx, UpdateSubscription, subscription.MonthCount, subscription.UserID, subscription.Id)
	var userIDtmp uuid.UUID
	var subNotification models.NotificationSubInfo
	if err := row.Scan(&userIDtmp); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		_ = tx.Rollback()
		return models.NotificationSubInfo{}, models.InternalError
	} else if errors.Is(err, sql.ErrNoRows) { // если нет, то добавляем о ней запись
		row = ur.db.QueryRowContext(ctx, CheckIfSubExists, subscription.Id)
		if err = row.Scan(&subNotification.SubscriptionName, &subNotification.CreatorID); err != nil && !errors.Is(err, sql.ErrNoRows) {
			ur.logger.Error(err)
			_ = tx.Rollback()
			return models.NotificationSubInfo{}, models.InternalError
		} else if errors.Is(err, sql.ErrNoRows) { // такой подписки нет
			_ = tx.Rollback()
			return models.NotificationSubInfo{}, models.WrongData
		}

		row = ur.db.QueryRowContext(ctx, Subscribe, subscription.UserID, subscription.Id, subscription.MonthCount)
		if err = row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
			ur.logger.Error(err)
			_ = tx.Rollback()
			return models.NotificationSubInfo{}, models.InternalError
		}
	}

	if err = tx.Commit(); err != nil {
		ur.logger.Error(err)
		return models.NotificationSubInfo{}, models.InternalError
	}

	return subNotification, nil
}

func (ur *UserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	row := ur.db.QueryRowContext(ctx, UpdatePassword, password, id)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (ur *UserRepo) UpdateProfileInfo(ctx context.Context, profileInfo models.UpdateProfileInfo, id uuid.UUID) error {
	row := ur.db.QueryRowContext(ctx, UpdateProfileInfo, profileInfo.Login, profileInfo.Name, id)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (ur *UserRepo) Donate(ctx context.Context, donateInfo models.Donate) (float32, error) {
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		ur.logger.Error(err)
		return 0, models.InternalError
	}
	var newMoney float32
	row := tx.QueryRowContext(ctx, UpdateAuthorAimMoney, donateInfo.MoneyCount, donateInfo.CreatorID)

	if err = row.Scan(&newMoney); err != nil && !errors.Is(sql.ErrNoRows, err) {
		ur.logger.Error(err)
		_ = tx.Rollback()
		return 0, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		_ = tx.Rollback()
		return 0, models.WrongData
	}

	tx.QueryRowContext(ctx, AddDonate, donateInfo.CreatorID, donateInfo.MoneyCount)

	if err = tx.Commit(); err != nil {
		ur.logger.Error(err)
		return 0, models.InternalError
	}

	return newMoney, nil
}

func (ur *UserRepo) CheckIfCreator(ctx context.Context, userId uuid.UUID) (uuid.UUID, bool, error) {
	idTmp := uuid.UUID{}
	row := ur.db.QueryRowContext(ctx, CheckIfCreator, userId)
	if err := row.Scan(&idTmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		ur.logger.Error(err)
		return uuid.Nil, false, models.InternalError
	} else if err == nil {
		return idTmp, true, nil
	}
	return uuid.Nil, false, nil
}

func (ur *UserRepo) BecomeCreator(ctx context.Context, creatorInfo models.BecameCreatorInfo, userId uuid.UUID) (uuid.UUID, error) {
	creatorId := uuid.New()
	row := ur.db.QueryRowContext(ctx, BecameCreator, creatorId, userId, creatorInfo.Name, creatorInfo.Description)
	if err := row.Scan(); err != nil && !errors.Is(sql.ErrNoRows, err) {
		ur.logger.Error(err)
		return uuid.Nil, models.InternalError
	}
	return creatorId, nil
}

func (ur *UserRepo) DeletePhoto(ctx context.Context, userId uuid.UUID) error {
	row := ur.db.QueryRowContext(ctx, DeletePhoto, userId)
	if err := row.Scan(); err != nil && !errors.Is(sql.ErrNoRows, err) {
		ur.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (ur *UserRepo) UserFollows(ctx context.Context, userId uuid.UUID) ([]models.Follow, error) {
	follows := make([]models.Follow, 0)
	rows, err := ur.db.QueryContext(ctx, FollowsList, userId)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		ur.logger.Error(err)
		return nil, models.InternalError
	}
	defer rows.Close()
	for rows.Next() {
		var follow models.Follow
		var descriptionTmp sql.NullString
		err = rows.Scan(&follow.Creator, &follow.CreatorName, &follow.CreatorPhoto, &descriptionTmp)
		if err != nil {
			ur.logger.Error(err)
			return nil, models.InternalError
		}
		follow.Description = descriptionTmp.String

		follows = append(follows, follow)
	}
	return follows, nil
}
