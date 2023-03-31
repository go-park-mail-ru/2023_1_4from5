package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	mockAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var testUser = models.User{
	Login:        "Dasha2003!",
	PasswordHash: "Dasha2003!",
	Name:         "Дарья Такташова",
}

func TestNewUserHandler(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock.NewMockUserUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	testHandler := NewUserHandler(mockUsecase, mockAuthUsecase)
	if testHandler.usecase != mockUsecase {
		t.Error("bad constructor")
	}
}

func TestGetProfile(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy, _ := tkn.GetJWTToken(models.User{Login: testUser.Login, Id: uuid.New()})

	usecaseMock := mock.NewMockUserUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)

	handler := NewUserHandler(usecaseMock, mockAuthUsecase)

	var r *http.Request
	var status int
	for i := 0; i < 4; i++ {
		value := bdy
		r = httptest.NewRequest("GET", "/user/profile", strings.NewReader(fmt.Sprint()))
		switch i {
		case 0:
			usecaseMock.EXPECT().GetProfile(gomock.Any()).Return(models.UserProfile{}, nil)
			status = http.StatusOK
		case 1:
			value = "body"
			status = http.StatusUnauthorized
		case 2:
			usecaseMock.EXPECT().GetProfile(gomock.Any()).Return(models.UserProfile{}, models.InternalError)
			status = http.StatusInternalServerError
		case 3:
			usecaseMock.EXPECT().GetProfile(gomock.Any()).Return(models.UserProfile{}, models.NotFound)
			status = http.StatusBadRequest
		}
		r.AddCookie(&http.Cookie{
			Name:     "SSID",
			Value:    value,
			Expires:  time.Time{},
			HttpOnly: true,
		})

		w := httptest.NewRecorder()

		handler.GetProfile(w, r)
		require.Equal(t, status, w.Code, fmt.Errorf("expected %d, got %d",
			status, w.Code))
	}
}

func TestGetHomePage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy, _ := tkn.GetJWTToken(models.User{Login: testUser.Login, Id: uuid.New()})

	usecaseMock := mock.NewMockUserUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)

	handler := NewUserHandler(usecaseMock, mockAuthUsecase)

	var r *http.Request
	var status int
	for i := 0; i < 4; i++ {
		value := bdy
		r = httptest.NewRequest("GET", "/user/homePage", strings.NewReader(fmt.Sprint()))
		switch i {
		case 0:
			usecaseMock.EXPECT().GetHomePage(gomock.Any()).Return(models.UserHomePage{}, nil)
			status = http.StatusOK
		case 1:
			value = "body"
			status = http.StatusUnauthorized
		case 2:
			usecaseMock.EXPECT().GetHomePage(gomock.Any()).Return(models.UserHomePage{}, models.InternalError)
			status = http.StatusInternalServerError
		case 3:
			usecaseMock.EXPECT().GetHomePage(gomock.Any()).Return(models.UserHomePage{}, models.NotFound)
			status = http.StatusBadRequest
		}
		r.AddCookie(&http.Cookie{
			Name:     "SSID",
			Value:    value,
			Expires:  time.Time{},
			HttpOnly: true,
		})

		w := httptest.NewRecorder()

		handler.GetHomePage(w, r)
		require.Equal(t, status, w.Code, fmt.Errorf("expected %d, got %d",
			status, w.Code))
	}

}
