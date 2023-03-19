package repo

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	CreatorInfo       = `SELECT user_id, name, cover_photo, followers_count, description, posts_count FROM "creator" WHERE creator_id=$1;`
	CreatorPosts      = `SELECT "post".post_id, creation_date, title, post_text, array_agg(attachment_path), array_agg(subscription_id) FROM "post" LEFT JOIN "attachment" a on "post".post_id = a.post_id JOIN "post_subscription" ps on "post".post_id = ps.post_id WHERE creator_id = $1 GROUP BY "post".post_id, creation_date, title, post_text ORDER BY creation_date DESC;`
	UserSubscriptions = `SELECT subscription_id FROM "user_subscription" WHERE user_id=$1;`
	IsLiked           = `SELECT post_id, user_id FROM "like_post" WHERE post_id = $1 AND user_id = $2`
)

type CreatorRepo struct {
	db *sql.DB
}

func NewCreatorRepo(db *sql.DB) *CreatorRepo {
	return &CreatorRepo{db: db}
}

func (ur *CreatorRepo) GetPage(userId uuid.UUID, creatorId uuid.UUID) (models.CreatorPage, error) {
	var creatorPage models.CreatorPage
	creatorPage.CreatorInfo.Id = creatorId
	creatorPage.Posts = make([]models.Post, 0)

	userSubscriptions := make([]uuid.UUID, 0)
	row := ur.db.QueryRow(CreatorInfo, creatorId)

	if err := row.Scan(&creatorPage.CreatorInfo.UserId, &creatorPage.CreatorInfo.Name, &creatorPage.CreatorInfo.CoverPhoto,
		&creatorPage.CreatorInfo.FollowersCount, &creatorPage.CreatorInfo.Description, &creatorPage.CreatorInfo.PostsCount); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return models.CreatorPage{}, models.InternalError
	} else if err == nil { //нашёл такого автора
		if creatorPage.CreatorInfo.UserId == userId { // страница автора принадлежит пользователю
			creatorPage.IsMyPage = true
		} else { // находим подписки пользователя
			row := ur.db.QueryRow(UserSubscriptions, userId)
			if err := row.Scan(pq.Array(&userSubscriptions)); err != nil && !errors.Is(sql.ErrNoRows, err) {
				return models.CreatorPage{}, models.InternalError
			}
		}
		// смотрим, какие посты доступны пользователю, исходя из его уровня подписки
		rows, err := ur.db.Query(CreatorPosts, creatorId)
		if err != nil && !errors.Is(sql.ErrNoRows, err) {
			return models.CreatorPage{}, models.InternalError
		}
		defer rows.Close()
		for rows.Next() {
			var post models.Post
			availableSubscriptions := make([]uuid.UUID, 0)
			post.Creator = creatorId
			err = rows.Scan(&post.Id, &post.Creation, &post.Title,
				&post.Text, pq.Array(&post.Attachments), pq.Array(&availableSubscriptions)) //подписки, при которыз пост доступен
			if err != nil {
				return models.CreatorPage{}, models.InternalError
			}

			if creatorPage.IsMyPage {
				post.IsAvailable = true
			}

			for _, availableSubscription := range availableSubscriptions {
				for _, userSubscription := range userSubscriptions {
					if availableSubscription == userSubscription {
						post.IsAvailable = true
						//проверяем, лайкнул ли его пользователь
						row := ur.db.QueryRow(IsLiked, post.Id, userId)
						if err = row.Scan(&post.Id, &userId); err != nil && !errors.Is(sql.ErrNoRows, err) {
							return models.CreatorPage{}, models.InternalError
						} else if err == nil {
							post.IsLiked = true
						}
						break
					}
				}
				if post.IsAvailable {
					break
				}
			}

			creatorPage.Posts = append(creatorPage.Posts, post)
		}

		return creatorPage, nil
	}
	return models.CreatorPage{}, models.WrongData // такого автора нет
}
