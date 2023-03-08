package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewCreatorUsecase(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockCreatorRepo := mock.NewMockCreatorRepo(ctl)
	testusecase := NewCreatorUsecase(mockCreatorRepo)
	if testusecase.repo != mockCreatorRepo {
		t.Error("bad constructor")
	}
}

var testUser models.AccessDetails = models.AccessDetails{Login: "Bashmak1!", Id: uuid.New()}

func TestCreatorUsecase_GetPage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("TOKEN_SECRET", "TESTS")
	mockCreatorRepo := mock.NewMockCreatorRepo(ctl)

	tests := []struct {
		name               string
		accessDetails      models.AccessDetails
		creatorID          string
		fields             *mock.MockCreatorRepo
		expectedStatusCode error
	}{
		{
			name:               "OK",
			accessDetails:      testUser,
			creatorID:          uuid.New().String(),
			fields:             mockCreatorRepo,
			expectedStatusCode: nil,
		},
		{
			name:               "WrongData: wrong creatorUUId",
			accessDetails:      testUser,
			creatorID:          "123",
			fields:             mockCreatorRepo,
			expectedStatusCode: models.WrongData,
		},
		{
			name:               "InternalError",
			accessDetails:      testUser,
			creatorID:          uuid.New().String(),
			fields:             mockCreatorRepo,
			expectedStatusCode: models.InternalError,
		},
		{
			name:               "WrongData: no such creator",
			accessDetails:      testUser,
			creatorID:          uuid.New().String(),
			fields:             mockCreatorRepo,
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			mockCreatorRepo.EXPECT().GetPage(gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, nil)
			continue
		}
		if tests[i].name == "WrongData: wrong creatorUUId" {
			continue
		}
		if tests[i].expectedStatusCode == models.InternalError {
			mockCreatorRepo.EXPECT().GetPage(gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, models.InternalError)
		} else {
			mockCreatorRepo.EXPECT().GetPage(gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, models.WrongData)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorUsecase{
				repo: mockCreatorRepo,
			}

			_, code := h.GetPage(&test.accessDetails, test.creatorID)
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %e, got %e,",
				test.name, test.expectedStatusCode, code))
		})
	}
}
