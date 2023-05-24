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

	test := []struct {
		name               string
		userID             uuid.UUID
		creatorID          uuid.UUID
		repo               *mock.MockCreatorRepo
		expectedStatusCode error
	}{
		{
			name:               "OK",
			userID:             uuid.New(),
			creatorID:          uuid.New(),
			repo:               mockCreatorRepo,
			expectedStatusCode: nil,
		},
	}

	test[0].repo.EXPECT().GetPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, nil)

	t.Run(test[0].name, func(t *testing.T) {
		h := &CreatorUsecase{
			repo:   mockCreatorRepo,
			logger: zapSugar,
		}

		_, code := h.GetPage(context.Background(), test[0].userID, test[0].creatorID)
		require.Equal(t, test[0].expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
			test[0].name, test[0].expectedStatusCode, code))
	})
}

func TestCreatorUsecase_UpdateCreatorData(t *testing.T) {
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

	test := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "OK",
			expectedErr: nil,
		},
	}

	mockCreatorRepo.EXPECT().UpdateCreatorData(gomock.Any(), gomock.Any()).Return(nil)

	t.Run(test[0].name, func(t *testing.T) {
		h := &CreatorUsecase{
			repo:   mockCreatorRepo,
			logger: zapSugar,
		}

		err := h.UpdateCreatorData(context.Background(), models.UpdateCreatorInfo{
			Description: "",
			CreatorName: "",
			CreatorID:   uuid.UUID{},
		})
		require.Equal(t, test[0].expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
			test[0].name, test[0].expectedErr, err))
	})
}

func TestCreatorUsecase_CheckIfCreator(t *testing.T) {
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

	test := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "OK",
			expectedErr: nil,
		},
	}

	mockCreatorRepo.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(uuid.Nil, nil)

	t.Run(test[0].name, func(t *testing.T) {
		h := &CreatorUsecase{
			repo:   mockCreatorRepo,
			logger: zapSugar,
		}

		_, err := h.CheckIfCreator(context.Background(), uuid.New())
		require.Equal(t, test[0].expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
			test[0].name, test[0].expectedErr, err))
	})

}

func TestCreatorUsecase_FindCreators(t *testing.T) {
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

	test := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "OK",
			expectedErr: nil,
		},
	}

	mockCreatorRepo.EXPECT().FindCreators(gomock.Any(), gomock.Any()).Return(nil, nil)

	t.Run(test[0].name, func(t *testing.T) {
		h := &CreatorUsecase{
			repo:   mockCreatorRepo,
			logger: zapSugar,
		}

		_, err := h.FindCreators(context.Background(), "test")
		require.Equal(t, test[0].expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
			test[0].name, test[0].expectedErr, err))
	})
}

func TestCreatorUsecase_GetFeed(t *testing.T) {
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

	test := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "OK",
			expectedErr: nil,
		},
	}

	mockCreatorRepo.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(nil, nil)

	t.Run(test[0].name, func(t *testing.T) {
		h := &CreatorUsecase{
			repo:   mockCreatorRepo,
			logger: zapSugar,
		}

		_, err := h.GetFeed(context.Background(), uuid.New())
		require.Equal(t, test[0].expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
			test[0].name, test[0].expectedErr, err))
	})
}

func TestCreatorUsecase_DeleteCoverPhoto(t *testing.T) {
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

	test := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "OK",
			expectedErr: nil,
		},
	}

	mockCreatorRepo.EXPECT().DeleteCoverPhoto(gomock.Any(), gomock.Any()).Return(nil)

	t.Run(test[0].name, func(t *testing.T) {
		h := &CreatorUsecase{
			repo:   mockCreatorRepo,
			logger: zapSugar,
		}

		err := h.DeleteCoverPhoto(context.Background(), uuid.New())
		require.Equal(t, test[0].expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
			test[0].name, test[0].expectedErr, err))
	})
}

func TestCreatorUsecase_DeleteProfilePhoto(t *testing.T) {
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

	test := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "OK",
			expectedErr: nil,
		},
	}

	mockCreatorRepo.EXPECT().DeleteProfilePhoto(gomock.Any(), gomock.Any()).Return(nil)

	t.Run(test[0].name, func(t *testing.T) {
		h := &CreatorUsecase{
			repo:   mockCreatorRepo,
			logger: zapSugar,
		}

		err := h.DeleteProfilePhoto(context.Background(), uuid.New())
		require.Equal(t, test[0].expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
			test[0].name, test[0].expectedErr, err))
	})
}

func TestCreatorUsecase_UpdateCoverPhoto(t *testing.T) {
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
		expectedErr error
		mock        func()
	}{
		{
			name:        "OK",
			expectedErr: nil,
			mock: func() {
				mockCreatorRepo.EXPECT().UpdateCoverPhoto(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:        "Internal Err",
			expectedErr: models.InternalError,
			mock: func() {
				mockCreatorRepo.EXPECT().UpdateCoverPhoto(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.InternalError)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorUsecase{
				repo:   mockCreatorRepo,
				logger: zapSugar,
			}
			test.mock()
			_, err := h.UpdateCoverPhoto(context.Background(), uuid.New())
			require.Equal(t, test.expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedErr, err))
		})
	}
}

func TestCreatorUsecase_UpdateProfilePhoto(t *testing.T) {
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
		expectedErr error
		mock        func()
	}{
		{
			name:        "OK",
			expectedErr: nil,
			mock: func() {
				mockCreatorRepo.EXPECT().UpdateProfilePhoto(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:        "Internal Err",
			expectedErr: models.InternalError,
			mock: func() {
				mockCreatorRepo.EXPECT().UpdateProfilePhoto(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.InternalError)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorUsecase{
				repo:   mockCreatorRepo,
				logger: zapSugar,
			}
			test.mock()
			_, err := h.UpdateProfilePhoto(context.Background(), uuid.New())
			require.Equal(t, test.expectedErr, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedErr, err))
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
