package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CreatorRepo struct {
	db *sql.DB
}

func NewCreatorRepo(db *sql.DB) *CreatorRepo {
	return &CreatorRepo{db: db}
}

const (
	CREATOR_INFO      = "SELECT user_id, name, cover_photo, followers_count, description, posts_count FROM public.creator WHERE creator_id=$1;"
	CREATOR_POSTS     = "SELECT post_id, creation_date, title, post_text, attachments, available_subscriptions  FROM public.post WHERE creator_id=$1;"
	USER_SUBSCRIPTION = "SELECT subscriptions FROM public.user WHERE user_id=$1;"
)

func (ur *CreatorRepo) GetPage(userId uuid.UUID, creatorId uuid.UUID) (models.CreatorPage, error) {
	var creatorPage models.CreatorPage
	userSubscriptions := make([]uuid.UUID, 0)
	fmt.Println("GET PAGE REPO")
	row := ur.db.QueryRow(CREATOR_INFO, creatorId)
	if err := row.Scan(&creatorPage.CreatorInfo.UserId, &creatorPage.CreatorInfo.Name, &creatorPage.CreatorInfo.CoverPhoto, &creatorPage.CreatorInfo.FollowersCount, &creatorPage.CreatorInfo.Description, &creatorPage.CreatorInfo.PostsCount); err != nil {
		return models.CreatorPage{}, errors.New("InternalError")
	}
	if creatorPage.CreatorInfo.UserId == userId {
		creatorPage.IsMyPage = true
	} else {
		fmt.Println("Not Author")
		row := ur.db.QueryRow(USER_SUBSCRIPTION, userId)
		if err := row.Scan(&userSubscriptions); err != nil {
			return models.CreatorPage{}, errors.New("InternalError")
		}
		fmt.Println(userSubscriptions)
	}

	rows, err := ur.db.Query(CREATOR_POSTS, creatorId)
	if err != nil {
		return models.CreatorPage{}, errors.New("InternalError")
	}
	posts := make([]models.Post, 0)
	for rows.Next() {
		var post models.Post
		AvailableSubscriptions := make([]uuid.UUID, 0)
		err = rows.Scan(&post.Id, &post.Creation, &post.Title,
			&post.Text, pq.Array(&post.Attachments), pq.Array(&AvailableSubscriptions))
		if err != nil {
			return models.CreatorPage{}, errors.New("InternalError")
		}

		for _, v1 := range AvailableSubscriptions {
			for _, v2 := range userSubscriptions {
				if v1 == v2 {
					post.IsAvailable = true
					break
				}
			}
		}
		if creatorPage.IsMyPage {
			post.IsAvailable = true
		}
		posts = append(posts, post)
	}
	creatorPage.Posts = make([]models.Post, len(posts))
	copy(creatorPage.Posts, posts)
	return creatorPage, nil
}
