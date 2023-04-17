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
	AddDonate            = `INSERT INTO "donation"(user_id, creator_id, money_count) VALUES ($1, $2, $3);`
	BecameCreator        = `INSERT INTO "creator"(creator_id, user_id, name, description) VALUES ($1, $2, $3, $4);`
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
	return profile, nil
}

func (ur *UserRepo) GetHomePage(ctx context.Context, id uuid.UUID) (models.UserHomePage, error) {
	var page models.UserHomePage

	row := ur.db.QueryRowContext(ctx, UserNamePhoto, id)
	if err := row.Scan(&page.Name, &page.ProfilePhoto); err != nil && !errors.Is(sql.ErrNoRows, err) {
		ur.logger.Error(err)
		return models.UserHomePage{}, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.UserHomePage{}, models.NotFound
	}
	row = ur.db.QueryRowContext(ctx, CheckIfCreator, id)
	if err := row.Scan(&page.CreatorId); err != nil && !errors.Is(sql.ErrNoRows, err) {
		ur.logger.Error(err)
		return models.UserHomePage{}, models.InternalError
	} else if err == nil {
		page.IsCreator = true
	}
	return page, nil
}

func (ur *UserRepo) UpdateProfilePhoto(ctx context.Context, userID uuid.UUID, path uuid.UUID) error {
	row := ur.db.QueryRowContext(ctx, UpdateProfilePhoto, path, userID)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return models.InternalError
	}
	return nil
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

func (ur *UserRepo) Donate(ctx context.Context, donateInfo models.Donate, userID uuid.UUID) (int, error) {
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		ur.logger.Error(err)
		return 0, models.InternalError
	}

	row := tx.QueryRowContext(ctx, UpdateAuthorAimMoney, donateInfo.MoneyCount, donateInfo.CreatorID)
	if err != nil {
		_ = tx.Rollback()
		ur.logger.Error(err)
		return 0, models.InternalError
	}
	var newMoney int
	if err = row.Scan(&newMoney); err != nil && !errors.Is(sql.ErrNoRows, err) {
		_ = tx.Rollback()
		return 0, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		_ = tx.Rollback()
		return 0, models.WrongData
	}

	tx.QueryRowContext(ctx, AddDonate, userID, donateInfo.CreatorID, donateInfo.MoneyCount)

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
