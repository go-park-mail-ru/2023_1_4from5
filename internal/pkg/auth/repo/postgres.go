package repo

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"time"
)

const (
	//SElECT_USER = "SELECT id, email, login, encrypted_password, created_at FROM public.user;"
	CHECK_USER  = "SELECT user_id, login FROM public.user WHERE login=$1;"
	CREATE_USER = "INSERT INTO public.user(user_id, login, display_name, profile_photo, password_hash, registration_date) VALUES($1, $2, $3, $4, $5, $6) RETURNING user_id;"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(user models.User) (models.User, error) {
	var id uuid.UUID
	user.Id = uuid.New()
	row := r.db.QueryRow(CREATE_USER, user.Id, user.Login, user.Name, user.ProfilePhoto, user.PasswordHash, time.Now())

	err := row.Scan(&id)
	if err != nil {
		return models.User{}, err
	}

	userOut := models.User{
		Id:           id,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}

	return userOut, nil
}

func (r *AuthRepo) CheckUser(user models.User) (models.User, error) {
	var (
		login string
		id    uuid.UUID
	)

	row := r.db.QueryRow(CHECK_USER, user.Login)
	if err := row.Scan(&id, &login); err != nil {
		return models.User{}, errors.New("InternalError")
	}

	if user.Login == login { //TODO: TROUBLE
		return models.User{}, errors.New("Conflict")
	}

	userOut := models.User{
		Id:           id,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}

	return userOut, nil
}
