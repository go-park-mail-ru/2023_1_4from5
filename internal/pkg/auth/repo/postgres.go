package repo

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"time"
)

const (
	UserAccessDetails = "SELECT user_id, password_hash FROM public.user WHERE login=$1;"
	AddUser           = "INSERT INTO public.user(user_id, login, display_name, profile_photo, password_hash, registration_date) VALUES($1, $2, $3, $4, $5, $6) RETURNING user_id;"
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
	row := r.db.QueryRow(AddUser, user.Id, user.Login, user.Name, user.ProfilePhoto, user.PasswordHash, time.Now().UTC())

	if err := row.Scan(&id); err != nil {
		return models.User{}, models.InternalError
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
		passwordHash string
		id           uuid.UUID
	)

	row := r.db.QueryRow(UserAccessDetails, user.Login) // Ищем пользователя с таким логином и берем его пароль и id
	if err := row.Scan(&id, &passwordHash); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return models.User{}, models.InternalError
	}

	if passwordHash == user.PasswordHash { // совпал логин и пародь
		userOut := models.User{
			Id:           id,
			Login:        user.Login,
			PasswordHash: user.PasswordHash,
		}

		return userOut, nil
	}
	// запрос ничего не вернул, т.е. нет пользователя с таким логином
	if passwordHash == "" {
		return models.User{}, models.NotFound
	}
	// совпал логин, но не совпал пароль
	return models.User{}, models.WrongPassword
}
