package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	CreatorInfo       = `SELECT user_id, name, cover_photo, followers_count, description, posts_count FROM "creator" WHERE creator_id=$1;`
	CreatorPosts      = `SELECT "post".post_id, creation_date, title, post_text, array_agg(attachment_id), array_agg(subscription_id) FROM "post" LEFT JOIN "attachment" a on "post".post_id = a.post_id JOIN "post_subscription" ps on "post".post_id = ps.post_id WHERE creator_id = $1 GROUP BY "post".post_id, creation_date, title, post_text ORDER BY creation_date DESC;`
	UserSubscriptions = `SELECT subscription_id FROM "user_subscription" WHERE user_id=$1;`
	IsLiked           = `SELECT post_id, user_id FROM "like_post" WHERE post_id = $1 AND user_id = $2`
)

type CreatorRepo struct {
	db *sql.DB
}

func NewCreatorRepo(db *sql.DB) *CreatorRepo {
	return &CreatorRepo{db: db}
}

func (r *CreatorRepo) GetUserSubscriptions(userId uuid.UUID) ([]uuid.UUID, error) {
	userSubscriptions := make([]uuid.UUID, 0)
	row := r.db.QueryRow(UserSubscriptions, userId)
	if err := row.Scan(pq.Array(&userSubscriptions)); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return nil, models.InternalError
	}
	return userSubscriptions, nil
}

func (r *CreatorRepo) IsLiked(userID uuid.UUID, postID uuid.UUID) (bool, error) {
	row := r.db.QueryRow(IsLiked, postID, userID)
	if err := row.Scan(&postID, &userID); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return false, models.InternalError
	} else if err == nil {
		return true, nil
	}
	return false, nil
}

func (r *CreatorRepo) CreatorInfo(creatorPage *models.CreatorPage, creatorID uuid.UUID) error {
	row := r.db.QueryRow(CreatorInfo, creatorID)
	if err := row.Scan(&creatorPage.CreatorInfo.UserId, &creatorPage.CreatorInfo.Name, &creatorPage.CreatorInfo.CoverPhoto,
		&creatorPage.CreatorInfo.FollowersCount, &creatorPage.CreatorInfo.Description, &creatorPage.CreatorInfo.PostsCount); err != nil && !errors.Is(sql.ErrNoRows, err) {
		fmt.Println(err)
		return models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.NotFound
	}
	return nil
}

func (r *CreatorRepo) GetPage(userId uuid.UUID, creatorId uuid.UUID) (models.CreatorPage, error) {
	var creatorPage models.CreatorPage
	creatorPage.CreatorInfo.Id = creatorId
	creatorPage.Posts = make([]models.Post, 0)
	userSubscriptions := make([]uuid.UUID, 0)

	if err := r.CreatorInfo(&creatorPage, creatorId); err == models.InternalError {
		fmt.Println("creatorInfo")
		return models.CreatorPage{}, models.InternalError
	} else if err == nil { //нашёл такого автора
		if creatorPage.CreatorInfo.UserId == userId { // страница автора принадлежит пользователю
			creatorPage.IsMyPage = true
		} else { // находим подписки пользователя
			tmp, err := r.GetUserSubscriptions(userId)
			copy(userSubscriptions, tmp)
			if err != nil {
				fmt.Println("userSubs")
				return models.CreatorPage{}, models.InternalError
			}
		}
		// смотрим, какие посты доступны пользователю, исходя из его уровня подписки
		rows, err := r.db.Query(CreatorPosts, creatorId)
		if err != nil && !errors.Is(sql.ErrNoRows, err) {
			fmt.Println(err)
			fmt.Println("Posts")
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
				fmt.Println("Read Post")
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
						if post.IsLiked, err = r.IsLiked(userId, post.Id); err != nil {
							fmt.Println("IsLiked")
							return models.CreatorPage{}, models.InternalError
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
