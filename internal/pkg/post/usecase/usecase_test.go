package usecase

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNewPostUsecase(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockPostRepo := mock.NewMockPostRepo(ctl)
	logger, err := zap.NewProduction()
	if err != nil {
		t.Error(err.Error())
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()
	testusecase := NewPostUsecase(mockPostRepo, zapSugar)
	if testusecase.repo != mockPostRepo {
		t.Error("bad constructor")
	}
}

func TestPostUsecase_AddLike(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockPostRepo := mock.NewMockPostRepo(ctl)

	tests := []struct {
		name               string
		repo               *mock.MockPostRepo
		userID             uuid.UUID
		postID             uuid.UUID
		expectedStatusCode error
	}{
		{
			name:               "OK",
			repo:               mockPostRepo,
			userID:             uuid.New(),
			postID:             uuid.New(),
			expectedStatusCode: nil,
		},
		{
			name:               "WrongData",
			repo:               mockPostRepo,
			userID:             uuid.New(),
			postID:             uuid.New(),
			expectedStatusCode: models.WrongData,
		},
		{
			name:               "InternalError",
			repo:               mockPostRepo,
			userID:             uuid.New(),
			postID:             uuid.New(),
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			tests[i].repo.EXPECT().AddLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, nil)
		}
		if tests[i].expectedStatusCode == models.InternalError {
			tests[i].repo.EXPECT().AddLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, models.WrongData)
		}
		if tests[i].expectedStatusCode == models.InternalError {
			tests[i].repo.EXPECT().AddLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, models.InternalError)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &PostUsecase{
				repo: mockPostRepo,
			}

			_, code := h.AddLike(context.Background(), test.userID, test.userID)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestPostUsecase_RemoveLike(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockPostRepo := mock.NewMockPostRepo(ctl)

	tests := []struct {
		name               string
		repo               *mock.MockPostRepo
		userID             uuid.UUID
		postID             uuid.UUID
		expectedStatusCode error
	}{
		{
			name:               "OK",
			repo:               mockPostRepo,
			userID:             uuid.New(),
			postID:             uuid.New(),
			expectedStatusCode: nil,
		},
		{
			name:               "WrongData",
			repo:               mockPostRepo,
			userID:             uuid.New(),
			postID:             uuid.New(),
			expectedStatusCode: models.WrongData,
		},
		{
			name:               "InternalError",
			repo:               mockPostRepo,
			userID:             uuid.New(),
			postID:             uuid.New(),
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			tests[i].repo.EXPECT().RemoveLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, nil)
		}
		if tests[i].expectedStatusCode == models.InternalError {
			tests[i].repo.EXPECT().RemoveLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, models.WrongData)
		}
		if tests[i].expectedStatusCode == models.InternalError {
			tests[i].repo.EXPECT().RemoveLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, models.InternalError)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &PostUsecase{
				repo: mockPostRepo,
			}

			_, code := h.RemoveLike(context.Background(), test.userID, test.userID)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}
