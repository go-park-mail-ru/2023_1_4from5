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

			_, code := u.GetProfile(context.Background(), test.accessDetails)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestUserUsecase_GetHomePage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tests := userUsecaseTestsSetup(ctl)

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			tests[i].mockUserRepo.EXPECT().GetHomePage(gomock.Any(), gomock.Any()).Return(models.UserHomePage{}, nil)
		} else if tests[i].expectedStatusCode == models.InternalError {
			tests[i].mockUserRepo.EXPECT().GetHomePage(gomock.Any(), gomock.Any()).Return(models.UserHomePage{}, models.InternalError)
		} else {
			tests[i].mockUserRepo.EXPECT().GetHomePage(gomock.Any(), gomock.Any()).Return(models.UserHomePage{}, models.NotFound)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: test.mockUserRepo,
			}

			_, code := h.GetHomePage(context.Background(), test.accessDetails)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}
