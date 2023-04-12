package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNewCreatorUsecase(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockCreatorRepo := mock.NewMockCreatorRepo(ctl)
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
	testusecase := NewCreatorUsecase(mockCreatorRepo, zapSugar)
	if testusecase.repo != mockCreatorRepo {
		t.Error("bad constructor")
	}
}

var testUser = &models.AccessDetails{Login: "Bashmak1!", Id: uuid.New()}

func TestCreatorUsecase_GetPage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	mockCreatorRepo := mock.NewMockCreatorRepo(ctl)

	tests := []struct {
		name               string
		accessDetails      *models.AccessDetails
		creatorID          string
		repo               *mock.MockCreatorRepo
		expectedStatusCode error
	}{
		{
			name:               "OK",
			accessDetails:      testUser,
			creatorID:          uuid.New().String(),
			repo:               mockCreatorRepo,
			expectedStatusCode: nil,
		},
		{
			name:               "WrongData: wrong creatorUUId",
			accessDetails:      testUser,
			creatorID:          "123",
			repo:               mockCreatorRepo,
			expectedStatusCode: models.WrongData,
		},
		{
			name:               "InternalError",
			accessDetails:      testUser,
			creatorID:          uuid.New().String(),
			repo:               mockCreatorRepo,
			expectedStatusCode: models.InternalError,
		},
		{
			name:               "WrongData: no such creator",
			accessDetails:      testUser,
			creatorID:          uuid.New().String(),
			repo:               mockCreatorRepo,
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			tests[i].repo.EXPECT().GetPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, nil)
			continue
		}
		if tests[i].name == "WrongData: wrong creatorUUId" {
			continue
		}
		if tests[i].expectedStatusCode == models.InternalError {
			tests[i].repo.EXPECT().GetPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, models.InternalError)
		} else {
			tests[i].repo.EXPECT().GetPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, models.WrongData)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorUsecase{
				repo:   mockCreatorRepo,
				logger: zapSugar,
			}

			_, code := h.GetPage(context.Background(), test.accessDetails, test.creatorID)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestCreatorUsecase_CreateAim(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	logger := zap.NewNop()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	mockCreatorRepo := mock.NewMockCreatorRepo(ctl)

	tests := []struct {
		name        string
		mock        func()
		input       models.Aim
		expectedErr error
	}{
		{
			name: "OK",
			mock: func() {
				mockCreatorRepo.EXPECT().CreateAim(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       models.Aim{},
			expectedErr: nil,
		},
		{
			name: "Error",
			mock: func() {
				mockCreatorRepo.EXPECT().CreateAim(gomock.Any(), gomock.Any()).Return(errors.New("test"))
			},
			input:       models.Aim{},
			expectedErr: errors.New("test"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorUsecase{
				repo:   mockCreatorRepo,
				logger: zapSugar,
			}
			test.mock()

			err := h.CreateAim(context.Background(), test.input)
			require.Equal(t, test.expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedErr, err))
		})
	}
}

var creators = make([]models.Creator, 1)

func TestCreatorUsecase_GetAllCreators(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	mockCreatorRepo := mock.NewMockCreatorRepo(ctl)

	tests := []struct {
		name        string
		mock        func()
		input       models.Aim
		expectedErr error
		expectedRes []models.Creator
	}{
		{
			name: "OK",
			mock: func() {
				mockCreatorRepo.EXPECT().GetAllCreators(gomock.Any()).Return(creators, nil)
			},
			input:       models.Aim{},
			expectedErr: nil,
			expectedRes: creators,
		},
		{
			name: "Error",
			mock: func() {
				mockCreatorRepo.EXPECT().GetAllCreators(gomock.Any()).Return(nil, errors.New("test"))
			},
			input:       models.Aim{},
			expectedErr: errors.New("test"),
			expectedRes: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorUsecase{
				repo:   mockCreatorRepo,
				logger: zapSugar,
			}
			test.mock()

			creators, err := h.GetAllCreators(context.Background())
			require.Equal(t, test.expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedErr, err))
			require.Equal(t, test.expectedRes, creators, fmt.Errorf("%s :  expected %v, got %v,",
				test.name, test.expectedRes, creators))
		})
	}
}
