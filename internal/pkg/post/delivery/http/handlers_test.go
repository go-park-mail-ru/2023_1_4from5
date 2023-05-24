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
	mockAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	mockCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/mocks"
	mockNotification "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/notification/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func bodyPrepare(like interface{}) []byte {
	userjson, err := json.Marshal(&like)
	if err != nil {
		return nil
	}
	return userjson
}

var testUser = models.User{
	Login:        "Dasha2003",
	PasswordHash: "Dasha2003!",
	Name:         "Дарья Такташова",
}

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

func TestNewPostHandler(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	creatorClient := mockCreator.NewMockCreatorServiceClient(ctl)
	notify := mockNotification.NewMockNotificationApp(ctl)
	logger, err := zap.NewProduction()
	if err != nil {
		t.Error(err.Error())
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	testHandler := NewPostHandler(authClient, creatorClient, zapSugar, notify)
	if testHandler.authClient != authClient || testHandler.creatorClient != creatorClient || testHandler.notificationApp != notify {
		t.Error("bad constructor")
	}
}

func TestPostHandler_AddLike(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	creatorClient := mockCreator.NewMockCreatorServiceClient(ctl)
	logger := zap.NewNop()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

	tests := []struct {
		name             string
		mock             func() *http.Request
		expectedResponse int
	}{
		{
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, "token")
				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name: "BadRequest wrong like format",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare("Trying to signIn")))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				return r
			},

			expectedResponse: http.StatusBadRequest,
		},
		{
			name: "Internal err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().IsPostOwner(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, errors.New("test"))

				return r
			},

			expectedResponse: http.StatusInternalServerError,
		},
		{
			name: "Internal err from is post owner",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().IsPostOwner(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "test",
				}, nil)

				return r
			},
			expectedResponse: http.StatusInternalServerError,
		},
		{
			name: "Bad Request because post owner",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().IsPostOwner(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)

				return r
			},
			expectedResponse: http.StatusBadRequest,
		},
		{
			name: "Internal Err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().IsPostOwner(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, nil)

				creatorClient.EXPECT().AddLike(gomock.Any(), gomock.Any()).Return(&generated.Like{
					LikesCount: 2,
					PostID:     uuid.New().String(),
					Error:      "",
				}, errors.New("test"))

				return r
			},

			expectedResponse: http.StatusInternalServerError,
		},
		{
			name: "Wrong data from add like",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().IsPostOwner(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, nil)

				creatorClient.EXPECT().AddLike(gomock.Any(), gomock.Any()).Return(&generated.Like{
					LikesCount: 2,
					PostID:     uuid.New().String(),
					Error:      models.WrongData.Error(),
				}, nil)

				return r
			},

			expectedResponse: http.StatusBadRequest,
		},
		{
			name: "Internal err from add like",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().IsPostOwner(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, nil)

				creatorClient.EXPECT().AddLike(gomock.Any(), gomock.Any()).Return(&generated.Like{
					LikesCount: 2,
					PostID:     uuid.New().String(),
					Error:      models.InternalError.Error(),
				}, nil)

				return r
			},

			expectedResponse: http.StatusInternalServerError,
		},
		{
			name: "Internal err from add like",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().IsPostOwner(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, nil)

				creatorClient.EXPECT().AddLike(gomock.Any(), gomock.Any()).Return(&generated.Like{
					LikesCount: 2,
					PostID:     uuid.New().String(),
					Error:      "",
				}, nil)

				return r
			},

			expectedResponse: http.StatusOK,
		},
		{
			name: "Internal err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))

				return r
			},

			expectedResponse: http.StatusInternalServerError,
		},
		{
			name: "Internal err from check user",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "test",
				}, nil)

				return r
			},

			expectedResponse: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &PostHandler{
				authClient:    authClient,
				creatorClient: creatorClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.AddLike(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func TestPostHandler_RemoveLike(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	creatorClient := mockCreator.NewMockCreatorServiceClient(ctl)
	logger := zap.NewNop()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

	tests := []struct {
		name             string
		mock             func() *http.Request
		expectedResponse int
	}{
		{
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, "token")
				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name: "Internal err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))

				return r
			},

			expectedResponse: http.StatusInternalServerError,
		},
		{
			name: "Internal err from check user",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "test",
				}, nil)

				return r
			},

			expectedResponse: http.StatusForbidden,
		},
		{
			name: "BadRequest wrong like format",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare("Trying to signIn")))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				return r
			},

			expectedResponse: http.StatusBadRequest,
		},
		{
			name: "Internal err from remove like",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().RemoveLike(gomock.Any(), gomock.Any()).Return(&generated.Like{
					LikesCount: 2,
					PostID:     uuid.New().String(),
					Error:      models.InternalError.Error(),
				}, nil)

				return r
			},

			expectedResponse: http.StatusInternalServerError,
		},
		{
			name: "Wrong data from remove like",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().RemoveLike(gomock.Any(), gomock.Any()).Return(&generated.Like{
					LikesCount: 2,
					PostID:     uuid.New().String(),
					Error:      models.WrongData.Error(),
				}, nil)

				return r
			},

			expectedResponse: http.StatusBadRequest,
		},
		{
			name: "Internal err from creator servic",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().RemoveLike(gomock.Any(), gomock.Any()).Return(&generated.Like{
					LikesCount: 2,
					PostID:     uuid.New().String(),
					Error:      "",
				}, errors.New("test"))

				return r
			},

			expectedResponse: http.StatusInternalServerError,
		},
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()})))

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)

				creatorClient.EXPECT().RemoveLike(gomock.Any(), gomock.Any()).Return(&generated.Like{
					LikesCount: 2,
					PostID:     uuid.New().String(),
					Error:      "",
				}, nil)

				return r
			},

			expectedResponse: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &PostHandler{
				authClient:    authClient,
				creatorClient: creatorClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.RemoveLike(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func TestPostHandler_CreatePost(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	os.Setenv("TOKEN_SECRET", "TEST")
	os.Setenv("CSRF_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	id := uuid.New()
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})

	logger := zap.NewNop()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	notify := mockNotification.NewMockNotificationApp(ctl)
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	creatorClient := mockCreator.NewMockCreatorServiceClient(ctl)

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/create",
					nil)

				setJWTToken(r, "111")

				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/create",
					nil)

				setJWTToken(r, bdy)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "test",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Internal err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/create",
					nil)

				setJWTToken(r, bdy)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/post/create",
					nil)

				setJWTToken(r, bdy)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Get CSRF with err",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/post/create",
					nil)
				os.Unsetenv("CSRF_SECRET")
				setJWTToken(r, bdy)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/create",
					nil)
				os.Setenv("CSRF_SECRET", "TEST")

				setJWTToken(r, bdy)
				setCSRFToken(r, "111")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Wrong data type",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/create",
					nil)

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
			name: "Wrong Multipart Fields",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte(uuid.Nil.String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Wrong data type for creator_id",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte("1323"))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Is creator internal err from service",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, errors.New("test"))
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Is creator wrong data",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: models.WrongData.Error(),
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Is creator err",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "11",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Is creator forbidden",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "No text field in multipart",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No title field in multipart",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partText, _ := writer.CreateFormField("text")
				_, err = partText.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "No title field in multipart",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partText, _ := writer.CreateFormField("text")
				_, err = partText.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "OK",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partText, _ := writer.CreateFormField("text")
				_, err = partText.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partTitle, _ := writer.CreateFormField("title")
				_, err = partTitle.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				creatorClient.EXPECT().CreatorNotificationInfo(gomock.Any(), gomock.Any()).Return(&generated.NotificationCreatorInfo{
					Name:  "test",
					Photo: uuid.New().String(),
					Error: "",
				}, nil)
				notify.EXPECT().SendUserNotification(gomock.Any(), gomock.Any()).Return(nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "InternalErr from creator info",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partText, _ := writer.CreateFormField("text")
				_, err = partText.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partTitle, _ := writer.CreateFormField("title")
				_, err = partTitle.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				creatorClient.EXPECT().CreatorNotificationInfo(gomock.Any(), gomock.Any()).Return(&generated.NotificationCreatorInfo{
					Name:  "test",
					Photo: uuid.New().String(),
					Error: "",
				}, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "InternalErr from creator info",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partText, _ := writer.CreateFormField("text")
				_, err = partText.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partTitle, _ := writer.CreateFormField("title")
				_, err = partTitle.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				creatorClient.EXPECT().CreatorNotificationInfo(gomock.Any(), gomock.Any()).Return(&generated.NotificationCreatorInfo{
					Name:  "test",
					Photo: uuid.New().String(),
					Error: "err",
				}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "err from create post service",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("creator")
				_, err := partPath.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partText, _ := writer.CreateFormField("text")
				_, err = partText.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}
				partTitle, _ := writer.CreateFormField("title")
				_, err = partTitle.Write([]byte(uuid.New().String()))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/post/create",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("test"))
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &PostHandler{
				authClient:      authClient,
				creatorClient:   creatorClient,
				logger:          zapSugar,
				notificationApp: notify,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.CreatePost(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

var newPost = models.PostEditData{Title: "test", Text: "testtest"}

func TestPostHandler_EditPost(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	os.Setenv("TOKEN_SECRET", "TEST")
	os.Setenv("CSRF_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	id := uuid.New()
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})

	logger := zap.NewNop()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	creatorClient := mockCreator.NewMockCreatorServiceClient(ctl)

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/post/edit/{post-uuid}",
					nil)

				setJWTToken(r, "111")

				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Internal err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/post/edit/{post-uuid}",
					nil)

				setJWTToken(r, bdy)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/post/edit/{post-uuid}",
					nil)

				setJWTToken(r, bdy)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "test",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/post/edit/{post-uuid}",
					nil)

				setJWTToken(r, bdy)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/post/edit/{post-uuid}",
					nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, "111")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Wrong data type",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/post/edit/{post-uuid}",
					nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				os.Setenv("CSRF_SECRET", "TEST")
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Wrong post id",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/edit/{post-uuid}",
					bytes.NewReader(bodyPrepare(newPost)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": "123",
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
			name: "Is post owner Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/edit/{post-uuid}",
					bytes.NewReader(bodyPrepare(newPost)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsPostOwner(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &PostHandler{
				authClient:    authClient,
				creatorClient: creatorClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.EditPost(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

/*
func TestPostHandler_GetPost(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	id := uuid.New()
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})

	logger := zap.NewNop()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)
	mockAttachUsecase := mockAttach.NewMockAttachmentUsecase(ctl)

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/get/{post-uuid}",
					nil)

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockPostUsecase.EXPECT().GetPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Post{}, nil)

				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "No postId",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/get/{post-uuid}",
					nil)
				setJWTToken(r, bdy)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Wrong uuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/get/{post-uuid}",
					nil)

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": "111",
				})

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/get/{post-uuid}",
					nil)

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockPostUsecase.EXPECT().GetPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Post{}, models.Forbbiden)

				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Wrong Data",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/get/{post-uuid}",
					nil)

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockPostUsecase.EXPECT().GetPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Post{}, models.WrongData)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal Error",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/get/{post-uuid}",
					nil)

				setJWTToken(r, bdy)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockPostUsecase.EXPECT().GetPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Post{}, models.InternalError)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &PostHandler{
				usecase:           mockPostUsecase,
				authUsecase:       mockAuthUsecase,
				attachmentUsecase: mockAttachUsecase,
				logger:            zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.GetPost(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}*/
