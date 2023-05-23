package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	mockAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	mockCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	mockUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func setCSRFToken(r *http.Request, token string) {
	r.Header.Set("X-CSRF-Token", token)
}

func setJWTToken(r *http.Request, token string) {
	r.AddCookie(&http.Cookie{
		Name:     "SSID",
		Value:    token,
		Expires:  time.Time{},
		HttpOnly: true,
	})
}

var testComment = models.Comment{
	CommentID:  uuid.New(),
	UserID:     uuid.New(),
	Username:   "TestUsername",
	UserPhoto:  uuid.New(),
	PostID:     uuid.New(),
	Text:       "TestText",
	Creation:   time.Now(),
	LikesCount: 10,
	IsLiked:    false,
	IsOwner:    true,
}

var testUser = models.User{
	Login:        "Dasha2003!",
	PasswordHash: "Dasha2003!",
	Name:         "Дарья Такташова",
}

func bodyPrepare(val interface{}) []byte {
	valJSON, err := json.Marshal(&val)
	if err != nil {
		return nil
	}
	return valJSON
}

//func TestNewCommentHandler(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//
//	authClient := mockAuth.NewMockAuthServiceClient(ctl)
//	creatorClient := mockCreator.NewMockCreatorServiceClient(ctl)
//	userClient := mockUser.NewMockUserServiceClient(ctl)
//
//	logger := zap.NewNop()
//
//	defer func(logger *zap.Logger) {
//		err := logger.Sync()
//		if err != nil {
//			return
//		}
//	}(logger)
//	zapSugar := logger.Sugar()
//
//	testHandler := NewCreatorHandler(creatorClient, authClient, notify, zapSugar)
//	if testHandler.authClient != authClient || testHandler.creatorClient != creatorClient || testHandler.notificationApp != notify {
//		t.Error("bad constructor")
//	}
//}

func TestCommentHandler_CreateComment(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	creatorClient := mockCreator.NewMockCreatorServiceClient(ctl)
	userClient := mockUser.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	id := uuid.New()
	os.Setenv("CSRF_SECRET", "TEST")
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/comment/create",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, nil)
				creatorClient.EXPECT().CreateComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		//{
		//	name: "Get CSRF",
		//	mock: func() *http.Request {
		//		r := httptest.NewRequest("GET", "/comment/create", nil)
		//
		//		setJWTToken(r, bdy)
		//		os.Setenv("CSRF_SECRET", "TEST")
		//
		//		authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
		//			UserVersion: int64(1),
		//			Error:       "",
		//		}, nil)
		//		return r
		//	},
		//	expectedStatus: http.StatusOK,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CommentHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				userClient:    userClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.CreateComment(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}
