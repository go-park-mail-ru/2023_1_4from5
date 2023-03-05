package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewUserUsecase(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctl)
	testusecase := NewUserUsecase(mockUserRepo)
	if testusecase.repo != mockUserRepo {
		t.Error("bad constructor")
	}
}

var testUser models.AccessDetails = models.AccessDetails{Login: "Bashmak1!", Id: uuid.New()}

func TestUserUsecase_GetProfile(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("SECRET", "TESTS")
	mockUserRepo := mock.NewMockUserRepo(ctl)

	tests := []struct {
		name               string
		accessDetails      models.AccessDetails
		fields             *mock.MockUserRepo
		expectedStatusCode error
	}{
		{
			name:               "OK",
			accessDetails:      testUser,
			fields:             mockUserRepo,
			expectedStatusCode: nil,
		},
		{
			name:               "Unauthorized",
			accessDetails:      testUser,
			fields:             mockUserRepo,
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			mockUserRepo.EXPECT().GetUserProfile(gomock.Any()).Return(models.UserProfile{}, nil)
		} else {
			mockUserRepo.EXPECT().GetUserProfile(gomock.Any()).Return(models.UserProfile{}, models.InternalError)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}

			_, code := h.GetProfile(test.accessDetails)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestUserUsecase_GetHomePage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("SECRET", "TESTS")
	mockUserRepo := mock.NewMockUserRepo(ctl)

	tests := []struct {
		name               string
		accessDetails      models.AccessDetails
		fields             *mock.MockUserRepo
		expectedStatusCode error
	}{
		{
			name:               "OK",
			accessDetails:      testUser,
			fields:             mockUserRepo,
			expectedStatusCode: nil,
		},
		{
			name:               "Unauthorized",
			accessDetails:      testUser,
			fields:             mockUserRepo,
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			mockUserRepo.EXPECT().GetHomePage(gomock.Any()).Return(models.UserHomePage{}, nil)
		} else {
			mockUserRepo.EXPECT().GetHomePage(gomock.Any()).Return(models.UserHomePage{}, models.InternalError)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserUsecase{
				repo: mockUserRepo,
			}

			_, code := h.GetHomePage(test.accessDetails)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}
