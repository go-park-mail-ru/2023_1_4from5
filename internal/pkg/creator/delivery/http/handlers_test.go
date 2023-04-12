package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	mockAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/mocks"
	mockPost "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func TestNewCreatorHandler(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock.NewMockCreatorUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)
	testHandler := NewCreatorHandler(mockUsecase, mockAuthUsecase, mockPostUsecase)
	if testHandler.usecase != mockUsecase {
		t.Error("bad constructor")
	}
}

func TestCreatorHandler_GetPage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

	usecaseMock := mock.NewMockCreatorUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)
	handler := NewCreatorHandler(usecaseMock, mockAuthUsecase, mockPostUsecase)
	var r *http.Request
	var status int
	for i := 0; i < 4; i++ {
		value := bdy
		r = httptest.NewRequest("GET", "/creator/page", strings.NewReader(fmt.Sprint()))
		r = mux.SetURLVars(r, map[string]string{
			"creator-uuid": uuid.NewString(),
		})
		switch i {
		case 0:
			usecaseMock.EXPECT().GetPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, nil)
			status = http.StatusOK
		case 1:
			r = mux.SetURLVars(r, map[string]string{
				"creator_uuid": uuid.NewString(),
			})
			status = http.StatusBadRequest
		case 2:
			value = "1"
			usecaseMock.EXPECT().GetPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, models.InternalError)
			status = http.StatusInternalServerError
		case 3:
			usecaseMock.EXPECT().GetPage(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, models.WrongData)
			status = http.StatusBadRequest
		}
		r.AddCookie(&http.Cookie{
			Name:     "SSID",
			Value:    value,
			Expires:  time.Time{},
			HttpOnly: true,
		})

		w := httptest.NewRecorder()

		handler.GetPage(w, r)
		require.Equal(t, status, w.Code, fmt.Errorf("expected %d, got %d",
			status, w.Code))
	}

}

var testAim = models.Aim{Creator: uuid.New(), Description: "test", MoneyNeeded: 500, MoneyGot: 0}
var testAimWithLongDescription = models.Aim{Creator: uuid.New(), Description: "testtesttesttesttesttesttesttesttesttesttesttesttesttesttest" +
	"testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest" +
	"testtesttesttesttesttesttest", MoneyNeeded: 500, MoneyGot: 0}

func bodyPrepare(aim models.Aim) []byte {
	aimJSON, err := json.Marshal(&aim)
	if err != nil {
		return nil
	}
	return aimJSON
}

func TestCreatorHandler_CreateAim(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

	usecaseMock := mock.NewMockCreatorUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				usecaseMock.EXPECT().CreateAim(gomock.Any(), gomock.Any()).Return(nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Wrong Token",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    "1",
					Expires:  time.Time{},
					HttpOnly: true,
				})

				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden, wrong UserVersion",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Bad Request",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader([]byte("11111")))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Bad request: wrong description length",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAimWithLongDescription)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal Error",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				usecaseMock.EXPECT().CreateAim(gomock.Any(), gomock.Any()).Return(models.InternalError)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Is Creator wrong data",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, models.WrongData)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Is Creator internal error",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, models.InternalError)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Is Creator wrong",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				usecase:     usecaseMock,
				authUsecase: mockAuthUsecase,
				postUsecase: mockPostUsecase,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.CreateAim(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

var testCreators = make([]models.Creator, 1)

func creatorsBodyPrepare(status int, creator ...models.Creator) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	if len(creator) == 0 {
		utils.Response(w, status, nil)
		return w
	}
	utils.Response(w, status, creator)
	return w
}

func TestCreatorHandler_GetAllCreators(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	usecaseMock := mock.NewMockCreatorUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
		expectedBody   []models.Creator
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create", nil)

				usecaseMock.EXPECT().GetAllCreators(gomock.Any()).Return(testCreators, nil)
				return r
			},
			expectedStatus: http.StatusOK,
			expectedBody:   testCreators,
		},
		{
			name: "InternalError",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create", nil)

				usecaseMock.EXPECT().GetAllCreators(gomock.Any()).Return(nil, models.InternalError)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				usecase:     usecaseMock,
				authUsecase: mockAuthUsecase,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.GetAllCreators(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
			require.Equal(t, creatorsBodyPrepare(test.expectedStatus, test.expectedBody...).Body.String(), w.Body.String(), fmt.Errorf("Wrong body"))
		})
	}
}
