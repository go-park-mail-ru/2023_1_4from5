package repo

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
)

const (
	InsertPost = `INSERT INTO "post"(post_id, creator_id, title, post_text) VALUES($1, $2, $3, $4);`
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
