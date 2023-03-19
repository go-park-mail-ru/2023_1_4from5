package post

//go:generate mockgen -source=interfaces.go -destination=./mocks/post_mock.go -package=mock

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
)

type PostUsecase interface {
	CreatePost(postData models.PostCreationData) error
}
type PostRepo interface {
	CreatePost(postData models.PostCreationData) error
}
