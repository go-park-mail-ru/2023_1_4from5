package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	mockAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	mockCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	mockUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func TestNewCommentHandler(t *testing.T) {
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

	testHandler := NewCommentHandler(authClient, userClient, creatorClient, zapSugar)
	if testHandler.authClient != authClient || testHandler.creatorClient != creatorClient || testHandler.userClient != userClient {
		t.Error("bad constructor")
	}
}

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
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/comment/create", nil)

				setJWTToken(r, bdy)
				os.Unsetenv("CSRF_SECRET")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/comment/create", nil)

				setJWTToken(r, bdy)
				os.Setenv("CSRF_SECRET", "TEST")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "No CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/create", nil)

				setJWTToken(r, bdy)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Err from Check User Version",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", nil)

				setJWTToken(r, bdy)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "1",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Internal Err from Check User Version",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", nil)

				setJWTToken(r, bdy)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("err"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Wrong Token",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/comment/create",
					bytes.NewReader(bodyPrepare(testComment)))

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
			name: "PostUnavailable",
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
					Error: "testErr",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "PostUnavailable Internal",
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
				}, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "bad comment err",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/comment/create",
					bytes.NewReader(bodyPrepare(models.Comment{
						CommentID:  uuid.UUID{},
						UserID:     uuid.UUID{},
						Username:   "",
						UserPhoto:  uuid.UUID{},
						PostID:     uuid.UUID{},
						Text:       "",
						Creation:   time.Time{},
						LikesCount: 0,
						IsLiked:    false,
						IsOwner:    false,
					})))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "CreateComment err",
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
				creatorClient.EXPECT().CreateComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: "error"}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "CreateComment err Internal",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/comment/create",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil).AnyTimes()
				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, nil)
				creatorClient.EXPECT().CreateComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("internal"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
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
				}, nil).AnyTimes()
				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, nil).AnyTimes()
				creatorClient.EXPECT().CreateComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/comment/create", nil)

				setJWTToken(r, bdy)
				os.Unsetenv("CSRF_SECRET")

				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
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

func TestCommentHandler_DeleteComment(t *testing.T) {
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
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/comment/delete/b72dd39d-e19b-4070-9200-71a0c92417ca", nil)

				setJWTToken(r, "token")

				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/comment/delete/b72dd39d-e19b-4070-9200-71a0c92417ca", nil)

				setJWTToken(r, bdy)
				os.Unsetenv("CSRF_SECRET")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/comment/delete/b72dd39d-e19b-4070-9200-71a0c92417ca", nil)

				setJWTToken(r, bdy)
				os.Setenv("CSRF_SECRET", "TEST")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "No CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/b72dd39d-e19b-4070-9200-71a0c92417ca", nil)

				setJWTToken(r, bdy)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "ErrUV",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "err",
				}, nil)

				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "ErrUVInternal",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("err"))

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Err bad uuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71aewewe0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "CommentOwner err internal",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "",
				}, errors.New("err"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "CommentOwner err wd",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "WrongData",
				}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "CommentOwner err wd",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "err",
				}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "CommentOwner err forbiden",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  false,
					Error: "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Wrong Data",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: "WrongData"}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal(micro)",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("err"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: "err"}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/comment/delete/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
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
			h.DeleteComment(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestCommentHandler_EditComment(t *testing.T) {
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
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/comment/edit/b72dd39d-e19b-4070-9200-71a0c92417ca", nil)

				setJWTToken(r, "token")

				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/comment/edit/b72dd39d-e19b-4070-9200-71a0c92417ca", nil)

				setJWTToken(r, bdy)
				os.Unsetenv("CSRF_SECRET")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/comment/edit/b72dd39d-e19b-4070-9200-71a0c92417ca", nil)

				setJWTToken(r, bdy)
				os.Setenv("CSRF_SECRET", "TEST")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "No CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/b72dd39d-e19b-4070-9200-71a0c92417ca", nil)

				setJWTToken(r, bdy)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "ErrUV",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "err",
				}, nil)

				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "ErrUVInternal",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("err"))

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Err bad uuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71aewewe0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Err no uuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "CommentOwner err internal",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "",
				}, errors.New("err"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "CommentOwner err wd",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "WrongData",
				}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "CommentOwner err wd",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "err",
				}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "CommentOwner err forbiden",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  false,
					Error: "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "out err",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().EditComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: "err"}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "internal edit err",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().EditComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("error"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/edit/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCommentOwner(gomock.Any(), gomock.Any()).Return(&generatedCreator.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().EditComment(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
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
			h.EditComment(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestCommentHandler_AddLike(t *testing.T) {
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
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "unauth",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/addLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, "token")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "badUuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/addLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, nil)
				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c9ewe2417ca",
				})

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "forbidden post",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/addLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "err",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "internal ispostavaliable",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/addLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, errors.New("err"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "no uuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/addLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal addLike",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/addLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, nil)
				creatorClient.EXPECT().AddLikeComment(gomock.Any(), gomock.Any()).Return(&generatedCreator.Like{Error: ""}, errors.New("err"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "WrongData addLike",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/addLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, nil)
				creatorClient.EXPECT().AddLikeComment(gomock.Any(), gomock.Any()).Return(&generatedCreator.Like{Error: "WrongData"}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal addLike",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/addLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, nil)
				creatorClient.EXPECT().AddLikeComment(gomock.Any(), gomock.Any()).Return(&generatedCreator.Like{Error: "err"}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/addLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().IsPostAvailable(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{
					Error: "",
				}, nil)
				creatorClient.EXPECT().AddLikeComment(gomock.Any(), gomock.Any()).Return(&generatedCreator.Like{Error: ""}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
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
			h.AddLike(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestCommentHandler_RemoveLike(t *testing.T) {
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
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "unauth",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/removeLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, "token")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "badUuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/removeLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c9ewe2417ca",
				})

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "noUuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/removeLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Bad Request",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/removeLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().RemoveLikeComment(gomock.Any(), gomock.Any()).Return(&generatedCreator.Like{Error: "WrongData"}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal DB",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/removeLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().RemoveLikeComment(gomock.Any(), gomock.Any()).Return(&generatedCreator.Like{Error: "err"}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal MS",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/removeLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().RemoveLikeComment(gomock.Any(), gomock.Any()).Return(&generatedCreator.Like{Error: ""}, errors.New("err"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/comment/removeLike/{comment-uuid}",
					bytes.NewReader(bodyPrepare(testComment)))

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"comment-uuid": "b72dd39d-e19b-4070-9200-71a0c92417ca",
				})

				creatorClient.EXPECT().RemoveLikeComment(gomock.Any(), gomock.Any()).Return(&generatedCreator.Like{Error: ""}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
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
			h.RemoveLike(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}
