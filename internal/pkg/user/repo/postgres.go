package repo

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

const (
	GET_PROFILE    = "SELECT login, display_name, profile_photo, registration_date FROM public.user WHERE user_id=$1;"
	GET_NAME_PHOTO = "SELECT display_name, profile_photo FROM public.user WHERE user_id=$1;"
	CHECK_CREATOR  = "SELECT creator_id FROM public.creator WHERE user_id=$1;"
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

	row := ur.db.QueryRow(GET_PROFILE, id)
	if err := row.Scan(&profile.Login, &profile.Name, &profile.ProfilePhoto, &profile.Registration); err != nil {
		return models.UserProfile{}, errors.New("InternalError")
	}

	return profile, nil
}

func (ur *UserRepo) GetHomePage(id uuid.UUID) (models.UserHomePage, error) {
	var page models.UserHomePage

	row := ur.db.QueryRow(GET_NAME_PHOTO, id)
	if err := row.Scan(&page.Name, &page.ProfilePhoto); err != nil {
		return models.UserHomePage{}, errors.New("InternalError")
	}

	row = ur.db.QueryRow(CHECK_CREATOR, id)
	if err := row.Scan(&page.CreatorId); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return models.UserHomePage{}, errors.New("InternalError")
	} else {
		page.IsCreator = true
	}

	return page, nil
}
