package repo

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"net/http"
)

const (
	//SElECT_USER = "SELECT id, email, login, encrypted_password, created_at FROM public.users;"
	CHECK_USER  = "SELECT password_hash FROM public.user WHERE login=$1;"
	CREATE_USER = "INSERT INTO public.user(user_id, login, display_name, profile_photo, password_hash) VALUES($1, $2, $3, $4, $5) RETURNING user_id;"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(user models.User) (models.User, int) {
	var id uuid.UUID
	user.Id = uuid.New()
	row := r.db.QueryRow(CREATE_USER, user.Id, user.Login, user.DisplayName, user.ProfilePhoto, user.PasswordHash)

	err := row.Scan(&id)
	if err != nil {
		return models.User{}, http.StatusInternalServerError
	}

	userOut := models.User{
		Id:           id,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}

	return userOut, http.StatusOK
}

func (r *AuthRepo) CheckUser(user models.User) (models.User, int) {
	var (
		pwd string
		id  uuid.UUID
	)
	row := r.db.QueryRow(CHECK_USER, user.Login)

	if err := row.Scan(&id, &pwd); err != nil {
		return models.User{}, http.StatusInternalServerError
	}

	if pwd != user.PasswordHash {
		return models.User{}, http.StatusUnauthorized
	}

	userOut := models.User{
		Id:           id,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}

	return userOut, http.StatusOK
}
