package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testUsecase struct {
	AuthRepo  auth.AuthRepo
	TokenGen  auth.TokenGenerator
	encrypter auth.Encrypter
}

var testUsers = []models.User{
	{
		Login:        "Alligator19",
		PasswordHash: "Password123!",
		Name:         "Alligator",
	},
	{
		Login:        "Donald123",
		PasswordHash: "Donald123!",
		Name:         "Donald Brown",
	},
	{
		Login:        "",
		PasswordHash: "",
		Name:         "Дарья Такташова",
	},
}

func TestNewAuthUsecase(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockAuthRepo := mock.NewMockAuthRepo(ctl)
	mockTokenGen := mock.NewMockTokenGenerator(ctl)
	mockEncrypter := mock.NewMockEncrypter(ctl)
	testusecase := NewAuthUsecase(mockAuthRepo, mockTokenGen, mockEncrypter)
	if testusecase.repo != mockAuthRepo {
		t.Error("bad constructor")
	}

	if testusecase.tokenator != mockTokenGen {
		t.Error("bad constructor")
	}

	if testusecase.encrypter != mockEncrypter {
		t.Error("bad constructor")
	}
}

func TestAuthUsecase_SignIn(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	os.Setenv("SECRET", "TESTS")
	mockAuthRepo := mock.NewMockAuthRepo(ctl)
	mockTokenGen := mock.NewMockTokenGenerator(ctl)
	mockEncrypter := mock.NewMockEncrypter(ctl)

	tests := []struct {
		name               string
		fields             testUsecase
		expectedStatusCode error
	}{
		{
			name:               "OK",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: nil,
		},
		{
			name:               "Unauthorized",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: models.NotFound,
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			mockAuthRepo.EXPECT().CheckUser(gomock.Any()).Return(models.User{}, nil)
		} else {
			mockAuthRepo.EXPECT().CheckUser(gomock.Any()).Return(models.User{}, models.WrongPassword)
		}
		mockEncrypter.EXPECT().EncryptPswd(gomock.Any()).Return("test")
		mockTokenGen.EXPECT().GetToken(gomock.Any()).Return("TEST TOKEN", nil)
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &AuthUsecase{
				repo:      mockAuthRepo,
				tokenator: mockTokenGen,
				encrypter: mockEncrypter,
			}

			_, code := h.SignIn(models.LoginUser{Login: testUsers[i].Login, PasswordHash: testUsers[i].PasswordHash})
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestAuthUsecase_SignUp(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	os.Setenv("SECRET", "TESTS")
	mockAuthRepo := mock.NewMockAuthRepo(ctl)
	mockTokenGen := mock.NewMockTokenGenerator(ctl)
	mockEncrypter := mock.NewMockEncrypter(ctl)

	tests := []struct {
		name               string
		fields             testUsecase
		expectedStatusCode error
	}{
		{
			name:               "Conflict, user with that login already exists",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: models.WrongData,
		},
		{
			name:               "OK",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: nil,
		},
		{
			name:               "ServerError",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		mockEncrypter.EXPECT().EncryptPswd(gomock.Any()).Return("test")
		switch tests[i].expectedStatusCode {
		case models.WrongData:
			mockAuthRepo.EXPECT().CheckUser(gomock.Any()).Return(models.User{}, models.WrongPassword)
			continue
		case nil:
			mockAuthRepo.EXPECT().CheckUser(gomock.Any()).Return(models.User{}, models.NotFound)
			mockAuthRepo.EXPECT().CreateUser(gomock.Any()).Return(models.User{}, nil)
		default:
			mockAuthRepo.EXPECT().CheckUser(gomock.Any()).Return(models.User{}, models.NotFound)
			mockAuthRepo.EXPECT().CreateUser(gomock.Any()).Return(models.User{}, models.InternalError)
		}
		mockTokenGen.EXPECT().GetToken(gomock.Any()).Return("TEST TOKEN", nil)
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &AuthUsecase{
				repo:      mockAuthRepo,
				tokenator: mockTokenGen,
				encrypter: mockEncrypter,
			}

			_, code := h.SignUp(testUsers[i])
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedStatusCode, code))
		})
	}
}
