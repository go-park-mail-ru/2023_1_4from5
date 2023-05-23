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

func TestUserUsecase_BecomeCreator(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)

	tests := []struct {
		name               string
		mock               func()
		expectedStatusCode error
	}{
		{
			name: "OK",
			mock: func() {
				mockUserRepo.EXPECT().BecomeCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), nil)
			},
			expectedStatusCode: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}
			test.mock()
			_, err := h.BecomeCreator(context.Background(), models.BecameCreatorInfo{}, uuid.New())
			require.Equal(t, test.expectedStatusCode, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, err))
		})
	}
}

func TestUserUsecase_Donate(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)

	tests := []struct {
		name               string
		mock               func()
		expectedStatusCode error
	}{
		{
			name: "OK",
			mock: func() {
				mockUserRepo.EXPECT().Donate(gomock.Any(), gomock.Any()).Return(float32(10.0), nil)
			},
			expectedStatusCode: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}
			test.mock()
			_, err := h.Donate(context.Background(), models.Donate{})
			require.Equal(t, test.expectedStatusCode, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, err))
		})
	}
}

func TestUserUsecase_UserSubscriptions(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)

	tests := []struct {
		name               string
		mock               func()
		expectedStatusCode error
	}{
		{
			name: "OK",
			mock: func() {
				mockUserRepo.EXPECT().UserSubscriptions(gomock.Any(), gomock.Any()).Return([]models.Subscription{}, nil)
			},
			expectedStatusCode: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}
			test.mock()
			_, err := h.UserSubscriptions(context.Background(), uuid.New())
			require.Equal(t, test.expectedStatusCode, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, err))
		})
	}
}

func TestUserUsecase_DeletePhoto(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)

	tests := []struct {
		name               string
		mock               func()
		expectedStatusCode error
	}{
		{
			name: "OK",
			mock: func() {
				mockUserRepo.EXPECT().DeletePhoto(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatusCode: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}
			test.mock()
			err := h.DeletePhoto(context.Background(), uuid.New())
			require.Equal(t, test.expectedStatusCode, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, err))
		})
	}
}

func TestUserUsecase_UserFollows(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)

	tests := []struct {
		name               string
		mock               func()
		expectedStatusCode error
	}{
		{
			name: "OK",
			mock: func() {
				mockUserRepo.EXPECT().UserFollows(gomock.Any(), gomock.Any()).Return([]models.Follow{}, nil)
			},
			expectedStatusCode: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}
			test.mock()
			_, err := h.UserFollows(context.Background(), uuid.New())
			require.Equal(t, test.expectedStatusCode, err, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, err))
		})
	}
}
