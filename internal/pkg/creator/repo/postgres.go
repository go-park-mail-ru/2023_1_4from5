package repo

import (
	"database/sql"
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
	CREATOR_POSTS     = "SELECT post_id, creation_date, title, post_text, attachments, available_subscriptions  FROM public.post WHERE creator_id=$1 ORDER BY creation_date DESC;"
	USER_SUBSCRIPTION = "SELECT subscriptions FROM public.user WHERE user_id=$1;"
)

func (ur *CreatorRepo) GetPage(userId uuid.UUID, creatorId uuid.UUID) (models.CreatorPage, error) {
	var creatorPage models.CreatorPage
	creatorPage.CreatorInfo.Id = creatorId

	userSubscriptions := make([]uuid.UUID, 0)
	row := ur.db.QueryRow(CREATOR_INFO, creatorId)
	if err := row.Scan(&creatorPage.CreatorInfo.UserId, &creatorPage.CreatorInfo.Name, &creatorPage.CreatorInfo.CoverPhoto, &creatorPage.CreatorInfo.FollowersCount, &creatorPage.CreatorInfo.Description, &creatorPage.CreatorInfo.PostsCount); err != nil {
		return models.CreatorPage{}, models.InternalError
	}
	if creatorPage.CreatorInfo.UserId == userId {
		creatorPage.IsMyPage = true
	} else {
		row := ur.db.QueryRow(USER_SUBSCRIPTION, userId)
		if err := row.Scan(pq.Array(&userSubscriptions)); err != nil {
			return models.CreatorPage{}, models.InternalError
		}
	}

	rows, err := ur.db.Query(CREATOR_POSTS, creatorId)
	defer rows.Close()
	if err != nil {
		return models.CreatorPage{}, models.InternalError
	}
	posts := make([]models.Post, 0)
	for rows.Next() {
		var post models.Post
		availableSubscriptions := make([]uuid.UUID, 0)
		post.Creator = creatorId
		err = rows.Scan(&post.Id, &post.Creation, &post.Title,
			&post.Text, pq.Array(&post.Attachments), pq.Array(&availableSubscriptions))
		if err != nil {
			return models.CreatorPage{}, models.InternalError
		}

		if creatorPage.IsMyPage {
			post.IsAvailable = true
		}

		for _, v1 := range availableSubscriptions {
			for _, v2 := range userSubscriptions {
				if v1 == v2 {
					post.IsAvailable = true
					break
				}
			}
			if post.IsAvailable {
				break
			}
		}

		posts = append(posts, post)
	}

	creatorPage.Posts = make([]models.Post, len(posts))
	copy(creatorPage.Posts, posts)
	return creatorPage, nil
}
