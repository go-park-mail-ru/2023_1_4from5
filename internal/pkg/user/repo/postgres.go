package repo

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

const (
	GET_PROFILE = "SELECT login, display_name, profile_photo, registration_date FROM public.user WHERE user_id=$1;"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) GetUserProfile(id uuid.UUID) (models.UserProfile, error) {
	var profile models.UserProfile

	row := ur.db.QueryRow(GET_PROFILE, id)
	if err := row.Scan(&profile.Login, &profile.Name, &profile.ProfilePhoto, &profile.Registration); err != nil {
		return models.UserProfile{}, errors.New("InternalError")
	}

	return profile, nil
}
