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
)

const (
	InsertPost                 = `INSERT INTO "post"(post_id, creator_id, title, post_text) VALUES($1, $2, $3, $4);`
	InsertAttach               = `INSERT INTO "attachment"(attachment_id, post_id, attachment_type) VALUES($1, $2, $3);`
	IncPostCount               = `UPDATE "creator" SET posts_count = posts_count+1 WHERE creator_id = $1;`
	UpdatePostInfo             = `UPDATE "post" SET title = $1, post_text = $2 WHERE post_id = $3;`
	DeletePostSubscriptions    = `DELETE FROM "post_subscription" WHERE post_id = $1;`
	AddSubscriptionsToPost     = `INSERT INTO "post_subscription"(post_id, subscription_id) VALUES($1,$2);`
	DeletePost                 = `DELETE FROM  "post" WHERE post_id = $1;`
	DeletePostSubscription     = `DELETE FROM "post_subscription" WHERE post_id = $1`
	GetUserId                  = `SELECT user_id FROM "post" JOIN "creator" c on c.creator_id = "post".creator_id WHERE post_id = $1`
	AddLike                    = `INSERT INTO "like_post"(post_id, user_id) VALUES($1, $2);`
	RemoveLike                 = `DELETE FROM "like_post" WHERE post_id = $1 AND user_id = $2;`
	UpdateLikeCount            = `UPDATE "post" SET likes_count = likes_count + $1 WHERE post_id = $2 RETURNING likes_count;`
	IsLiked                    = `SELECT post_id, user_id FROM "like_post" WHERE post_id = $1 AND user_id = $2;`
	DeleteLikes                = `DELETE FROM "like_post" WHERE post_id = $1;`
	DeleteComments             = `DELETE FROM "comment" WHERE post_id = $1;`
	IsPostAvailableWithSub     = `SELECT user_id FROM "user_subscription" INNER JOIN "post_subscription" p on "user_subscription".subscription_id = p.subscription_id WHERE user_id = $1 AND post_id = $2 AND expire_date > now()`
	IsPostAvailableForEveryone = `SELECT post_id FROM post_subscription WHERE post_id = $1`
	IsCreator                  = `SELECT user_id FROM "creator" WHERE creator_id = $1;`
	GetPost                    = `SELECT "post".post_id, "post".creator_id, creation_date, title, post_text, likes_count, "post".comments_count, array_agg(attachment_id), array_agg(attachment_type), array_agg(DISTINCT subscription_id) FROM "post" LEFT JOIN "attachment" a on "post".post_id = a.post_id LEFT JOIN "post_subscription" ps on "post".post_id = ps.post_id WHERE "post".post_id = $1 GROUP BY "post".post_id, creation_date, title, post_text;`
	GetSubInfo                 = `SELECT creator_id, month_cost, title, description FROM "subscription" WHERE subscription_id = $1;`
	GetComments                = `SELECT comment_id, u.user_id, u.display_name, u.profile_photo, c.post_id, c.comment_text, c.creation_date, c.likes_count FROM comment c JOIN "user" u on c.user_id = u.user_id WHERE post_id = $1;`
	IsLikedComment             = `SELECT comment_id FROM "like_comment" WHERE comment_id = $1 AND user_id = $2;`
	GetUserIdComments          = `SELECT user_id FROM "comment" WHERE comment_id = $1;`
	GetCreatorPhoto            = `SELECT profile_photo FROM "creator" WHERE creator_id = $1`
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
	var postIDtmp uuid.UUID
	row := r.db.QueryRowContext(ctx, IsPostAvailableWithSub, userID, postID)
	if err := row.Scan(&userIDtmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		if ok, err := r.IsPostOwner(ctx, userID, postID); err == models.InternalError {
			return models.InternalError
		} else if ok {
			return nil
		}
		row = r.db.QueryRowContext(ctx, IsPostAvailableForEveryone, postID)
		if err = row.Scan(&postIDtmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
			r.logger.Error(err)
			return models.InternalError
		} else if errors.Is(sql.ErrNoRows, err) {
			return nil
		}
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
		if _, err = tx.ExecContext(ctx, InsertAttach, attach.Id, postData.Id, attach.Type); err != nil {
			_ = tx.Rollback()
			r.logger.Error(err)
			return models.InternalError
		}
	}

	if _, err = tx.ExecContext(ctx, IncPostCount, postData.Creator); err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	for _, sub := range postData.AvailableSubscriptions {
		if _, err = tx.ExecContext(ctx, AddSubscriptionsToPost, postData.Id, sub); err != nil {
			_ = tx.Rollback()
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
	subsInfo := make([]models.Subscription, 0)
	for i, v := range subsIDs {
		var sub = models.Subscription{}
		if v == uuid.Nil {
			subsIDs = append(subsIDs[:i], subsIDs[i+1:]...)
			continue
		}
		row := r.db.QueryRowContext(ctx, GetSubInfo, v)
		err := row.Scan(&sub.Creator, &sub.MonthCost, &sub.Title,
			&sub.Description)
		if err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
		}
		sub.Id = v
		subsInfo = append(subsInfo, sub)
	}
	return subsInfo, nil
}

func (r *PostRepo) GetComments(ctx context.Context, postID, userID uuid.UUID) ([]models.Comment, error) {
	comments := make([]models.Comment, 0)
	rows, err := r.db.QueryContext(ctx, GetComments, postID)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return nil, models.InternalError
	}
	defer rows.Close()
	for rows.Next() {
		comment := models.Comment{}

		err = rows.Scan(&comment.CommentID, &comment.UserID, &comment.Username, &comment.UserPhoto, &comment.PostID, &comment.Text, &comment.Creation, &comment.LikesCount)
		if err != nil {
			r.logger.Error(err)
			return nil, models.InternalError
		}

		//check for liked
		var commentId = uuid.UUID{}
		row := r.db.QueryRowContext(ctx, IsLikedComment, comment.CommentID, userID)
		if err := row.Scan(&commentId); err != nil && !errors.Is(sql.ErrNoRows, err) {
			r.logger.Error(err)
			return nil, models.InternalError
		} else if err == nil {
			comment.IsLiked = true
		}

		//check for owning
		isCommentOwner, err := r.IsCommentOwner(ctx, comment.CommentID, userID)
		if err != nil {
			return nil, err
		}
		comment.IsOwner = isCommentOwner

		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *PostRepo) GetPost(ctx context.Context, postID, userID uuid.UUID) (models.Post, error) {
	var post models.Post
	var postTextTmp sql.NullString
	attachs := make([]uuid.UUID, 0)
	types := make([]sql.NullString, 0)
	subs := make([]uuid.UUID, 0)
	row := r.db.QueryRowContext(ctx, GetPost, postID)
	err := row.Scan(&post.Id, &post.Creator, &post.Creation, &post.Title,
		&postTextTmp, &post.LikesCount, &post.CommentsCount, pq.Array(&attachs), pq.Array(&types), pq.Array(&subs)) //подписки, при которыз пост доступен
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return models.Post{}, models.WrongData
	}
	if err != nil {
		r.logger.Error(err)
		return models.Post{}, models.InternalError
	}
	post.Text = postTextTmp.String

	row = r.db.QueryRowContext(ctx, GetCreatorPhoto, post.Creator)
	if err = row.Scan(&post.CreatorPhoto); err != nil {
		r.logger.Error(err)
		return models.Post{}, models.InternalError
	}

	for i, v := range attachs {
		if v == uuid.Nil {
			continue
		}
		post.Attachments = append(post.Attachments, models.Attachment{
			Id:   v,
			Type: types[i].String,
		})
	}
	post.Subscriptions, err = r.GetSubsByID(ctx, subs...)
	fmt.Println(post.Subscriptions)
	row = r.db.QueryRowContext(ctx, IsLiked, postID, userID)
	var postUUID, userUUID uuid.UUID
	if err := row.Scan(&postUUID, &userUUID); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.Post{}, models.InternalError
	} else if err == nil { // уже есть запись об этом лайке
		post.IsLiked = true
	}

	return post, err
}

func (r *PostRepo) IsCreator(ctx context.Context, userID, creatorID uuid.UUID) (bool, error) {
	var userIdtmp uuid.UUID
	row := r.db.QueryRowContext(ctx, IsCreator, creatorID)
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	_, err = tx.ExecContext(ctx, DeleteComments, postID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	_, err = tx.ExecContext(ctx, DeleteLikes, postID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	_, err = tx.ExecContext(ctx, DeletePostSubscription, postID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	_, err = tx.ExecContext(ctx, DeletePost, postID)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	if err = tx.Commit(); err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	return nil
}

func (r *PostRepo) IsPostOwner(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
	var userIdtmp uuid.UUID
	row := r.db.QueryRowContext(ctx, GetUserId, postId)
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
	row := r.db.QueryRowContext(ctx, IsLiked, postID, userID)
	if err := row.Scan(&postUUID, &userUUID); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.Like{}, models.InternalError
	} else if err == nil { // уже есть запись об этом лайке
		return models.Like{}, models.WrongData
	}
	// проверяем, есть ли доступ к этому посту
	if ok, err := r.IsPostOwner(ctx, userID, postID); err != nil {
		r.logger.Error(err)
		return models.Like{}, err
	} else if !ok {
		if err := r.IsPostAvailable(ctx, userID, postID); err != nil {
			r.logger.Error(err)
			return models.Like{}, err
		}
	}
	// обновляем кол-во лайков, заодно смотрим, есть ли вообще такой пост
	var like models.Like
	like.PostID = postID
	row = r.db.QueryRowContext(ctx, UpdateLikeCount, 1, postID)

	if err := row.Scan(&like.LikesCount); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.Like{}, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return models.Like{}, models.WrongData
	}

	row = r.db.QueryRowContext(ctx, AddLike, postID, userID)

	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
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
	row := r.db.QueryRowContext(ctx, IsLiked, postID, userID)
	if err := row.Scan(&postUUID, &userUUID); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return models.Like{}, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) { // нет такого лайка
		return models.Like{}, models.WrongData
	}

	var like models.Like
	like.PostID = postID
	row = r.db.QueryRowContext(ctx, UpdateLikeCount, -1, postID)

	if err := row.Scan(&like.LikesCount); err != nil {
		r.logger.Error(err)
		return models.Like{}, models.InternalError
	}

	row = r.db.QueryRowContext(ctx, RemoveLike, postID, userID)

	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		r.logger.Error(err)
		return models.Like{}, models.InternalError
	}

	return like, nil
}

func (r *PostRepo) EditPost(ctx context.Context, postData models.PostEditData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error(err)
		return models.InternalError
	}

	_, err = tx.ExecContext(ctx, UpdatePostInfo, postData.Title, postData.Text, postData.Id)
	if err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	if _, err = tx.ExecContext(ctx, DeletePostSubscriptions, postData.Id); err != nil {
		_ = tx.Rollback()
		r.logger.Error(err)
		return models.InternalError
	}

	for _, sub := range postData.AvailableSubscriptions {
		_, err = tx.ExecContext(ctx, AddSubscriptionsToPost, postData.Id, sub)
		if err != nil {
			_ = tx.Rollback()
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

func (r *PostRepo) IsCommentOwner(ctx context.Context, commentID, userID uuid.UUID) (bool, error) {
	row := r.db.QueryRowContext(ctx, GetUserIdComments, commentID)

	userIdTmp := uuid.UUID{}

	if err := row.Scan(&userIdTmp); err != nil && !errors.Is(sql.ErrNoRows, err) {
		r.logger.Error(err)
		return false, models.InternalError
	} else if errors.Is(sql.ErrNoRows, err) {
		return false, models.WrongData
	}
	if userIdTmp != userID {
		return false, nil
	}
	return true, nil
}
