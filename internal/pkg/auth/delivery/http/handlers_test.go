package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func bodyPrepare(user models.User) []byte {
	userjson, err := json.Marshal(&user)
	if err != nil {
		return nil
	}
	return userjson
}

type args struct {
	r                *http.Request
	expectedResponse http.Response
}

var testUsers = []models.User{
	{
		Login:        "Dasha2003",
		PasswordHash: "Dasha2003!",
		Name:         "Дарья Такташова",
	},
	{
		Login:        "Donald123",
		PasswordHash: "Donald123!",
		Name:         "Donald Brown",
	},
	{
		Login:        "Alligator19",
		PasswordHash: "Password123!",
		Name:         "Alligator",
	},
	{
		Login:        "Bad",
		PasswordHash: "User",
		Name:         "KKK",
	},
}

func TestNewAuthHandler(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockClient := mock.NewMockAuthServiceClient(ctl)
	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()
	testHandler := NewAuthHandler(mockClient, zapSugar)
	if testHandler.client != mockClient {
		t.Error("bad constructor")
	}
}

func TestAuthHandler_SignIn(t *testing.T) {
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

	mockClient := mock.NewMockAuthServiceClient(ctl)

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "OK",
			args: args{
				r: httptest.NewRequest("POST", "/signIn",
					bytes.NewReader(bodyPrepare(testUsers[0]))),
				expectedResponse: http.Response{StatusCode: http.StatusOK},
			},
			mock: func() {
				mockClient.EXPECT().
					SignIn(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.Token{
						Cookie: "test",
						Error:  "",
					}, nil)
			},
		},
		{
			name: "Unauthorized",
			args: args{
				r: httptest.NewRequest("POST", "/signIn",
					bytes.NewReader(bodyPrepare(testUsers[1]))),
				expectedResponse: http.Response{StatusCode: http.StatusUnauthorized},
			},
			mock: func() {
				mockClient.EXPECT().
					SignIn(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.Token{
						Cookie: "",
						Error:  models.NotFound.Error(),
					}, nil)
			},
		},
		{
			name: "BadRequest",
			args: args{
				r: httptest.NewRequest("POST", "/signIn",
					bytes.NewReader([]byte("Trying to signIn"))),
				expectedResponse: http.Response{StatusCode: http.StatusBadRequest},
			},
			mock: func() {
			},
		},
		{
			name: "InternalErr",
			args: args{
				r: httptest.NewRequest("POST", "/signIn",
					bytes.NewReader(bodyPrepare(testUsers[1]))),
				expectedResponse: http.Response{StatusCode: http.StatusInternalServerError},
			},
			mock: func() {
				mockClient.EXPECT().
					SignIn(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.Token{
						Cookie: "",
						Error:  "",
					}, errors.New("test"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &AuthHandler{
				client: mockClient,
				logger: zapSugar,
			}
			w := httptest.NewRecorder()
			test.mock()

			h.SignIn(w, test.args.r)
			require.Equal(t, test.args.expectedResponse.StatusCode, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.args.expectedResponse.StatusCode, w.Code))
		})
	}
}

func TestAuthHandler_SignUp(t *testing.T) {
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

	mockClient := mock.NewMockAuthServiceClient(ctl)

	tests := []struct {
		name  string
		Login string
		args  args
		mock  func()
	}{
		{
			name:  "OK",
			Login: testUsers[0].Login,
			args: args{
				r: httptest.NewRequest("POST", "/signUp",
					bytes.NewReader(bodyPrepare(testUsers[0]))),
				expectedResponse: http.Response{StatusCode: http.StatusOK},
			},
			mock: func() {
				mockClient.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.Token{
						Cookie: "test",
						Error:  "",
					}, nil)
			},
		},
		{
			name:  "Conflict",
			Login: testUsers[1].Login,
			args: args{
				r: httptest.NewRequest("POST", "/signUp",
					bytes.NewReader(bodyPrepare(testUsers[1]))),
				expectedResponse: http.Response{StatusCode: http.StatusConflict},
			},
			mock: func() {
				mockClient.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.Token{
						Cookie: "",
						Error:  models.WrongData.Error(),
					}, nil)
			},
		},
		{
			name:  "InternalErr from DB",
			Login: testUsers[1].Login,
			args: args{
				r: httptest.NewRequest("POST", "/signUp",
					bytes.NewReader(bodyPrepare(testUsers[1]))),
				expectedResponse: http.Response{StatusCode: http.StatusInternalServerError},
			},
			mock: func() {
				mockClient.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.Token{
						Cookie: "",
						Error:  models.InternalError.Error(),
					}, nil)
			},
		},
		{
			name:  "BadRequest because no json",
			Login: testUsers[2].Login,
			args: args{
				r: httptest.NewRequest("POST", "/signUp",
					bytes.NewReader([]byte("ppppp"))),
				expectedResponse: http.Response{StatusCode: http.StatusBadRequest},
			},
			mock: func() {
			},
		},
		{
			name:  "InternalError",
			Login: testUsers[2].Login,
			args: args{
				r: httptest.NewRequest("POST", "/signUp",
					bytes.NewReader(bodyPrepare(testUsers[2]))),
				expectedResponse: http.Response{StatusCode: http.StatusInternalServerError},
			},
			mock: func() {
				mockClient.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.Token{
						Cookie: "",
						Error:  "",
					}, errors.New("test"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &AuthHandler{
				client: mockClient,
				logger: zapSugar,
			}
			w := httptest.NewRecorder()
			test.mock()

			h.SignUp(w, test.args.r)
			require.Equal(t, test.args.expectedResponse.StatusCode, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.args.expectedResponse.StatusCode, w.Code))
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
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

	mockClient := mock.NewMockAuthServiceClient(ctl)

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUsers[1].Login, Id: uuid.New()})
	var name = "SSID"
	expires := time.Now().UTC().Add(time.Hour)
	value := bdy

	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "BadRequest wrong value for token",
			args: args{
				r:                httptest.NewRequest("POST", "/logout", nil),
				expectedResponse: http.Response{StatusCode: http.StatusBadRequest},
			},
			mock: func() {
				value = "token"
			},
		},
		{
			name: "OK",
			args: args{
				r:                httptest.NewRequest("POST", "/logout", nil),
				expectedResponse: http.Response{StatusCode: http.StatusOK},
			},
			mock: func() {
				name = "SSID"
				expires = time.Now().UTC().Add(time.Hour)
				value = bdy

				mockClient.EXPECT().
					IncUserVersion(gomock.Any(), gomock.Any()).
					Return(&generatedCommon.Empty{
						Error: "",
					}, nil)
			},
		},
		{
			name: "BadRequest wrong Cookie name",
			args: args{
				r:                httptest.NewRequest("POST", "/logout", nil),
				expectedResponse: http.Response{StatusCode: http.StatusBadRequest},
			},
			mock: func() {
				name = "test"
			},
		},
		{
			name: "InternalError from service",
			args: args{
				r:                httptest.NewRequest("POST", "/logout", nil),
				expectedResponse: http.Response{StatusCode: http.StatusInternalServerError},
			},
			mock: func() {
				name = "SSID"
				expires = time.Now().UTC().Add(time.Hour)
				value = bdy

				mockClient.EXPECT().
					IncUserVersion(gomock.Any(), gomock.Any()).
					Return(&generatedCommon.Empty{
						Error: "",
					}, errors.New("error"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &AuthHandler{
				client: mockClient,
				logger: zapSugar,
			}
			w := httptest.NewRecorder()
			test.mock()
			test.args.r.AddCookie(&http.Cookie{
				Name:     name,
				Value:    value,
				Expires:  expires,
				HttpOnly: true,
			})
			h.Logout(w, test.args.r)

			require.Equal(t, test.args.expectedResponse.StatusCode, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.args.expectedResponse.StatusCode, w.Code))
		})
	}
}
