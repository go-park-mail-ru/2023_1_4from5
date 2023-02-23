package models

// easyjson -all ./internal/models/creator.go

type Creator struct {
	Id             int    `json:"id"`
	UserId         int    `json:"user_id"`
	CoverPhoto     string `json:"cover_photo"`
	FollowersCount int    `json:"followers_count"`
	Description    string `json:"description"`
	PostsCount     int    `json:"posts_count"`
}
