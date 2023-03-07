package repo

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"time"
)

const (
	//SElECT_USER = "SELECT id, email, login, encrypted_password, created_at FROM public.user;"
	CHECK_USER        = "SELECT user_id, password_hash FROM test.user WHERE login=$1;"
	CREATE_USER       = "INSERT INTO test.user(user_id, login, display_name, profile_photo, password_hash, registration_date) VALUES($1, $2, $3, $4, $5, $6) RETURNING user_id;"
	INC_USERVERSION   = "UPDATE test.user SET user_version = user_version + 1 WHERE user_id=$1 RETURNING user_version;"
	CHECK_USERVERSION = "SELECT user_version FROM test.user WHERE user_id = $1"
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
	row := r.db.QueryRow(CHECK_USER, user.Login) // Ищем пользователя с таким логином и берем его пароль и id
	if err := row.Scan(&id, &passwordHash); err != nil {
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

func (r *AuthRepo) IncUserVersion(userId uuid.UUID) (int, error) {
	row := r.db.QueryRow(INC_USERVERSION, userId)
	var userVersion int

	if err := row.Scan(&userVersion); err != nil {
		return 0, models.InternalError
	}

	return userVersion, nil
}

func (r *AuthRepo) CheckUserVersion(details models.AccessDetails) (int, error) {
	row := r.db.QueryRow(CHECK_USERVERSION, details.Id)
	var userVersion int

	if err := row.Scan(&userVersion); err != nil {
		return 0, models.InternalError
	}

	if userVersion != details.UserVersion {
		return 0, models.Unauthorized
	}

	return userVersion, nil
}
