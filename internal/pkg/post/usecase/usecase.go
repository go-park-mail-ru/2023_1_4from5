package usecase

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/google/uuid"
)

type PostUsecase struct {
	repo post.PostRepo
}

func NewPostUsecase(repo post.PostRepo) *PostUsecase {
	return &PostUsecase{repo: repo}
}

func (u *PostUsecase) CreatePost(postData models.PostCreationData) (uuid.UUID, error) {
	postData.Id = uuid.New()
	if err := u.repo.CreatePost(postData); err != nil {
		return uuid.Nil, models.InternalError
	}
	return postData.Id, nil
}

func (u *PostUsecase) AddLike(userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
	like, err := u.repo.AddLike(userID, postID)
	if err != nil {
		return models.Like{}, err
	}
	return like, nil
}

func (u *PostUsecase) RemoveLike(userID uuid.UUID, postID uuid.UUID) (models.Like, error) {
	like, err := u.repo.RemoveLike(userID, postID)
	if err != nil {
		return models.Like{}, err
	}
	return like, nil
}
