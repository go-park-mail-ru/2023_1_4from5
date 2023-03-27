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

func (u *PostUsecase) CreatePost(postData models.PostCreationData) error {
	if err := u.repo.CreatePost(postData); err != nil {
		return models.InternalError
	}
	return nil
}
func (u *PostUsecase) DeletePost(postID uuid.UUID) error {
	if err := u.repo.DeletePost(postID); err != nil {
		return models.InternalError
	}
	return nil
}
func (u *PostUsecase) IsPostOwner(userId uuid.UUID, postId uuid.UUID) (bool, error) {
	return u.repo.IsPostOwner(userId, postId)
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
