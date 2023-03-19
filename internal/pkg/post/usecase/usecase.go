package usecase

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
)

type PostUsecase struct {
	repo post.PostRepo
}

func NewPostUsecase(repo post.PostRepo) *PostUsecase {
	return &PostUsecase{repo: repo}
}

func (u *PostUsecase) CreatePost(postData models.PostCreationData) error {
	if err := u.repo.CreatePost(postData); err != nil {
		return models.InternalError
	}
	return nil
}
