package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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

type fields struct {
	Usecase auth.AuthUsecase
}

type args struct {
	r                  *http.Request
	expectedResponse   http.Response
	expectedStatusCode error
}

var testUsers []models.User = []models.User{
	{
		Login:        "Dasha2003!",
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

	mockUsecase := mock.NewMockAuthUsecase(ctl)
	testHandler := NewAuthHandler(mockUsecase)
	if testHandler.usecase != mockUsecase {
		t.Error("bad constructor")
	}
}

func TestAuthHandler_SignIn(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock.NewMockAuthUsecase(ctl)

	tests := []struct {
		name   string
		Login  string
		fields fields
		args   args
	}{
		{
			name:   "OK",
			Login:  testUsers[0].Login,
			fields: fields{Usecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("POST", "/signIn",
					bytes.NewReader(bodyPrepare(testUsers[0]))),
				expectedStatusCode: nil,
				expectedResponse:   http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name:   "Unauthorized",
			Login:  testUsers[1].Login,
			fields: fields{Usecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("POST", "/signIn",
					bytes.NewReader(bodyPrepare(testUsers[1]))),
				expectedStatusCode: models.NotFound,
				expectedResponse:   http.Response{StatusCode: http.StatusUnauthorized},
			},
		},
		{
			name:   "BadRequest",
			Login:  testUsers[2].Login,
			fields: fields{Usecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("POST", "//signIn",
					bytes.NewReader([]byte("Trying to signIN"))),
				expectedStatusCode: models.NoAuthData,
				expectedResponse:   http.Response{StatusCode: http.StatusBadRequest},
			},
		},
	}

	for i := 0; i < len(tests); i++ {
		LoginUserCopy := models.LoginUser{Login: testUsers[i].Login, PasswordHash: testUsers[i].PasswordHash}
		if tests[i].args.expectedStatusCode != models.NoAuthData {
			mockUsecase.EXPECT().
				SignIn(LoginUserCopy).
				Return("", tests[i].args.expectedStatusCode)
		}
	}

	for _, test := range tests {
		t.Run(test.Login, func(t *testing.T) {
			h := &AuthHandler{
				usecase: test.fields.Usecase,
			}
			w := httptest.NewRecorder()

			h.SignIn(w, test.args.r)
			require.Equal(t, test.args.expectedResponse.StatusCode, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for login:%s", test.name, test.args.expectedResponse.StatusCode, w.Code, test.Login))
		})
	}
}

func TestAuthHandler_SignUp(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock.NewMockAuthUsecase(ctl)

	tests := []struct {
		name   string
		Login  string
		fields fields
		args   args
	}{
		{
			name:   "OK",
			Login:  testUsers[0].Login,
			fields: fields{Usecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("POST", "/signUp",
					bytes.NewReader(bodyPrepare(testUsers[0]))),
				expectedStatusCode: nil,
				expectedResponse:   http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name:   "Conflict",
			Login:  testUsers[1].Login,
			fields: fields{Usecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("POST", "/signUp",
					bytes.NewReader(bodyPrepare(testUsers[1]))),
				expectedStatusCode: models.ConflictData,
				expectedResponse:   http.Response{StatusCode: http.StatusConflict},
			},
		},
		{
			name:   "BadRequest because no json",
			Login:  testUsers[2].Login,
			fields: fields{Usecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("POST", "/signUp",
					bytes.NewReader([]byte("ppppp"))),
				expectedStatusCode: models.NoAuthData,
				expectedResponse:   http.Response{StatusCode: http.StatusBadRequest},
			},
		},
		{
			name:   "BadRequest because user is not valid",
			Login:  testUsers[2].Login,
			fields: fields{Usecase: mockUsecase},
			args: args{
				r: httptest.NewRequest("POST", "/signUp",
					bytes.NewReader(bodyPrepare(testUsers[3]))),
				expectedStatusCode: models.NoAuthData,
				expectedResponse:   http.Response{StatusCode: http.StatusBadRequest},
			},
		},
	}

	for i := 0; i < len(tests); i++ {
		if tests[i].args.expectedStatusCode == models.NoAuthData {
			continue
		}
		if tests[i].args.expectedStatusCode == nil {
			mockUsecase.EXPECT().SignUp(testUsers[i]).Return("token", tests[i].args.expectedStatusCode)
			continue
		}
		mockUsecase.EXPECT().
			SignUp(models.User{Login: testUsers[i].Login, PasswordHash: testUsers[i].PasswordHash, Name: testUsers[i].Name}).
			Return("", tests[i].args.expectedStatusCode)
	}

	for _, tt := range tests {
		t.Run(tt.Login, func(t *testing.T) {
			h := &AuthHandler{
				usecase: tt.fields.Usecase,
			}
			w := httptest.NewRecorder()

			h.SignUp(w, tt.args.r)
			require.Equal(t, tt.args.expectedResponse.StatusCode, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for login:%s", tt.name, tt.args.expectedResponse.StatusCode, w.Code, tt.Login))
		})
	}
}

type argsLogout struct {
	r                  *http.Request
	expectedResponse   http.Response
	expectedStatusCode int
}

func TestAuthHandler_Logout(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy := tkn.GetToken(models.User{Login: testUsers[1].Login, Id: uuid.New()})

	tests := []struct {
		name  string
		Login string
		body  []byte
		args  argsLogout
	}{
		{
			name:  "BadRequest wrong value for token",
			Login: testUsers[0].Login,
			args: argsLogout{
				r:                  httptest.NewRequest("POST", "/logout", nil),
				expectedStatusCode: http.StatusBadRequest,
				expectedResponse:   http.Response{StatusCode: http.StatusBadRequest},
			},
		},
		{
			name:  "OK",
			Login: testUsers[1].Login,
			args: argsLogout{
				r:                  httptest.NewRequest("POST", "/logout", nil),
				expectedStatusCode: http.StatusOK,
				expectedResponse:   http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name:  "BadRequest wrong Cookie name",
			Login: testUsers[2].Login,
			args: argsLogout{
				r:                  httptest.NewRequest("POST", "/logout", nil),
				expectedStatusCode: http.StatusBadRequest,
				expectedResponse:   http.Response{StatusCode: http.StatusBadRequest},
			},
		},
	}

	for i := 0; i < len(tests); i++ {
		name := "SSID"
		expires := time.Now().Add(time.Hour)
		value := bdy
		if tests[i].args.expectedStatusCode == http.StatusBadRequest {
			switch i {
			case 0:
				value = "token"
			case 2:
				name = "ssss"
			}
			fmt.Println(expires.Unix())
			tests[i].args.r.AddCookie(&http.Cookie{
				Name:     name,
				Value:    value,
				Expires:  expires,
				HttpOnly: true,
			})
			continue
		}
		tests[i].args.r.AddCookie(&http.Cookie{
			Name:     "SSID",
			Value:    bdy,
			Expires:  time.Time{},
			HttpOnly: true,
		})
	}

	for _, tt := range tests {
		t.Run(tt.Login, func(t *testing.T) {
			h := &AuthHandler{}

			w := httptest.NewRecorder()
			h.Logout(w, tt.args.r)
			require.Equal(t, tt.args.expectedResponse.StatusCode, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for login:%s", tt.name, tt.args.expectedResponse.StatusCode, w.Code, tt.Login))
		})
	}
}
