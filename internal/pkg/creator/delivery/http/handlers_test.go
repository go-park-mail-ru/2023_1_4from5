package http

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/mocks"
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

var testUsers []models.User = []models.User{
	{
		Login:        "Dasha2003!",
		PasswordHash: "Dasha2003!",
		Name:         "Дарья Такташова",
	},
}

func TestGetPage(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy := tkn.GetToken(models.User{Login: testUsers[0].Login, Id: uuid.New()})

	usecaseMock := mock.NewMockCreatorUsecase(ctl)

	handler := NewCreatorHandler(usecaseMock)
	var r *http.Request
	var status int
	for i := 0; i < 4; i++ {
		value := bdy
		r = httptest.NewRequest("GET", "/creator/page", strings.NewReader(fmt.Sprint()))
		switch i {
		case 0:
			usecaseMock.EXPECT().GetPage(gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, nil)
			r = mux.SetURLVars(r, map[string]string{
				"creator-uuid": uuid.NewString(),
			})
			status = http.StatusOK
		case 1:
			r = mux.SetURLVars(r, map[string]string{
				"creator_uuid": uuid.NewString(),
			})
			status = http.StatusBadRequest
		case 2:
			value = "body"
			r = mux.SetURLVars(r, map[string]string{
				"creator-uuid": uuid.NewString(),
			})
			status = http.StatusUnauthorized
		case 3:
			usecaseMock.EXPECT().GetPage(gomock.Any(), gomock.Any()).Return(models.CreatorPage{}, errors.New("test"))
			r = mux.SetURLVars(r, map[string]string{
				"creator-uuid": uuid.NewString(),
			})
			status = http.StatusInternalServerError
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
