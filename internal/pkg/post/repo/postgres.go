package repo

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

const (
	InsertPost      = `INSERT INTO "post"(post_id, creator_id, title, post_text) VALUES($1, $2, $3, $4);`
	AddLike         = `INSERT INTO "like_post"(post_id, user_id) VALUES($1, $2);`
	RemoveLike      = `DELETE FROM "like_post" WHERE post_id = $1 AND user_id = $2`
	UpdateLikeCount = `UPDATE "post" SET likes_count = likes_count + $1 WHERE post_id = $2 RETURNING likes_count;`
	IsLiked         = `SELECT post_id, user_id FROM "like_post" WHERE post_id = $1 AND user_id = $2`
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) CreatePost(postData models.PostCreationData) error {
	row := r.db.QueryRow(InsertPost, postData.Id, postData.Creator, postData.Title, postData.Text)

	if err := row.Err(); err != nil {
		return models.InternalError
	}

	return nil
}

func (r *PostRepo) AddLike(userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
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
		return models.Like{}, models.InternalError
	}

	return like, nil
}

func (r *PostRepo) RemoveLike(userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
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
		return models.Like{}, models.InternalError
	}

	row = r.db.QueryRow(RemoveLike, postID, userID)

	if err := row.Err(); err != nil {
		return models.Like{}, models.InternalError
	}

	return like, nil
}