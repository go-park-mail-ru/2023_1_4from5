package repo

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

type CreatorRepo struct {
	db *sql.DB
}

func NewCreatorRepo(db *sql.DB) *CreatorRepo {
	return &CreatorRepo{db: db}
}

const (
	CREATOR_INFO = "SELECT user_id, name, cover_photo, followers_count, description, posts_count FROM public.creator WHERE creator_id=$1;"
	//CREATOR_POSTS = "SELECT name, cover_photo, followers_count, description, posts_count FROM public.creator WHERE user_id=$1;"
)

func (ur *CreatorRepo) GetPage(userId uuid.UUID, creatorId uuid.UUID) (models.CreatorPage, error) {
	var creatorPage models.CreatorPage
	row := ur.db.QueryRow(CREATOR_INFO, creatorId)
	if err := row.Scan(&creatorPage.CreatorInfo.UserId, &creatorPage.CreatorInfo.Name, &creatorPage.CreatorInfo.CoverPhoto, &creatorPage.CreatorInfo.FollowersCount, &creatorPage.CreatorInfo.Description, &creatorPage.CreatorInfo.PostsCount); err != nil {
		return models.CreatorPage{}, errors.New("InternalError")
	}
	if creatorPage.CreatorInfo.UserId == userId {
		creatorPage.IsMyPage = true
	}
	return creatorPage, nil
}
