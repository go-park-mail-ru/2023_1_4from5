package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

const (
	UserAccessDetails = `SELECT user_id, password_hash, user_version FROM "user" WHERE login=$1;`
	AddUser           = `INSERT INTO "user"(user_id, login, display_name, profile_photo, password_hash) VALUES($1, $2, $3, $4, $5) RETURNING user_id;`
	IncUserVersion    = `UPDATE "user" SET user_version = user_version + 1 WHERE user_id=$1 RETURNING user_version;`
	CheckUserVersion  = `SELECT user_version FROM "user" WHERE user_id = $1`
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var id uuid.UUID
	user.Id = uuid.New()
	row := r.db.QueryRow(AddUser, user.Id, user.Login, user.Name, user.ProfilePhoto, user.PasswordHash)

	if err := row.Scan(&id); err != nil {
		fmt.Println(err)
		return models.User{}, models.InternalError
	}

	userOut := models.User{
		Id:           id,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		UserVersion:  user.UserVersion,
	}

	return userOut, nil
}

func (r *AuthRepo) CheckUser(ctx context.Context, user models.User) (models.User, error) {
	var (
		passwordHash string
		userVersion  int
		id           uuid.UUID
	)

	row := r.db.QueryRow(UserAccessDetails, user.Login) // Ищем пользователя с таким логином и берем его пароль и id и юзерверсию
	if err := row.Scan(&id, &passwordHash, &userVersion); err != nil && !errors.Is(sql.ErrNoRows, err) {
		fmt.Println(err)
		return models.User{}, models.InternalError
	}

	if passwordHash == user.PasswordHash { // совпал логин и пародь
		userOut := models.User{
			Id:           id,
			Login:        user.Login,
			PasswordHash: user.PasswordHash,
			UserVersion:  userVersion,
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

func (r *AuthRepo) IncUserVersion(ctx context.Context, userId uuid.UUID) (int, error) {
	row := r.db.QueryRow(IncUserVersion, userId)
	var userVersion int

	if err := row.Scan(&userVersion); err != nil {
		return 0, models.InternalError
	}

	return userVersion, nil
}

func (r *AuthRepo) CheckUserVersion(ctx context.Context, details models.AccessDetails) (int, error) {
	row := r.db.QueryRow(CheckUserVersion, details.Id)
	var userVersion int

	if err := row.Scan(&userVersion); err != nil {
		return 0, models.InternalError
	}

	if userVersion != details.UserVersion {
		return 0, models.Unauthorized
	}

	return userVersion, nil
}
