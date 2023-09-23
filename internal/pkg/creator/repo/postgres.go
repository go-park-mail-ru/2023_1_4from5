package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

const (
	CreatorInfo             = `SELECT user_id, name, cover_photo, followers_count, description, posts_count, aim, money_got, money_needed, profile_photo FROM "creator" WHERE creator_id=$1;`
	GetCreatorSubs          = `SELECT subscription_id, month_cost, title, description, is_available FROM "subscription" WHERE creator_id=$1;`
	GetAllCreators          = `SELECT creator_id, user_id, name, cover_photo, followers_count, description, posts_count, profile_photo FROM "creator" LIMIT 100;`
	CreatorPosts            = `SELECT "post".post_id,creation_date,title,post_text,likes_count,comments_count,array_agg(attachment_id),array_agg(attachment_type),subs FROM "post" LEFT JOIN "attachment" a on "post".post_id = a.post_id LEFT JOIN ( SELECT array_agg(DISTINCT subscription_id) as "subs", post_id FROM post_subscription s GROUP BY s.post_id ) as b on "post".post_id = b.post_id WHERE creator_id = $1 GROUP BY "post".post_id, creation_date, title, post_text, subs ORDER BY creation_date DESC;`
	UserSubscriptions       = `SELECT array_agg(subscription_id) FROM "user_subscription" WHERE user_id=$1;`
	IsLiked                 = `SELECT post_id, user_id FROM "like_post" WHERE post_id = $1 AND user_id = $2`
	GetSubInfo              = `SELECT creator_id, month_cost, title, description FROM "subscription" WHERE subscription_id = $1;`
	AddAim                  = `UPDATE creator SET aim = $1,  money_got = $2, money_needed = $3 WHERE creator_id = $4;`
	CheckIfFollow           = `SELECT user_id FROM "follow" WHERE user_id = $1 AND creator_id = $2;`
	FindCreators            = `SELECT creator_id, user_id, name, cover_photo, followers_count, description, posts_count, profile_photo FROM creator WHERE (make_tsvector(name, 'A'::"char") || make_tsvector(description, 'B'::"char")) @@ (plainto_tsquery('ru', $1) || plainto_tsquery('english', $1)) or LOWER(name) like LOWER('%' || $1 || '%') or LOWER(description) like LOWER('%' || $1 || '%') ORDER BY make_tsrank(name, $1, 'ru'::regconfig), make_tsrank(description, $1, 'ru'::regconfig) DESC LIMIT 30;`
	CheckIfCreator          = `SELECT creator_id FROM "creator" WHERE user_id = $1`
	UpdateCreatorData       = `UPDATE creator SET name = $1, description = $2 WHERE creator_id = $3`
	Feed                    = `SELECT t.post_id, t.creator_id, creation_date, title, post_text, array_agg(attachment_id), array_agg(attachment_type), t.name, t.profile_photo, t.likes_count, t.comments_count FROM (SELECT DISTINCT p.post_id, p.creator_id, creation_date, title, post_text, c.name, c.profile_photo, p.likes_count, p.comments_count FROM follow f JOIN post p on p.creator_id = f.creator_id JOIN creator c on f.creator_id = c.creator_id LEFT JOIN post_subscription ps on p.post_id = ps.post_id LEFT JOIN user_subscription us on f.user_id = us.user_id and (ps.subscription_id = us.subscription_id or ps.subscription_id is null) WHERE f.user_id = $1 AND ((ps.subscription_id is null) OR (us.subscription_id is not null)) GROUP BY c.name, p.creator_id, creation_date, title, post_text, p.post_id, c.profile_photo, c.creator_id LIMIT 50) as t LEFT JOIN attachment a on a.post_id = t.post_id GROUP BY t.name, t.creator_id, creation_date, title, post_text, t.post_id, t.profile_photo, t.likes_count, t.comments_count ORDER BY creation_date DESC;`
	UpdateProfilePhoto      = `UPDATE "creator" SET profile_photo = $1 WHERE creator_id = $2;`
	UpdateCoverPhoto        = `UPDATE "creator" SET cover_photo = $1 WHERE creator_id = $2;`
	DeleteCoverPhoto        = `UPDATE "creator" SET cover_photo = null WHERE creator_id = $1`
	DeleteProfilePhoto      = `UPDATE "creator" SET profile_photo = null WHERE creator_id = $1`
	GetStatistics           = `SELECT coalesce(sum(posts_per_month), 0), coalesce(sum(subscriptions_bought), 0), coalesce(sum(donations_count), 0), coalesce(sum(money_from_donations), 0), coalesce(sum(money_from_subscriptions),0), coalesce(sum(new_followers), 0), coalesce(sum(likes_count), 0), coalesce(sum(comments_count), 0) FROM "statistics" AS s WHERE creator_id = $1 AND  date_trunc('month'::text, s.month::date)::date BETWEEN date_trunc('month'::text, $2::date)::date AND  date_trunc('month'::text, $3::date)::date;`
	CreatorNotificationInfo = `SELECT profile_photo, name FROM creator WHERE creator_id = $1;`
	FirstStatisticsDate     = `SELECT MIN(month) FROM statistics WHERE creator_id = $1;`
	CreatorBalance          = `SELECT balance FROM creator WHERE creator_id = $1;`
	UpdateBalance           = `UPDATE creator SET balance = balance - $1 WHERE creator_id = $2 RETURNING balance;`
)

type CreatorRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewCreatorRepo(db *sql.DB, logger *zap.SugaredLogger) *CreatorRepo {
	return &CreatorRepo{
		db:     db,
		logger: logger,
	}
}

func (r *CreatorRepo) CheckIfFollow(ctx context.Context, userId, creatorId uuid.UUID) (bool, error) {
	row := r.db.QueryRowContext(ctx, CheckIfFollow, userId, creatorId)
	if err := row.Scan(&userId); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return false, models.InternalError
	} else if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, nil
}

func (r *CreatorRepo) CreatorNotificationInfo(ctx context.Context, creatorID uuid.UUID) (models.NotificationCreatorInfo, error) {
	var info models.NotificationCreatorInfo
	row := r.db.QueryRowContext(ctx, CreatorNotificationInfo, creatorID)
	if err := row.Scan(&info.Photo, &info.Name); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.NotificationCreatorInfo{}, models.InternalError
	}
	return info, nil
}

func (ur *CreatorRepo) UpdateBalance(ctx context.Context, transfer models.CreatorTransfer) (float32, error) {
	var newBalance float32
	row := ur.db.QueryRowContext(ctx, UpdateBalance, transfer.Money, transfer.CreatorID)
	if err := row.Scan(&newBalance); err != nil && !errors.Is(err, sql.ErrNoRows) {
		ur.logger.Error(err)
		return 0, models.InternalError
	}
	return newBalance, nil
}

func (r *CreatorRepo) GetUserSubscriptions(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	userSubscriptions := make([]uuid.UUID, 0)
	row := r.db.QueryRowContext(ctx, UserSubscriptions, userId)
	if err := row.Scan(pq.Array(&userSubscriptions)); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return nil, models.InternalError
	}
	return userSubscriptions, nil
}

func (r *CreatorRepo) StatisticsFirstDate(ctx context.Context, creatorID uuid.UUID) (string, error) {
	var firstDate string
	row := r.db.QueryRowContext(ctx, FirstStatisticsDate, creatorID)
	if err := row.Scan(&firstDate); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return "", models.InternalError
	}
	return firstDate, nil
}

func (r *CreatorRepo) GetCreatorBalance(ctx context.Context, creatorID uuid.UUID) (float32, error) {
	var balance float32
	row := r.db.QueryRowContext(ctx, CreatorBalance, creatorID)
	if err := row.Scan(&balance); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return 0, models.InternalError
	}
	return balance, nil
}

func (r *CreatorRepo) IsLiked(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (bool, error) {
	row := r.db.QueryRowContext(ctx, IsLiked, postID, userID)
	if err := row.Scan(&postID, &userID); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return false, models.InternalError
	} else if err == nil {
		return true, nil
	}
	return false, nil
}

func (r *CreatorRepo) CreatorInfo(ctx context.Context, creatorPage *models.CreatorPage, creatorID uuid.UUID) error {
	row := r.db.QueryRowContext(ctx, CreatorInfo, creatorID)
	creatorPage.CreatorInfo.Id = creatorID
	var tmpAim sql.NullString
	if err := row.Scan(&creatorPage.CreatorInfo.UserId, &creatorPage.CreatorInfo.Name, &creatorPage.CreatorInfo.CoverPhoto,
		&creatorPage.CreatorInfo.FollowersCount, &creatorPage.CreatorInfo.Description, &creatorPage.CreatorInfo.PostsCount,
		&tmpAim, &creatorPage.Aim.MoneyGot, &creatorPage.Aim.MoneyNeeded, &creatorPage.CreatorInfo.ProfilePhoto); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.NotFound
	}
	creatorPage.Aim.Creator = creatorID
	creatorPage.Aim.Description = tmpAim.String
	return nil
}

func (r *CreatorRepo) GetCreatorSubs(ctx context.Context, creatorID uuid.UUID) ([]models.Subscription, error) {
	subs := make([]models.Subscription, 0)
	var tmpTitle, tmpDescr sql.NullString
	rows, err := r.db.QueryContext(ctx, GetCreatorSubs, creatorID)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return nil, models.InternalError
	}
	defer rows.Close()
	for rows.Next() {
		tmpSub := models.Subscription{}
		var isAvailable bool
		err = rows.Scan(&tmpSub.Id, &tmpSub.MonthCost, &tmpTitle, &tmpDescr, &isAvailable)
		if err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
		}
		if !isAvailable {
			continue
		}
		tmpSub.Title = tmpTitle.String
		tmpSub.Description = tmpDescr.String
		tmpSub.Creator = creatorID

		subs = append(subs, tmpSub)
	}
	return subs, nil
}

func (r *CreatorRepo) CreatorPosts(ctx context.Context, creatorId uuid.UUID) ([]models.Post, error) {
	var posts = make([]models.Post, 0)
	rows, err := r.db.QueryContext(ctx, CreatorPosts, creatorId)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return nil, models.InternalError
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		availableSubscriptions := make([]uuid.UUID, 0)
		post.Creator = creatorId
		attachs := make([]uuid.UUID, 0)
		types := make([]sql.NullString, 0)
		err = rows.Scan(&post.Id, &post.Creation, &post.Title,
			&post.Text, &post.LikesCount, &post.CommentsCount, pq.Array(&attachs), pq.Array(&types), pq.Array(&availableSubscriptions)) //подписки, при которыз пост доступен
		if err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
		}
		post.Subscriptions = make([]models.Subscription, len(availableSubscriptions))
		if post.Subscriptions, err = r.GetSubsByID(ctx, availableSubscriptions...); err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
		}
		post.Attachments = make([]models.Attachment, 0, len(attachs))
		for i, v := range attachs {
			if v != uuid.Nil {
				post.Attachments = append(post.Attachments, models.Attachment{Id: v, Type: types[i].String})
			}
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (r *CreatorRepo) GetPage(ctx context.Context, userId uuid.UUID, creatorId uuid.UUID) (models.CreatorPage, error) {
	var creatorPage models.CreatorPage
	creatorPage.Posts = make([]models.Post, 0)
	var userSubscriptions []uuid.UUID
	if err := r.CreatorInfo(ctx, &creatorPage, creatorId); err == models.InternalError {
		return models.CreatorPage{}, models.InternalError
	} else if err == nil { //нашёл такого автора
		if creatorPage.Follows, err = r.CheckIfFollow(ctx, userId, creatorId); err != nil {
			return models.CreatorPage{}, models.InternalError
		}
		if creatorPage.CreatorInfo.UserId == userId { // страница автора принадлежит пользователю
			creatorPage.IsMyPage = true
		} else { // находим подписки пользователя
			tmp, err := r.GetUserSubscriptions(ctx, userId)
			if err != nil {
				fmt.Println("user subs")
				r.logger.Error(err)
				return models.CreatorPage{}, models.InternalError
			}
			userSubscriptions = make([]uuid.UUID, len(tmp))
			copy(userSubscriptions, tmp)
		}
		creatorPage.Posts, err = r.CreatorPosts(ctx, creatorId)
		if err != nil {
			return models.CreatorPage{}, err
		}

		for i := range creatorPage.Posts {
			if creatorPage.IsMyPage {
				creatorPage.Posts[i].IsAvailable = true
			}
			if len(creatorPage.Posts[i].Subscriptions) == 0 {
				creatorPage.Posts[i].IsAvailable = true
			}
			for _, availableSubscription := range creatorPage.Posts[i].Subscriptions {
				for _, userSubscription := range userSubscriptions {
					if availableSubscription.Id == userSubscription {
						creatorPage.Posts[i].IsAvailable = true
						break
					}
				}
				if creatorPage.Posts[i].IsAvailable {
					break
				}
			}
			if creatorPage.Posts[i].IsLiked, err = r.IsLiked(ctx, userId, creatorPage.Posts[i].Id); err != nil {
				fmt.Println("is liked")
				return models.CreatorPage{}, models.InternalError
			}
			if !creatorPage.Posts[i].IsAvailable {
				creatorPage.Posts[i].Text = ""
				creatorPage.Posts[i].Attachments = nil
			}
		}

		if creatorPage.Subscriptions, err = r.GetCreatorSubs(ctx, creatorId); err != nil {
			return models.CreatorPage{}, err
		}

		return creatorPage, nil
	}
	return models.CreatorPage{}, models.WrongData // такого автора нет
}

func (r *CreatorRepo) GetSubsByID(ctx context.Context, subsIDs ...uuid.UUID) ([]models.Subscription, error) {
	subsInfo := make([]models.Subscription, 0)
	var sub models.Subscription
	for _, v := range subsIDs {
		row := r.db.QueryRowContext(ctx, GetSubInfo, v)
		err := row.Scan(&sub.Creator, &sub.MonthCost, &sub.Title,
			&sub.Description)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			r.logger.Error(err)
			return nil, models.InternalError
		} else if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		sub.Id = v
		subsInfo = append(subsInfo, sub)
	}
	return subsInfo, nil
}

func (r *CreatorRepo) CreateAim(ctx context.Context, aimInfo models.Aim) error {
	row := r.db.QueryRowContext(ctx, AddAim, aimInfo.Description, aimInfo.MoneyGot, aimInfo.MoneyNeeded, aimInfo.Creator)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (r *CreatorRepo) GetAllCreators(ctx context.Context) ([]models.Creator, error) {
	var creators = make([]models.Creator, 0)
	rows, err := r.db.QueryContext(ctx, GetAllCreators)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return nil, models.InternalError
	}
	defer rows.Close()
	for rows.Next() {
		var creator models.Creator
		var tmpDescr sql.NullString
		err = rows.Scan(&creator.Id, &creator.UserId, &creator.Name,
			&creator.CoverPhoto, &creator.FollowersCount, &tmpDescr, &creator.PostsCount, &creator.ProfilePhoto)
		if err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
		}
		creator.Description = tmpDescr.String
		creators = append(creators, creator)
	}

	return creators, nil
}

func (r *CreatorRepo) FindCreators(ctx context.Context, keyword string) ([]models.Creator, error) {
	var creators = make([]models.Creator, 0)
	rows, err := r.db.QueryContext(ctx, FindCreators, keyword)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return nil, models.InternalError
	}

	defer rows.Close()
	for rows.Next() {
		var creator models.Creator
		var tmpDescr sql.NullString
		err = rows.Scan(&creator.Id, &creator.UserId, &creator.Name,
			&creator.CoverPhoto, &creator.FollowersCount, &tmpDescr, &creator.PostsCount, &creator.ProfilePhoto)
		if err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
		}
		creator.Description = tmpDescr.String
		creators = append(creators, creator)
	}

	return creators, nil
}

func (r *CreatorRepo) UpdateCreatorData(ctx context.Context, updateData models.UpdateCreatorInfo) error {
	row := r.db.QueryRowContext(ctx, UpdateCreatorData, updateData.CreatorName, updateData.Description, updateData.CreatorID)
	if err := row.Scan(); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (r *CreatorRepo) CheckIfCreator(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var creatorID uuid.UUID
	row := r.db.QueryRowContext(ctx, CheckIfCreator, userID)
	if err := row.Scan(&creatorID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return uuid.Nil, models.InternalError
	} else if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, models.NotFound
	}
	return creatorID, nil
}

func (r *CreatorRepo) GetFeed(ctx context.Context, userID uuid.UUID) ([]models.Post, error) {
	var feed = make([]models.Post, 0)

	rows, err := r.db.QueryContext(ctx, Feed, userID)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return nil, models.InternalError
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		attachs := make([]uuid.UUID, 0)
		types := make([]sql.NullString, 0)
		err = rows.Scan(&post.Id, &post.Creator, &post.Creation,
			&post.Title, &post.Text, pq.Array(&attachs), pq.Array(&types), &post.CreatorName, &post.CreatorPhoto, &post.LikesCount, &post.CommentsCount)
		if err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
		}

		post.IsAvailable = true

		if post.IsLiked, err = r.IsLiked(ctx, userID, post.Id); err != nil {
			return nil, models.InternalError
		}

		post.Attachments = make([]models.Attachment, 0, len(attachs))
		for i, v := range attachs {
			if v != uuid.Nil {
				post.Attachments = append(post.Attachments, models.Attachment{Id: v, Type: types[i].String})
			}
		}

		feed = append(feed, post)
	}

	return feed, nil
}

func (r *CreatorRepo) UpdateProfilePhoto(ctx context.Context, creatorId, path uuid.UUID) error {
	row := r.db.QueryRowContext(ctx, UpdateProfilePhoto, path, creatorId)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (r *CreatorRepo) UpdateCoverPhoto(ctx context.Context, creatorId, path uuid.UUID) error {
	row := r.db.QueryRowContext(ctx, UpdateCoverPhoto, path, creatorId)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (r *CreatorRepo) DeleteCoverPhoto(ctx context.Context, creatorId uuid.UUID) error {
	row := r.db.QueryRowContext(ctx, DeleteCoverPhoto, creatorId)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (r *CreatorRepo) DeleteProfilePhoto(ctx context.Context, creatorId uuid.UUID) error {
	row := r.db.QueryRowContext(ctx, DeleteProfilePhoto, creatorId)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (r *CreatorRepo) Statistics(ctx context.Context, statsInput models.StatisticsDates) (models.Statistics, error) {
	var stat models.Statistics

	row := r.db.QueryRowContext(ctx, GetStatistics, statsInput.CreatorId, statsInput.FirstMonth.Format(time.RFC3339), statsInput.SecondMonth.Format(time.RFC3339))
	err := row.Scan(&stat.PostsPerMonth, &stat.SubscriptionsBought, &stat.DonationsCount, &stat.MoneyFromDonations, &stat.MoneyFromSubscriptions, &stat.NewFollowers, &stat.LikesCount, &stat.CommentsCount)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return models.Statistics{}, models.WrongData
	}
	if err != nil {
		fmt.Println(err)
		r.logger.Error(err)
		return models.Statistics{}, models.InternalError
	}

	return stat, nil
}
