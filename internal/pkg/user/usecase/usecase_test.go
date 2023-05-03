package usecase

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"os"
	"testing"
)

type test struct {
	name               string
	accessDetails      models.AccessDetails
	mockUserRepo       *mock.MockUserRepo
	expectedStatusCode error
}

var testUser models.AccessDetails = models.AccessDetails{Login: "Bashmak1!", Id: uuid.New()}

func userUsecaseTestsSetup(ctl *gomock.Controller) []test {
	os.Setenv("TOKEN_SECRET", "TESTS")
	mockUserRepo := mock.NewMockUserRepo(ctl)

	tests := []test{
		{
			name:               "OK",
			accessDetails:      testUser,
			mockUserRepo:       mockUserRepo,
			expectedStatusCode: nil,
		},
		{
			name:               "InternalError",
			accessDetails:      testUser,
			mockUserRepo:       mockUserRepo,
			expectedStatusCode: models.InternalError,
		},
		{
			name:               "NotFound",
			accessDetails:      testUser,
			mockUserRepo:       mockUserRepo,
			expectedStatusCode: models.NotFound,
		},
	}
	return tests
}

func TestNewUserUsecase(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)
	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()
	testusecase := NewUserUsecase(mockUserRepo, zapSugar)
	if testusecase.repo != mockUserRepo {
		t.Error("bad constructor")
	}
}

func TestUserUsecase_GetProfile(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tests := userUsecaseTestsSetup(ctl)
	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			tests[i].mockUserRepo.EXPECT().GetUserProfile(gomock.Any(), gomock.Any()).Return(models.UserProfile{}, nil)
		} else if tests[i].expectedStatusCode == models.InternalError {
			tests[i].mockUserRepo.EXPECT().GetUserProfile(gomock.Any(), gomock.Any()).Return(models.UserProfile{}, models.InternalError)
		} else {
			tests[i].mockUserRepo.EXPECT().GetUserProfile(gomock.Any(), gomock.Any()).Return(models.UserProfile{}, models.NotFound)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u := &UserUsecase{
				repo: test.mockUserRepo,
			}

			_, code := u.GetProfile(context.Background(), test.accessDetails.Id)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestUserUsecase_UpdatePhoto(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)

	testsUpdatePhoto := []struct {
		name               string
		accessDetails      models.AccessDetails
		mock               func()
		expectedStatusCode error
	}{
		{
			name:          "OK",
			accessDetails: testUser,
			mock: func() {
				mockUserRepo.EXPECT().UpdateProfilePhoto(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatusCode: nil,
		},
		{
			name:          "Internal Error",
			accessDetails: testUser,
			mock: func() {
				mockUserRepo.EXPECT().UpdateProfilePhoto(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.InternalError)
			},
			expectedStatusCode: models.InternalError,
		},
	}

	for _, test := range testsUpdatePhoto {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}
			test.mock()
			_, code := h.UpdatePhoto(context.Background(), test.accessDetails.Id)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestUserUsecase_UpdatePassword(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)

	testsUpdatePhoto := []struct {
		name               string
		id                 uuid.UUID
		password           string
		mock               func()
		expectedStatusCode error
	}{
		{
			name:     "OK",
			id:       uuid.New(),
			password: "1234567aa",
			mock: func() {
				mockUserRepo.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatusCode: nil,
		},
		{
			name:     "Internal Error",
			id:       uuid.New(),
			password: "1234567aa",
			mock: func() {
				mockUserRepo.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.InternalError)
			},
			expectedStatusCode: models.InternalError,
		},
	}

	for _, test := range testsUpdatePhoto {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}
			test.mock()
			err := h.UpdatePassword(context.Background(), test.id, test.password)
			require.Equal(t, test.expectedStatusCode, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, err))
		})
	}
}

func TestUserUsecase_UpdateProfileInfo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)

	testsUpdateProfileInfo := []struct {
		name               string
		id                 uuid.UUID
		profileInfo        models.UpdateProfileInfo
		mock               func()
		expectedStatusCode error
	}{
		{
			name:        "OK",
			id:          uuid.New(),
			profileInfo: models.UpdateProfileInfo{},
			mock: func() {
				mockUserRepo.EXPECT().UpdateProfileInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatusCode: nil,
		},
		{
			name:        "Internal Error",
			id:          uuid.New(),
			profileInfo: models.UpdateProfileInfo{},
			mock: func() {
				mockUserRepo.EXPECT().UpdateProfileInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.InternalError)
			},
			expectedStatusCode: models.InternalError,
		},
	}

	for _, test := range testsUpdateProfileInfo {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}
			test.mock()
			err := h.UpdateProfileInfo(context.Background(), test.profileInfo, test.id)
			require.Equal(t, test.expectedStatusCode, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, err))
		})
	}
}

func TestUserUsecase_Donate(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)

	testsDonateInfo := []struct {
		name               string
		id                 uuid.UUID
		donateInfo         models.Donate
		mock               func()
		expectedStatusCode error
		expectedRes        int64
	}{
		{
			name:       "OK",
			id:         uuid.New(),
			donateInfo: models.Donate{},
			mock: func() {
				mockUserRepo.EXPECT().Donate(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(100), nil)
			},
			expectedStatusCode: nil,
			expectedRes:        100,
		},
		{
			name:       "Internal Error",
			id:         uuid.New(),
			donateInfo: models.Donate{},
			mock: func() {
				mockUserRepo.EXPECT().Donate(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), models.InternalError)
			},
			expectedStatusCode: models.InternalError,
			expectedRes:        0,
		},
	}

	for _, test := range testsDonateInfo {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}
			test.mock()
			got, err := h.Donate(context.Background(), test.donateInfo, test.id)
			require.Equal(t, test.expectedStatusCode, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, err))
			require.Equal(t, test.expectedRes, got, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedRes, err))
		})
	}
}
