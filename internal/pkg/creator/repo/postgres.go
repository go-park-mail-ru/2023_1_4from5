package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	CreatorInfo       = `SELECT user_id, name, cover_photo, followers_count, description, posts_count, aim, money_got, money_needed, profile_photo FROM "creator" WHERE creator_id=$1;`
	GetCreatorSubs    = `SELECT subscription_id, month_cost, title, description FROM "subscription" WHERE creator_id=$1;`
	GetAllCreators    = `SELECT creator_id, user_id, name, cover_photo, followers_count, description, posts_count, profile_photo FROM "creator" LIMIT 100;`
	CreatorPosts      = `SELECT "post".post_id, creation_date, title, post_text, likes_count, array_agg(attachment_id), array_agg(attachment_type), array_agg(DISTINCT subscription_id) FROM "post" LEFT JOIN "attachment" a on "post".post_id = a.post_id LEFT JOIN "post_subscription" ps on "post".post_id = ps.post_id WHERE creator_id = $1 GROUP BY "post".post_id, creation_date, title, post_text ORDER BY creation_date DESC;`
	UserSubscriptions = `SELECT array_agg(subscription_id) FROM "user_subscription" WHERE user_id=$1;`
	IsLiked           = `SELECT post_id, user_id FROM "like_post" WHERE post_id = $1 AND user_id = $2`
	GetSubInfo        = `SELECT creator_id, month_cost, title, description FROM "subscription" WHERE subscription_id = $1;`
	AddAim            = `UPDATE creator SET aim = $1,  money_got = $2, money_needed = $3 WHERE creator_id = $4;`
	CheckIfFollow     = `SELECT user_id FROM "follow" WHERE user_id = $1 AND creator_id = $2;`
	FindCreators      = `SELECT creator_id, user_id, name, cover_photo, followers_count, description, posts_count, profile_photo FROM creator WHERE (setweight(to_tsvector('russian', name),'A') || setweight(to_tsvector('russian', description), 'B') || setweight(to_tsvector('english', name),'A') || setweight(to_tsvector('english', description), 'B'))  @@ (plainto_tsquery('russian',$1)|| plainto_tsquery('english',$1)) LIMIT 30;`
	CheckIfCreator    = `SELECT creator_id FROM "creator" WHERE user_id = $1`
	UpdateCreatorData = `UPDATE creator SET name = $1, description = $2 WHERE creator_id = $3`
	Feed              = `SELECT DISTINCT p.post_id, p.creator_id, creation_date, title, post_text, array_agg(attachment_id), array_agg(attachment_type), c.name, c.profile_photo, c.creator_id, p.likes_count FROM follow f JOIN post p on p.creator_id = f.creator_id JOIN creator c on f.creator_id = c.creator_id LEFT JOIN post_subscription ps on p.post_id = ps.post_id JOIN user_subscription us on f.user_id = us.user_id and (ps.subscription_id = us.subscription_id or ps.subscription_id is null) LEFT JOIN "attachment" a on p.post_id = a.post_id WHERE f.user_id = $1 GROUP BY c.name, p.creator_id, creation_date, title, post_text, p.post_id, c.profile_photo, c.creator_id ORDER BY creation_date DESC LIMIT 50;`
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

func (r *CreatorRepo) GetUserSubscriptions(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	userSubscriptions := make([]uuid.UUID, 0)
	row := r.db.QueryRowContext(ctx, UserSubscriptions, userId)
	if err := row.Scan(pq.Array(&userSubscriptions)); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return nil, models.InternalError
	}
	return userSubscriptions, nil
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
		err = rows.Scan(&tmpSub.Id, &tmpSub.MonthCost, &tmpTitle, &tmpDescr)
		if err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
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
			&post.Text, &post.LikesCount, pq.Array(&attachs), pq.Array(&types), pq.Array(&availableSubscriptions)) //подписки, при которыз пост доступен
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
			&post.Title, &post.Text, pq.Array(&attachs), pq.Array(&types), &post.CreatorName, &post.CreatorPhoto, &post.Creator, &post.LikesCount)
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
