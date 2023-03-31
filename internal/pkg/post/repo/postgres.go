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
	InsertPost             = `INSERT INTO "post"(post_id, creator_id, title, post_text) VALUES($1, $2, $3, $4);`
	InsertAttach           = `INSERT INTO "attachment"(attachment_id, post_id, attachment_type) VALUES($1, $2, $3);`
	IncPostCount           = `UPDATE "creator" SET posts_count = posts_count+1 WHERE creator_id = $1;`
	AddSubscriptionsToPost = `INSERT INTO "post_subscription"(post_id, subscription_id) VALUES($1,$2);`
	DeletePost             = `DELETE FROM  "post" WHERE post_id = $1;`
	GetUserId              = `SELECT user_id FROM "post" JOIN "creator" c on c.creator_id = "post".creator_id WHERE post_id = $1`
	AddLike                = `INSERT INTO "like_post"(post_id, user_id) VALUES($1, $2);`
	RemoveLike             = `DELETE FROM "like_post" WHERE post_id = $1 AND user_id = $2;`
	UpdateLikeCount        = `UPDATE "post" SET likes_count = likes_count + $1 WHERE post_id = $2 RETURNING likes_count;`
	IsLiked                = `SELECT post_id, user_id FROM "like_post" WHERE post_id = $1 AND user_id = $2;`
	IsPostAvailable        = `SELECT user_id FROM "user_subscription" INNER JOIN "post_subscription" p on "user_subscription".subscription_id = p.subscription_id WHERE user_id = $1 AND post_id = $2 AND expire_date > now()`
	IsCreator              = `SELECT user_id FROM "creator" WHERE creator_id = $1;`
	GetPost                = `SELECT "post".post_id, "post".creator_id, creation_date, title, post_text, array_agg(attachment_id), array_agg(attachment_type), array_agg(DISTINCT subscription_id) FROM "post" JOIN "attachment" a on "post".post_id = a.post_id LEFT JOIN "post_subscription" ps on "post".post_id = ps.post_id WHERE "post".post_id = $1 GROUP BY "post".post_id, creation_date, title, post_text;`
	GetSubInfo             = `SELECT creator_id, month_cost, title, description FROM "subscription" WHERE subscription_id = $1;`
)

type PostRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewPostRepo(db *sql.DB, logger *zap.SugaredLogger) *PostRepo {
	return &PostRepo{
		db:     db,
		logger: logger,
	}
}

func (r *PostRepo) IsPostAvailable(ctx context.Context, userID, postID uuid.UUID) error {
	var userIDtmp uuid.UUID
	row := r.db.QueryRow(IsPostAvailable, userID, postID)
	if err := row.Scan(&userIDtmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.WrongData
	}
	return nil
}

func (r *PostRepo) CreatePost(ctx context.Context, postData models.PostCreationData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	row, err := tx.QueryContext(ctx, InsertPost, postData.Id, postData.Creator, postData.Title, postData.Text)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	if err = row.Close(); err != nil {
		return models.InternalError
	}

	for _, attach := range postData.Attachments {
		if row, err = tx.QueryContext(ctx, InsertAttach, attach.Id, postData.Id, attach.Type); err != nil {
			_ = tx.Rollback()
			r.logger.Error(err)
			return models.InternalError
		}
		if err = row.Close(); err != nil {
			r.logger.Error(err)
			return models.InternalError
		}
	}

	if row, err = tx.QueryContext(ctx, IncPostCount, postData.Creator); err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}
	if err = row.Close(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	for _, sub := range postData.AvailableSubscriptions {
		if row, err = tx.QueryContext(ctx, AddSubscriptionsToPost, postData.Id, sub); err != nil {
			_ = tx.Rollback()
			r.logger.Error(err)
			return models.InternalError
		}
		if err = row.Close(); err != nil {
			r.logger.Error(err)
			return models.InternalError
		}
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	return nil
}

func (r *PostRepo) GetSubsByID(ctx context.Context, subsIDs ...uuid.UUID) ([]models.Subscription, error) {
	subsInfo := make([]models.Subscription, len(subsIDs))
	for i, v := range subsIDs {
		row := r.db.QueryRow(GetSubInfo, v)
		err := row.Scan(&subsInfo[i].Creator, &subsInfo[i].MonthConst, &subsInfo[i].Title,
			&subsInfo[i].Description)
		if err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
		}
		subsInfo[i].Id = v
	}
	return subsInfo, nil
}

func (r *PostRepo) GetPost(ctx context.Context, postID, userID uuid.UUID) (models.Post, error) {
	var post models.Post
	attachs := make([]uuid.UUID, 0)
	types := make([]sql.NullString, 0)
	subs := make([]uuid.UUID, 0)
	row := r.db.QueryRow(GetPost, postID)
	err := row.Scan(&post.Id, &post.Creator, &post.Creation, &post.Title,
		&post.Text, pq.Array(&attachs), pq.Array(&types), pq.Array(&subs)) //подписки, при которыз пост доступен
	if err != nil {
		r.logger.Error(err)
		return models.Post{}, models.InternalError
	}
	attachs = attachs[:len(attachs)/2] //TODO !!!!!!!!!!!!!!!!
	post.Attachments = make([]models.Attachment, len(attachs))
	for i, v := range attachs {
		post.Attachments[i].Type = types[i].String
		post.Attachments[i].Id = v
	}

	post.Subscriptions, err = r.GetSubsByID(ctx, subs...)
	return post, err
}

func (r *PostRepo) IsCreator(ctx context.Context, userID, creatorID uuid.UUID) (bool, error) {
	var userIdtmp uuid.UUID
	row := r.db.QueryRow(IsCreator, creatorID)
	if err := row.Scan(&userIdtmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return false, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return false, models.WrongData
	}
	if userIdtmp != userID {
		return false, nil
	}
	return true, nil
}

func (r *PostRepo) DeletePost(ctx context.Context, postID uuid.UUID) error {
	//TODO: delete also in another table (or cascade delete?)
	row := r.db.QueryRow(DeletePost, postID)
	if err := row.Err(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}
	return nil
}

func (r *PostRepo) IsPostOwner(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
	var userIdtmp uuid.UUID
	row := r.db.QueryRow(GetUserId, postId)
	if err := row.Scan(&userIdtmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return false, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return false, models.WrongData
	}
	if userIdtmp != userId {
		return false, nil
	}
	return true, nil
}

func (r *PostRepo) AddLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
	var (
		userUUID uuid.UUID
		postUUID uuid.UUID
	)
	//проверяем, лайкнул ли уже
	row := r.db.QueryRow(IsLiked, postID, userID)
	if err := row.Scan(&postUUID, &userUUID); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return models.Like{}, models.InternalError
	} else if err == nil { // уже есть запись об этом лайке
		return models.Like{}, models.WrongData
	}
	// проверяем, есть ли доступ к этому посту
	if err := r.IsPostAvailable(ctx, userID, postID); err != nil {
		r.logger.Error(err)
		return models.Like{}, err
	}
	// обновляем кол-во лайков, заодно смотрим, есть ли вообще такой пост
	var like models.Like
	like.PostID = postID
	row = r.db.QueryRow(UpdateLikeCount, 1, postID)

	if err := row.Scan(&like.LikesCount); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return models.Like{}, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.Like{}, models.WrongData
	}

	row = r.db.QueryRow(AddLike, postID, userID)

	if err := row.Err(); err != nil {
		r.logger.Error(err)
		return models.Like{}, models.InternalError
	}

	return like, nil
}

func (r *PostRepo) RemoveLike(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
	var (
		userUUID uuid.UUID
		postUUID uuid.UUID
	)
	row := r.db.QueryRow(IsLiked, postID, userID)
	if err := row.Scan(&postUUID, &userUUID); err != nil && !errors.Is(sql.ErrNoRows, err) {
		return models.Like{}, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) { // нет такого лайка
		return models.Like{}, models.WrongData
	}

	var like models.Like
	like.PostID = postID
	row = r.db.QueryRow(UpdateLikeCount, -1, postID)

	if err := row.Scan(&like.LikesCount); err != nil {
		r.logger.Error(err)
		return models.Like{}, models.InternalError
	}

	row = r.db.QueryRow(RemoveLike, postID, userID)

	if err := row.Err(); err != nil {
		r.logger.Error(err)
		return models.Like{}, models.InternalError
	}

	return like, nil
}
