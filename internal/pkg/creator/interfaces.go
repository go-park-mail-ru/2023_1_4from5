package creator

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/creator_mock.go -package=mock

type CreatorUsecase interface {
	GetPage(details models.AccessDetails, creatorUUID string) (models.CreatorPage, error)
}

type CreatorRepo interface {
	GetPage(userId uuid.UUID, creatorId uuid.UUID) (models.CreatorPage, error)
}
