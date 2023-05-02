package usecase

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

	logger := zap.NewNop()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	testusecase := NewAuthUsecase(mockAuthRepo, mockTokenGen, mockEncrypter, zapSugar)
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

func TestNewEncryptor(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	os.Setenv("ENCRYPTER_SECRET", "TESTS")

	enc, err := NewEncryptor()
	require.Equal(t, "TESTS", enc.salt, fmt.Errorf("expected %s, got %s",
		"TESTS", enc.salt))
	require.Equal(t, nil, err, fmt.Errorf("error wasnt expected, got %s",
		err))
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
			mockAuthRepo.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, nil)
		} else {
			mockAuthRepo.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, models.WrongPassword)
		}
		mockEncrypter.EXPECT().EncryptPswd(gomock.Any(), gomock.Any()).Return("test")
		mockTokenGen.EXPECT().GetJWTToken(gomock.Any(), gomock.Any()).Return("TEST TOKEN", nil)
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &AuthUsecase{
				repo:      mockAuthRepo,
				tokenator: mockTokenGen,
				encrypter: mockEncrypter,
			}

			_, code := h.SignIn(context.Background(), models.LoginUser{Login: testUsers[i].Login, PasswordHash: testUsers[i].PasswordHash})
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
		mockEncrypter.EXPECT().EncryptPswd(gomock.Any(), gomock.Any()).Return("test")
		switch tests[i].expectedStatusCode {
		case models.WrongData:
			mockAuthRepo.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, models.WrongPassword)
			continue
		case nil:
			mockAuthRepo.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, models.NotFound)
			mockAuthRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(models.User{}, nil)
		default:
			mockAuthRepo.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, models.NotFound)
			mockAuthRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(models.User{}, models.InternalError)
		}
		mockTokenGen.EXPECT().GetJWTToken(gomock.Any(), gomock.Any()).Return("TEST TOKEN", nil)
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &AuthUsecase{
				repo:      mockAuthRepo,
				tokenator: mockTokenGen,
				encrypter: mockEncrypter,
			}

			_, code := h.SignUp(context.Background(), testUsers[i])
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestAuthUsecase_Logout(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockAuthRepo := mock.NewMockAuthRepo(ctl)
	mockTokenGen := mock.NewMockTokenGenerator(ctl)
	mockEncrypter := mock.NewMockEncrypter(ctl)

	tests := []struct {
		name               string
		fields             testUsecase
		expectedStatusCode error
	}{
		{
			name:               "Nil err",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: nil,
		},
		{
			name:               "InternalErr",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			mockAuthRepo.EXPECT().IncUserVersion(gomock.Any(), gomock.Any()).Return(int64(1), nil)
		} else {
			mockAuthRepo.EXPECT().IncUserVersion(gomock.Any(), gomock.Any()).Return(int64(0), models.InternalError)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u := &AuthUsecase{
				repo:      mockAuthRepo,
				tokenator: mockTokenGen,
				encrypter: mockEncrypter,
			}

			_, code := u.IncUserVersion(context.Background(), models.AccessDetails{})
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestAuthUsecase_CheckUserVersion(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockAuthRepo := mock.NewMockAuthRepo(ctl)
	mockTokenGen := mock.NewMockTokenGenerator(ctl)
	mockEncrypter := mock.NewMockEncrypter(ctl)

	tests := []struct {
		name               string
		fields             testUsecase
		expectedStatusCode error
	}{
		{
			name:               "Nil err",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: nil,
		},
		{
			name:               "InternalErr",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].expectedStatusCode == nil {
			mockAuthRepo.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(int64(1), nil)
		} else {
			mockAuthRepo.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(int64(0), models.InternalError)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &AuthUsecase{
				repo:      mockAuthRepo,
				tokenator: mockTokenGen,
				encrypter: mockEncrypter,
			}

			_, code := h.CheckUserVersion(context.Background(), models.AccessDetails{})
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestAuthUsecase_CheckUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockAuthRepo := mock.NewMockAuthRepo(ctl)
	mockTokenGen := mock.NewMockTokenGenerator(ctl)
	mockEncrypter := mock.NewMockEncrypter(ctl)

	tests := []struct {
		name               string
		fields             testUsecase
		expectedStatusCode error
	}{
		{
			name:               "Nil err",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: nil,
		},
		{
			name:               "InternalErr",
			fields:             testUsecase{mockAuthRepo, mockTokenGen, mockEncrypter},
			expectedStatusCode: models.InternalError,
		},
	}

	for i := 0; i < len(tests); i++ {
		mockEncrypter.EXPECT().EncryptPswd(gomock.Any(), gomock.Any()).Return("testPasswordHash")
		if tests[i].expectedStatusCode == nil {
			mockAuthRepo.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, nil)
		} else {
			mockAuthRepo.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, models.InternalError)
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &AuthUsecase{
				repo:      mockAuthRepo,
				tokenator: mockTokenGen,
				encrypter: mockEncrypter,
			}

			_, code := h.CheckUser(context.Background(), models.User{})
			require.Equal(t, test.expectedStatusCode, code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedStatusCode, code))
		})
	}
}

func TestAuthUsecase_EncryptPwd(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockAuthRepo := mock.NewMockAuthRepo(ctl)
	mockTokenGen := mock.NewMockTokenGenerator(ctl)
	mockEncrypter := mock.NewMockEncrypter(ctl)

	mockEncrypter.EXPECT().EncryptPswd(gomock.Any(), gomock.Any()).Return("testPasswordHash")

	h := &AuthUsecase{
		repo:      mockAuthRepo,
		tokenator: mockTokenGen,
		encrypter: mockEncrypter,
	}

	res := h.EncryptPwd(context.Background(), "test")
	require.Equal(t, "testPasswordHash", res, fmt.Errorf("expected %s, got %s",
		"testPasswordHash", res))
}
