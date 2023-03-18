package repo

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

const (
	UserProfile    = `SELECT login, display_name, profile_photo, registration_date FROM "user" WHERE user_id=$1;`
	UserNamePhoto  = `SELECT display_name, profile_photo FROM "user" WHERE user_id=$1;`
	CheckIfCreator = `SELECT creator_id FROM "creator" WHERE user_id=$1;`
	//GET_USER_POSTS = "SELECT post_id, creator_id, creation_date, title, post_text, attachments, available_subscriptions FROM public.post WHERE UNNEST(available_subscriptions) IN (SELECT subscriptions FROM public.user WHERE user_id = $1);"
	//SELECT array_agg(aa.id)::int[] FROM UNNEST(ARRAY[1,2,3,45,67,8,8]) AS aa (id) WHERE aa.id = ANY(ARRAY[1,2,3,4,56,56,56,56]);
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) GetUserProfile(id uuid.UUID) (models.UserProfile, error) {
	var profile models.UserProfile

	row := ur.db.QueryRow(UserProfile, id)
	if err := row.Scan(&profile.Login, &profile.Name, &profile.ProfilePhoto, &profile.Registration); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return models.UserProfile{}, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.UserProfile{}, models.NotFound
	}
	return profile, nil
}

func (ur *UserRepo) GetHomePage(id uuid.UUID) (models.UserHomePage, error) {
	var page models.UserHomePage

	row := ur.db.QueryRow(UserNamePhoto, id)
	if err := row.Scan(&page.Name, &page.ProfilePhoto); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return models.UserHomePage{}, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.UserHomePage{}, models.NotFound
	}
	row = ur.db.QueryRow(CheckIfCreator, id)
	if err := row.Scan(&page.CreatorId); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return models.UserHomePage{}, models.InternalError
	} else if err == nil {
		page.IsCreator = true
	}
	return page, nil
}
