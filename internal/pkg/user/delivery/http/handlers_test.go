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
	mockNotification "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/notification/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var follow = []*generated.Follow{{
	Creator:      uuid.New().String(),
	CreatorName:  "test",
	CreatorPhoto: uuid.New().String(),
	Description:  "test",
}}

var followWithErr1 = []*generated.Follow{{
	Creator:      "11",
	CreatorName:  "test",
	CreatorPhoto: uuid.New().String(),
	Description:  "test",
}}

var followWithErr2 = []*generated.Follow{{
	Creator:      uuid.New().String(),
	CreatorName:  "test",
	CreatorPhoto: "11",
	Description:  "test",
}}

var testUser = models.User{
	Login:        "Dasha2003!",
	PasswordHash: "Dasha2003!",
	Name:         "Дарья Такташова",
}

func TestNewUserHandler(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)
	notify := mockNotification.NewMockNotificationApp(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	testHandler := NewUserHandler(userClient, authClient, notify, zapSugar)
	if testHandler.userClient != userClient || testHandler.authClient != authClient {
		t.Error("bad constructor")
	}
}

func TestUserHandler_UserSubscriptions(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name             string
		expectedResponse int
		mock             func() *http.Request
	}{
		{
			name:             "OK",
			expectedResponse: http.StatusOK,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/subscriptions", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, token)
				userClient.EXPECT().
					UserSubscriptions(gomock.Any(), gomock.Any()).
					Return(&generated.SubscriptionsMessage{
						Subscriptions: []*generatedCommon.Subscription{{
							Id:           uuid.New().String(),
							Creator:      uuid.New().String(),
							CreatorName:  uuid.New().String(),
							CreatorPhoto: uuid.New().String(),
							MonthCost:    100,
							Title:        "test",
							Description:  "test",
						}},
						Error: "",
					}, nil)
				return r

			},
		},
		{
			name:             "Unauthorized",
			expectedResponse: http.StatusUnauthorized,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/subscriptions", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, "1")
				return r
			},
		},
		{
			name:             "Error from user service",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/subscriptions", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, token)
				userClient.EXPECT().
					UserSubscriptions(gomock.Any(), gomock.Any()).
					Return(&generated.SubscriptionsMessage{
						Subscriptions: []*generatedCommon.Subscription{{
							Id:           uuid.New().String(),
							Creator:      uuid.New().String(),
							CreatorName:  uuid.New().String(),
							CreatorPhoto: uuid.New().String(),
							MonthCost:    100,
							Title:        "test",
							Description:  "test",
						}},
						Error: "",
					}, errors.New("test"))
				return r

			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()

			h.UserSubscriptions(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func bodyPrepare(val interface{}) []byte {
	valJSON, err := json.Marshal(&val)
	if err != nil {
		return nil
	}
	return valJSON
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

func TestGetProfile(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name             string
		expectedResponse int
		mock             func() *http.Request
	}{
		{
			name:             "OK",
			expectedResponse: http.StatusOK,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userProfile", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, token)
				userClient.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generated.UserProfile{
						Login:        "testLogin",
						Name:         "testName",
						ProfilePhoto: uuid.New().String(),
						Registration: "2006-01-02 15:04:05 -0700 -0700",
						IsCreator:    false,
						CreatorID:    uuid.New().String(),
						Error:        "",
					}, nil)
				return r

			},
		},
		{
			name:             "Unauthorized",
			expectedResponse: http.StatusUnauthorized,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userProfile", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, "1")
				return r
			},
		},
		{
			name:             "Error from user service",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userProfile", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, token)
				userClient.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generated.UserProfile{
						Login:        "testLogin",
						Name:         "testName",
						ProfilePhoto: uuid.New().String(),
						Registration: "2006-01-02 15:04:05 -0700 -0700",
						IsCreator:    false,
						CreatorID:    uuid.New().String(),
						Error:        "",
					}, errors.New("test"))
				return r

			},
		},
		{
			name:             "Not Found Error from get profile",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userProfile", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, token)
				userClient.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generated.UserProfile{
						Login:        "testLogin",
						Name:         "testName",
						ProfilePhoto: uuid.New().String(),
						Registration: "2006-01-02 15:04:05 -0700 -0700",
						IsCreator:    false,
						CreatorID:    uuid.New().String(),
						Error:        models.NotFound.Error(),
					}, nil)
				return r
			},
		},
		{
			name:             "Internal Error from get profile",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userProfile", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, token)
				userClient.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generated.UserProfile{
						Login:        "testLogin",
						Name:         "testName",
						ProfilePhoto: uuid.New().String(),
						Registration: "2006-01-02 15:04:05 -0700 -0700",
						IsCreator:    false,
						CreatorID:    uuid.New().String(),
						Error:        models.InternalError.Error(),
					}, nil)
				return r
			},
		},
		{
			name:             "Internal Error from date format",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userProfile", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, token)
				userClient.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generated.UserProfile{
						Login:        "testLogin",
						Name:         "testName",
						ProfilePhoto: uuid.New().String(),
						Registration: "2006-01-02",
						IsCreator:    false,
						CreatorID:    uuid.New().String(),
						Error:        "",
					}, nil)
				return r
			},
		},
		{
			name:             "Internal Error from uuid1",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userProfile", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, token)
				userClient.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generated.UserProfile{
						Login:        "testLogin",
						Name:         "testName",
						ProfilePhoto: "1",
						Registration: time.Now().Format(time.RFC3339),
						IsCreator:    false,
						CreatorID:    uuid.New().String(),
						Error:        "",
					}, nil)
				return r
			},
		},
		{
			name:             "Internal Error from uuid2",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userProfile", strings.NewReader(fmt.Sprint()))
				setJWTToken(r, token)
				userClient.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generated.UserProfile{
						Login:        "testLogin",
						Name:         "testName",
						ProfilePhoto: uuid.New().String(),
						Registration: "2006-01-02 15:04:05 -0700 -0700",
						IsCreator:    false,
						CreatorID:    "1",
						Error:        "",
					}, nil)
				return r
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()

			h.GetProfile(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedResponse, w.Code))
		})
	}
}
func TestUpdatePassword(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	id := uuid.New()
	os.Setenv("CSRF_SECRET", "TEST")
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name             string
		expectedResponse int
		mock             func() *http.Request
	}{
		{
			name:             "OK",
			expectedResponse: http.StatusOK,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!!!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				authClient.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(&generatedAuth.User{}, nil)
				authClient.EXPECT().EncryptPwd(gomock.Any(), gomock.Any()).Return(&generatedAuth.EncryptPwdMg{Password: "testpassaasasasda"}, nil)
				userClient.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				authClient.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(&generatedAuth.Token{
					Cookie: "token",
					Error:  "",
				}, nil)
				return r

			},
		},
		{
			name:             "err from auth service",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!!!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))
				return r
			},
		},
		{
			name:             "err from auth service2",
			expectedResponse: http.StatusForbidden,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!!!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "test",
				}, nil)
				return r
			},
		},
		{
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/UpdatePassword", nil)

				setJWTToken(r, token)
				os.Unsetenv("CSRF_SECRET")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/UpdatePassword", nil)

				setJWTToken(r, token)
				os.Setenv("CSRF_SECRET", "TEST")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedResponse: http.StatusOK,
		},
		{
			name: "Wrong Token",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword",
					nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    "1",
					Expires:  time.Time{},
					HttpOnly: true,
				})

				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name:             "Err from signIn",
			expectedResponse: http.StatusUnauthorized,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!!!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				authClient.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(&generatedAuth.User{}, nil)
				authClient.EXPECT().EncryptPwd(gomock.Any(), gomock.Any()).Return(&generatedAuth.EncryptPwdMg{Password: "testpassaasasasda"}, nil)
				userClient.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				authClient.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(&generatedAuth.Token{
					Cookie: "token",
					Error:  "err",
				}, nil)
				return r

			},
		},
		{
			name:             "Internal err from auth service",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!!!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				authClient.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(&generatedAuth.User{}, nil)
				authClient.EXPECT().EncryptPwd(gomock.Any(), gomock.Any()).Return(&generatedAuth.EncryptPwdMg{Password: "testpassaasasasda"}, nil)
				userClient.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				authClient.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(&generatedAuth.Token{
					Cookie: "token",
					Error:  "",
				}, errors.New("test"))
				return r

			},
		},
		{
			name:             "Internal err from update password",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!!!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				authClient.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(&generatedAuth.User{}, nil)
				authClient.EXPECT().EncryptPwd(gomock.Any(), gomock.Any()).Return(&generatedAuth.EncryptPwdMg{Password: "testpassaasasasda"}, nil)
				userClient.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: "test"}, nil)
				return r
			},
		},
		{
			name:             "Internal err from user service",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!!!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				authClient.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(&generatedAuth.User{}, nil)
				authClient.EXPECT().EncryptPwd(gomock.Any(), gomock.Any()).Return(&generatedAuth.EncryptPwdMg{Password: "testpassaasasasda"}, nil)
				userClient.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("test"))
				return r
			},
		},
		{
			name:             "Err from encrypt pwd",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!!!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				authClient.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(&generatedAuth.User{}, nil)
				authClient.EXPECT().EncryptPwd(gomock.Any(), gomock.Any()).Return(&generatedAuth.EncryptPwdMg{Password: "testpassaasasasda"}, errors.New("test"))
				return r
			},
		},
		{
			name:             "Error from check user",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!!!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				authClient.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(&generatedAuth.User{}, errors.New("test"))
				return r

			},
		},
		{
			name:             "Same passwords",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader(bodyPrepare(models.UpdatePasswordInfo{
					NewPassword: "Dasha3003!",
					OldPassword: "Dasha3003!",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r

			},
		},
		{
			name:             "Wrong data",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/UpdatePassword", bytes.NewReader([]byte("11")))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()

			h.UpdatePassword(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func TestUserHandler_Follow(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	id := uuid.New()
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name             string
		expectedResponse int
		mock             func() *http.Request
	}{
		{
			name:             "OK",
			expectedResponse: http.StatusOK,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.New().String(),
				})
				userClient.EXPECT().Follow(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				return r
			},
		},
		{
			name:             "err from user service",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.New().String(),
				})
				userClient.EXPECT().Follow(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("test"))
				return r
			},
		},
		{
			name:             "err from follow",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.New().String(),
				})
				userClient.EXPECT().Follow(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: models.InternalError.Error()}, nil)
				return r
			},
		},
		{
			name:             "err from follow",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.New().String(),
				})
				userClient.EXPECT().Follow(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: models.WrongData.Error()}, nil)
				return r
			},
		},
		{
			name:             "Wrong uuid",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": "11",
				})
				return r
			},
		},
		{
			name:             "no uuid",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				return r
			},
		},
		{
			name: "Wrong Token",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow",
					nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    "1",
					Expires:  time.Time{},
					HttpOnly: true,
				})

				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()

			h.Follow(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func TestUserHandler_Unfollow(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	id := uuid.New()
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name             string
		expectedResponse int
		mock             func() *http.Request
	}{
		{
			name:             "OK",
			expectedResponse: http.StatusOK,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.New().String(),
				})
				userClient.EXPECT().Unfollow(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				return r
			},
		},
		{
			name:             "err from user service",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.New().String(),
				})
				userClient.EXPECT().Unfollow(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("test"))
				return r
			},
		},
		{
			name:             "err from follow",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.New().String(),
				})
				userClient.EXPECT().Unfollow(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: models.InternalError.Error()}, nil)
				return r
			},
		},
		{
			name:             "err from follow",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.New().String(),
				})
				userClient.EXPECT().Unfollow(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: models.WrongData.Error()}, nil)
				return r
			},
		},
		{
			name:             "Wrong uuid",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": "11",
				})
				return r
			},
		},
		{
			name:             "no uuid",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow", nil)
				setJWTToken(r, token)
				return r
			},
		},
		{
			name: "Wrong Token",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/Follow",
					nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    "1",
					Expires:  time.Time{},
					HttpOnly: true,
				})

				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()

			h.Unfollow(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func TestUserHandler_UpdateData(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	id := uuid.New()
	os.Setenv("CSRF_SECRET", "TEST")
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name             string
		expectedResponse int
		mock             func() *http.Request
	}{
		{
			name:             "Err from CheckUserVersion",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateData", bytes.NewReader(bodyPrepare(models.Donate{
					MoneyCount: 100,
					CreatorID:  uuid.New(),
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))
				return r

			},
		},
		{
			name:             "Err from CheckUserVersion",
			expectedResponse: http.StatusForbidden,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateData", bytes.NewReader(bodyPrepare(models.Donate{
					MoneyCount: 100,
					CreatorID:  uuid.New(),
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "test",
				}, nil)
				return r

			},
		},
		{
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/updateData", nil)

				setJWTToken(r, token)
				os.Unsetenv("CSRF_SECRET")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/updateData", nil)

				setJWTToken(r, token)
				os.Setenv("CSRF_SECRET", "TEST")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedResponse: http.StatusOK,
		},
		{
			name: "Wrong Token",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateData",
					nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    "1",
					Expires:  time.Time{},
					HttpOnly: true,
				})

				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name:             "Wrong data",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateData", bytes.NewReader([]byte("11")))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r

			},
		},
		{
			name:             "OK",
			expectedResponse: http.StatusOK,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateData", bytes.NewReader(bodyPrepare(models.UpdateProfileInfo{
					Login: "TestLogin",
					Name:  "New TestName",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				userClient.EXPECT().UpdateProfileInfo(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
				return r
			},
		},
		{
			name:             "Err from user service",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateData", bytes.NewReader(bodyPrepare(models.UpdateProfileInfo{
					Login: "TestLogin",
					Name:  "New TestName",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				userClient.EXPECT().UpdateProfileInfo(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("test"))
				return r
			},
		},
		{
			name:             "Err from updateProfileInfo",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateData", bytes.NewReader(bodyPrepare(models.UpdateProfileInfo{
					Login: "TestLogin",
					Name:  "New TestName",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				userClient.EXPECT().UpdateProfileInfo(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: "test"}, nil)
				return r
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()

			h.UpdateData(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func TestUserHandler_BecomeCreator(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	id := uuid.New()
	os.Setenv("CSRF_SECRET", "TEST")
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name             string
		expectedResponse int
		mock             func() *http.Request
	}{
		{
			name:             "Err from CheckUserVersion",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", nil)
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))
				return r

			},
		},
		{
			name:             "Err from CheckUserVersion",
			expectedResponse: http.StatusForbidden,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", nil)
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "test",
				}, nil)
				return r

			},
		},
		{
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/becomeCreator", nil)

				setJWTToken(r, token)
				os.Unsetenv("CSRF_SECRET")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/becomeCreator", nil)

				setJWTToken(r, token)
				os.Setenv("CSRF_SECRET", "TEST")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedResponse: http.StatusOK,
		},
		{
			name: "Wrong Token",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator",
					nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    "1",
					Expires:  time.Time{},
					HttpOnly: true,
				})

				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name:             "Wrong data",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", bytes.NewReader([]byte("11")))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				userClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generated.CheckCreatorMessage{
					ID:        uuid.New().String(),
					IsCreator: false,
					Error:     "",
				}, nil)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r

			},
		},
		{
			name:             "Check if Creator internal service err",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", bytes.NewReader([]byte("11")))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				userClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generated.CheckCreatorMessage{
					ID:        uuid.New().String(),
					IsCreator: false,
					Error:     "",
				}, errors.New("test"))
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r

			},
		},
		{
			name:             "Check if Creator internal err",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", bytes.NewReader([]byte("11")))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				userClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generated.CheckCreatorMessage{
					ID:        uuid.New().String(),
					IsCreator: false,
					Error:     "test",
				}, nil)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r

			},
		},
		{
			name:             "Check if Creator internal err",
			expectedResponse: http.StatusConflict,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", bytes.NewReader([]byte("11")))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				userClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generated.CheckCreatorMessage{
					ID:        uuid.New().String(),
					IsCreator: true,
					Error:     "",
				}, nil)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
		},
		{
			name:             "OK",
			expectedResponse: http.StatusOK,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", bytes.NewReader(bodyPrepare(models.BecameCreatorInfo{
					Name:        "test name",
					Description: "some test",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				userClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generated.CheckCreatorMessage{
					ID:        uuid.New().String(),
					IsCreator: false,
					Error:     "",
				}, nil)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				userClient.EXPECT().BecomeCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{Value: uuid.New().String(), Error: ""}, nil)
				return r
			},
		},
		{
			name:             "err from user service",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", bytes.NewReader(bodyPrepare(models.BecameCreatorInfo{
					Name:        "test name",
					Description: "some test",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				userClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generated.CheckCreatorMessage{
					ID:        uuid.New().String(),
					IsCreator: false,
					Error:     "",
				}, nil)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				userClient.EXPECT().BecomeCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{Value: uuid.New().String(), Error: ""}, errors.New("test"))
				return r
			},
		},
		{
			name:             "err from becomeCreator",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", bytes.NewReader(bodyPrepare(models.BecameCreatorInfo{
					Name:        "test name",
					Description: "some test",
				})))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				userClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generated.CheckCreatorMessage{
					ID:        uuid.New().String(),
					IsCreator: false,
					Error:     "",
				}, nil)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				userClient.EXPECT().BecomeCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{Value: uuid.New().String(), Error: "test"}, nil)
				return r
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()

			h.BecomeCreator(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func TestUserHandler_DeleteProfilePhoto(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	id := uuid.New()
	os.Setenv("CSRF_SECRET", "TEST")
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name             string
		expectedResponse int
		mock             func() *http.Request
	}{
		{
			name:             "Err from CheckUserVersion",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", nil)
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))
				return r

			},
		},
		{
			name:             "Err from CheckUserVersion",
			expectedResponse: http.StatusForbidden,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", nil)
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "test",
				}, nil)
				return r

			},
		},
		{
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/becomeCreator", nil)

				setJWTToken(r, token)
				os.Unsetenv("CSRF_SECRET")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/becomeCreator", nil)

				setJWTToken(r, token)
				os.Setenv("CSRF_SECRET", "TEST")

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedResponse: http.StatusOK,
		},
		{
			name: "Wrong Token",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator",
					nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    "1",
					Expires:  time.Time{},
					HttpOnly: true,
				})

				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name:             "Wrong data",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", bytes.NewReader([]byte("11")))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r

			},
		},
		{
			name:             "Wrong data2",
			expectedResponse: http.StatusBadRequest,
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator", bytes.NewReader([]byte("11")))
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				r = mux.SetURLVars(r, map[string]string{
					"image-uuid": "11",
				})
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()

			h.DeleteProfilePhoto(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func TestUserHandler_UserFollows(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	id := uuid.New()
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name             string
		expectedResponse int
		mock             func() *http.Request
	}{
		{
			name: "Wrong Token",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/becomeCreator",
					nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    "1",
					Expires:  time.Time{},
					HttpOnly: true,
				})

				return r
			},
			expectedResponse: http.StatusUnauthorized,
		},
		{
			name:             "OK",
			expectedResponse: http.StatusOK,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userFollows", nil)
				setJWTToken(r, token)
				userClient.EXPECT().UserFollows(gomock.Any(), gomock.Any()).Return(&generated.FollowsMessage{Follows: follow, Error: ""}, nil)
				return r
			},
		},
		{
			name:             "Err from userFollows",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userFollows", nil)
				setJWTToken(r, token)
				userClient.EXPECT().UserFollows(gomock.Any(), gomock.Any()).Return(&generated.FollowsMessage{Follows: follow, Error: ""}, errors.New("test"))
				return r
			},
		},
		{
			name:             "Err from userFollows2",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userFollows", nil)
				setJWTToken(r, token)
				userClient.EXPECT().UserFollows(gomock.Any(), gomock.Any()).Return(&generated.FollowsMessage{Follows: follow, Error: models.InternalError.Error()}, nil)
				return r
			},
		},
		{
			name:             "Wrong creator uuid",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userFollows", nil)
				setJWTToken(r, token)
				userClient.EXPECT().UserFollows(gomock.Any(), gomock.Any()).Return(&generated.FollowsMessage{Follows: followWithErr1, Error: ""}, nil)
				return r
			},
		},
		{
			name:             "Wrong creatorPhoto uuid",
			expectedResponse: http.StatusInternalServerError,
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/userFollows", nil)
				setJWTToken(r, token)
				userClient.EXPECT().UserFollows(gomock.Any(), gomock.Any()).Return(&generated.FollowsMessage{Follows: followWithErr2, Error: ""}, nil)
				return r
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()

			h.UserFollows(w, r)
			require.Equal(t, test.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.expectedResponse, w.Code))
		})
	}
}

func TestUserHandler_UpdateProfilePhoto(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	id := uuid.New()
	os.Setenv("CSRF_SECRET", "TEST")
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte(uuid.Nil.String()))
				if err != nil {
					t.Error(err)
				}

				// We create the form data field 'upload'
				// which returns another writer to write the actual file
				part, err := writer.CreateFormFile("upload", "img.png")
				if err != nil {
					t.Error(err)
				}

				width := 200
				height := 100
				upLeft := image.Point{0, 0}
				lowRight := image.Point{width, height}

				img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

				// Colors are defined by Red, Green, Blue, Alpha uint8 values.
				cyan := color.RGBA{100, 200, 200, 0xff}

				// Set color for each pixel.
				for x := 0; x < width; x++ {
					for y := 0; y < height; y++ {
						switch {
						case x < width/2 && y < height/2: // upper left quadrant
							img.Set(x, y, cyan)
						case x >= width/2 && y >= height/2: // lower right quadrant
							img.Set(x, y, color.White)
						default:
							// Use zero value.
						}
					}
				}

				// Encode() takes an io.Writer. We pass the multipart field 'upload' that we defined
				// earlier which, in turn, writes to our io.Pipe
				err = png.Encode(part, img)
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					body)

				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				userClient.EXPECT().UpdatePhoto(gomock.Any(), gomock.Any()).Return(&generated.ImageID{
					Value: uuid.New().String(),
					Error: "test",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, "111")

				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.UpdateProfilePhoto(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestUserHandler_UpdateProfilePhoto2(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	id := uuid.New()
	os.Setenv("CSRF_SECRET", "TEST")
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})
	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	userClient := mock.NewMockUserServiceClient(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{

		{
			name: "err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "err from check user version",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, token)
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
				r := httptest.NewRequest("GET", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "no multipart",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Unauthorized error while creating CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, token)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				os.Unsetenv("CSRF_SECRET")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "No upload file",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)
				go func() {
					defer writer.Close()

					partPath, _ := writer.CreateFormField("path")
					_, err := partPath.Write([]byte(uuid.Nil.String()))
					if err != nil {
						t.Error(err)
					}

				}()

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					body)
				os.Setenv("CSRF_SECRET", "TEST")
				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "wrong uuid",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte("11"))
				if err != nil {
					t.Error(err)
				}

				// We create the form data field 'upload'
				// which returns another writer to write the actual file
				part, err := writer.CreateFormFile("upload", "img.png")
				if err != nil {
					t.Error(err)
				}

				width := 200
				height := 100
				upLeft := image.Point{0, 0}
				lowRight := image.Point{width, height}

				img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

				// Colors are defined by Red, Green, Blue, Alpha uint8 values.
				cyan := color.RGBA{100, 200, 200, 0xff}

				// Set color for each pixel.
				for x := 0; x < width; x++ {
					for y := 0; y < height; y++ {
						switch {
						case x < width/2 && y < height/2: // upper left quadrant
							img.Set(x, y, cyan)
						case x >= width/2 && y >= height/2: // lower right quadrant
							img.Set(x, y, color.White)
						default:
							// Use zero value.
						}
					}
				}

				// Encode() takes an io.Writer. We pass the multipart field 'upload' that we defined
				// earlier which, in turn, writes to our io.Pipe
				err = png.Encode(part, img)
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					body)

				setJWTToken(r, token)
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
			name: "err from user service",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte(uuid.Nil.String()))
				if err != nil {
					t.Error(err)
				}

				// We create the form data field 'upload'
				// which returns another writer to write the actual file
				part, err := writer.CreateFormFile("upload", "img.png")
				if err != nil {
					t.Error(err)
				}

				width := 200
				height := 100
				upLeft := image.Point{0, 0}
				lowRight := image.Point{width, height}

				img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

				// Colors are defined by Red, Green, Blue, Alpha uint8 values.
				cyan := color.RGBA{100, 200, 200, 0xff}

				// Set color for each pixel.
				for x := 0; x < width; x++ {
					for y := 0; y < height; y++ {
						switch {
						case x < width/2 && y < height/2: // upper left quadrant
							img.Set(x, y, cyan)
						case x >= width/2 && y >= height/2: // lower right quadrant
							img.Set(x, y, color.White)
						default:
							// Use zero value.
						}
					}
				}

				// Encode() takes an io.Writer. We pass the multipart field 'upload' that we defined
				// earlier which, in turn, writes to our io.Pipe
				err = png.Encode(part, img)
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					body)

				setJWTToken(r, token)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				userClient.EXPECT().UpdatePhoto(gomock.Any(), gomock.Any()).Return(&generated.ImageID{
					Value: uuid.New().String(),
					Error: "",
				}, errors.New("test"))
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				userClient: userClient,
				authClient: authClient,
				logger:     zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.UpdateProfilePhoto(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}
