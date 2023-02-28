package repo

import "database/sql"

type CreatorRepo struct {
	db *sql.DB
}

func NewCreatorRepo(db *sql.DB) *CreatorRepo {
	return &CreatorRepo{db: db}
}

const (
	AUTHOR_INFO = "SELECT name, cover_photo, followers_count, description, posts_count FROM public.creator WHERE user_id=$1;"
	//AUTHOR_POSTS = "SELECT name, cover_photo, followers_count, description, posts_count FROM public.creator WHERE user_id=$1;"

)
