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
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/notifications/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
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

var subs = []*generatedCommon.Subscription{&generatedCommon.Subscription{
	Id:           uuid.New().String(),
	Creator:      uuid.New().String(),
	CreatorName:  "test",
	CreatorPhoto: uuid.New().String(),
	MonthCost:    0,
	Title:        "test",
	Description:  "test",
}}
var attachs = []*generated.Attachment{&generated.Attachment{
	ID:   uuid.New().String(),
	Type: "test",
}}

var post = &generated.Post{
	Id:              uuid.New().String(),
	CreatorID:       uuid.New().String(),
	Creation:        "2006-01-02 15:04:05 -0700 -0700",
	LikesCount:      2,
	CreatorPhoto:    uuid.New().String(),
	CreatorName:     "test",
	Title:           "test",
	Text:            "",
	IsAvailable:     false,
	IsLiked:         false,
	PostAttachments: attachs,
	Subscriptions:   subs,
}

var postWithErr = *post

var posts = []*generated.Post{post}

var postsMes = &generated.PostsMessage{
	Posts: posts,
	Error: "",
}
var postsMesWithErr = &generated.PostsMessage{
	Posts: posts,
	Error: "11",
}
var testUser = models.User{
	Login:        "Dasha2003!",
	PasswordHash: "Dasha2003!",
	Name:         "Дарья Такташова",
}

var aim = &generated.Aim{
	Creator:     uuid.New().String(),
	Description: "test",
	MoneyNeeded: 100,
	MoneyGot:    10,
}

var aimWithErr = &generated.Aim{
	Creator:     "111",
	Description: "test",
	MoneyNeeded: 100,
	MoneyGot:    10,
}

var creator = &generated.Creator{
	Id:             uuid.New().String(),
	UserID:         uuid.New().String(),
	CreatorName:    "test",
	CreatorPhoto:   uuid.New().String(),
	CoverPhoto:     uuid.New().String(),
	FollowersCount: 0,
	Description:    "test",
	PostsCount:     0,
}

var creatorWithErr = &generated.Creator{
	Id:             "11",
	UserID:         uuid.New().String(),
	CreatorName:    "test",
	CreatorPhoto:   uuid.New().String(),
	CoverPhoto:     uuid.New().String(),
	FollowersCount: 0,
	Description:    "test",
	PostsCount:     0,
}

var creators = &generated.CreatorsMessage{
	Creators: []*generated.Creator{creator},
	Error:    "",
}

var creatorsWithErr = &generated.CreatorsMessage{
	Creators: []*generated.Creator{creator},
	Error:    "err",
}

var creatorsWithErr2 = &generated.CreatorsMessage{
	Creators: []*generated.Creator{creatorWithErr},
	Error:    "",
}

func TestNewCreatorHandler(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	authClient := mockAuth.NewMockAuthServiceClient(ctl)
	creatorClient := mockCreator.NewMockCreatorServiceClient(ctl)
	notify := mock.NewMockNotificationApp(ctl)

	logger := zap.NewNop()

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()

	testHandler := NewCreatorHandler(creatorClient, authClient, notify, zapSugar)
	if testHandler.authClient != authClient || testHandler.creatorClient != creatorClient || testHandler.notificationApp != notify {
		t.Error("bad constructor")
	}
}

type args struct {
	r                *http.Request
	expectedResponse http.Response
}

func TestCreatorHandler_GetPage(t *testing.T) {
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
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

	tests := []struct {
		name string
		args args
		mock func(r *http.Request) *http.Request
	}{
		{
			name: "OK",
			args: args{
				r:                httptest.NewRequest("GET", "/creator/page", strings.NewReader(fmt.Sprint())),
				expectedResponse: http.Response{StatusCode: http.StatusOK},
			},
			mock: func(r *http.Request) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.NewString(),
				})
				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				creatorClient.EXPECT().
					GetPage(gomock.Any(), gomock.Any()).
					Return(&generated.CreatorPage{
						Error: "", CreatorInfo: creator, AimInfo: aim, Posts: posts, Subscriptions: subs,
					}, nil)
				return r
			},
		},
		{
			name: "Wrong Id",
			args: args{
				r:                httptest.NewRequest("GET", "/creator/page", strings.NewReader(fmt.Sprint())),
				expectedResponse: http.Response{StatusCode: http.StatusBadRequest},
			},
			mock: func(r *http.Request) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": "uuid.NewString()",
				})
				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				return r
			},
		},
		{
			name: "Wrong Id",
			args: args{
				r:                httptest.NewRequest("GET", "/creator/page", strings.NewReader(fmt.Sprint())),
				expectedResponse: http.Response{StatusCode: http.StatusBadRequest},
			},
			mock: func(r *http.Request) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"creator-uud": "uuid.NewString()",
				})
				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				return r
			},
		},
		{
			name: "InternalErr from service",
			args: args{
				r:                httptest.NewRequest("GET", "/creator/page", strings.NewReader(fmt.Sprint())),
				expectedResponse: http.Response{StatusCode: http.StatusInternalServerError},
			},
			mock: func(r *http.Request) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.NewString(),
				})
				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				creatorClient.EXPECT().
					GetPage(gomock.Any(), gomock.Any()).
					Return(&generated.CreatorPage{
						Error: "", CreatorInfo: creator, AimInfo: aim, Posts: posts, Subscriptions: subs,
					}, errors.New("err"))
				return r
			},
		},
		{
			name: "InternalErr",
			args: args{
				r:                httptest.NewRequest("GET", "/creator/page", strings.NewReader(fmt.Sprint())),
				expectedResponse: http.Response{StatusCode: http.StatusInternalServerError},
			},
			mock: func(r *http.Request) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.NewString(),
				})
				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				creatorClient.EXPECT().
					GetPage(gomock.Any(), gomock.Any()).
					Return(&generated.CreatorPage{
						Error: models.InternalError.Error(), CreatorInfo: creator, AimInfo: aim, Posts: posts, Subscriptions: subs,
					}, nil)
				return r
			},
		},
		{
			name: "WrongData",
			args: args{
				r:                httptest.NewRequest("GET", "/creator/page", strings.NewReader(fmt.Sprint())),
				expectedResponse: http.Response{StatusCode: http.StatusBadRequest},
			},
			mock: func(r *http.Request) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.NewString(),
				})
				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				creatorClient.EXPECT().
					GetPage(gomock.Any(), gomock.Any()).
					Return(&generated.CreatorPage{
						Error: models.WrongData.Error(), CreatorInfo: creator, AimInfo: aim, Posts: posts, Subscriptions: subs,
					}, nil)
				return r
			},
		},
		{
			name: "Err while converting",
			args: args{
				r:                httptest.NewRequest("GET", "/creator/page", strings.NewReader(fmt.Sprint())),
				expectedResponse: http.Response{StatusCode: http.StatusInternalServerError},
			},
			mock: func(r *http.Request) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"creator-uuid": uuid.NewString(),
				})
				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				creatorClient.EXPECT().
					GetPage(gomock.Any(), gomock.Any()).
					Return(&generated.CreatorPage{
						Error: "", CreatorInfo: creator, AimInfo: aimWithErr, Posts: posts, Subscriptions: subs,
					}, nil)
				return r
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			test.args.r = test.mock(test.args.r)

			h.GetPage(w, test.args.r)
			require.Equal(t, test.args.expectedResponse.StatusCode, w.Code, fmt.Errorf("%s :  expected %d, got %d,",
				test.name, test.args.expectedResponse.StatusCode, w.Code))
		})
	}

}

var testAim = models.Aim{Creator: uuid.New(), Description: "test", MoneyNeeded: 500, MoneyGot: 0}
var testAimWithLongDescription = models.Aim{Creator: uuid.New(), Description: "testtesttesttesttesttesttesttesttesttesttesttesttesttesttest" +
	"testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest" +
	"testtesttesttesttesttesttest", MoneyNeeded: 500, MoneyGot: 0}

func bodyPrepare(val interface{}) []byte {
	valJSON, err := json.Marshal(&val)
	if err != nil {
		return nil
	}
	return valJSON
}

func TestCreatorHandler_CreateAim(t *testing.T) {
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
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

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

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().CreateAim(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)
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
			name: "Err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
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

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "err",
				}, nil)
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

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
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

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal Error from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, errors.New("err"))
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Error from is creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: models.InternalError.Error(),
				}, nil)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Wrong data from is creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: models.WrongData.Error(),
				}, nil)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Not a creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  false,
					Error: "",
				}, nil)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Internal err from creator service Create Aim",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().CreateAim(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal err from creator service Create Aim",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/aim/create",
					bytes.NewReader(bodyPrepare(testAim)))

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().IsCreator(gomock.Any(), gomock.Any()).Return(&generated.FlagMessage{
					Flag:  true,
					Error: "",
				}, nil)
				creatorClient.EXPECT().CreateAim(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: models.InternalError.Error()}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
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

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/creator/list", nil)

				creatorClient.EXPECT().GetAllCreators(gomock.Any(), gomock.Any()).Return(creators, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Internal err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/creator/list", nil)

				creatorClient.EXPECT().GetAllCreators(gomock.Any(), gomock.Any()).Return(creators, errors.New("err"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal err from GetAllCreators",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/creator/list", nil)

				creatorClient.EXPECT().GetAllCreators(gomock.Any(), gomock.Any()).Return(creatorsWithErr, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal err from converting",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/creator/list", nil)

				creatorClient.EXPECT().GetAllCreators(gomock.Any(), gomock.Any()).Return(creatorsWithErr2, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.GetAllCreators(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestCreatorHandler_FindCreator(t *testing.T) {
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

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/find/", nil)
				r = mux.SetURLVars(r, map[string]string{
					"keyword": "test author",
				})
				creatorClient.EXPECT().FindCreators(gomock.Any(), gomock.Any()).Return(creators, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "no keyword",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/find/", nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal err from find creators",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/find/", nil)
				r = mux.SetURLVars(r, map[string]string{
					"keyword": "test author",
				})
				creatorClient.EXPECT().FindCreators(gomock.Any(), gomock.Any()).Return(creatorsWithErr, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "internal err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/find/", nil)
				r = mux.SetURLVars(r, map[string]string{
					"keyword": "test author",
				})
				creatorClient.EXPECT().FindCreators(gomock.Any(), gomock.Any()).Return(creatorsWithErr, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "wrong creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/find/", nil)
				r = mux.SetURLVars(r, map[string]string{
					"keyword": "test author",
				})
				creatorClient.EXPECT().FindCreators(gomock.Any(), gomock.Any()).Return(creatorsWithErr2, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.FindCreator(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestCreatorHandler_GetFeed(t *testing.T) {
	postWithErr.Id = "11"
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
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/feed", nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(postsMes, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Err while converting",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/feed", nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(&generated.PostsMessage{
					Posts: []*generated.Post{&postWithErr},
					Error: "",
				}, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/feed", nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(postsMes, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Err from Get Feed",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/feed", nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(postsMesWithErr, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/feed", nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Err from Check User Version",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/feed", nil)

				r.AddCookie(&http.Cookie{
					Name:     "SSID",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "1",
				}, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/feed", nil)

				r.AddCookie(&http.Cookie{
					Name:     "S",
					Value:    bdy,
					Expires:  time.Time{},
					HttpOnly: true,
				})
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.GetFeed(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
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

func TestCreatorHandler_UpdateProfilePhoto(t *testing.T) {
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
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateProfilePhoto", nil)

				setJWTToken(r, "111")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "BadRequest",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateProfilePhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateProfilePhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, errors.New("test"))

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Not Found err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateProfilePhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.NotFound.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateProfilePhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.InternalError.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal err no multipart",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateProfilePhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.New().String(),
					Error: "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateProfilePhoto", nil)

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
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/updateProfilePhoto", nil)

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
				r := httptest.NewRequest("GET", "/updateProfilePhoto", nil)

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
				r := httptest.NewRequest("PUT", "/updateProfilePhoto", nil)

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
			name: "Internal Err from Check User Version",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/feed", nil)

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
			name: "No photo",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte("111"))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("PUT", "/updateProfilePhoto",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.New().String(),
					Error: "",
				}, nil)

				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Wrong uuid",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte("111"))
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

				r := httptest.NewRequest("PUT", "/updateProfilePhoto",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.New().String(),
					Error: "",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.UpdateProfilePhoto(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestCreatorHandler_UpdateCoverPhoto(t *testing.T) {
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
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCoverPhoto", nil)

				setJWTToken(r, "111")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "BadRequest",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCoverPhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCoverPhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, errors.New("test"))

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Not Found err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCoverPhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.NotFound.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCoverPhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.InternalError.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCoverPhoto", nil)

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
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/updateCoverPhoto", nil)

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
				r := httptest.NewRequest("GET", "/updateCoverPhoto", nil)

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
				r := httptest.NewRequest("PUT", "/updateCoverPhoto", nil)

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
			name: "Internal Err from Check User Version",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCoverPhoto", nil)

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
			name: "No photo",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte("111"))
				if err != nil {
					t.Error(err)
				}

				r := httptest.NewRequest("POST", "/user/updateCoverPhoto",
					body)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.New().String(),
					Error: "",
				}, nil)

				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal err no multipart",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCoverPhoto", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.New().String(),
					Error: "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Wrong uuid",
			mock: func() *http.Request {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)

				defer writer.Close()

				partPath, _ := writer.CreateFormField("path")
				_, err := partPath.Write([]byte("111"))
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

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.New().String(),
					Error: "",
				}, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.UpdateCoverPhoto(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestCreatorHandler_UpdateCreatorData(t *testing.T) {
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
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", nil)

				setJWTToken(r, "111")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "BadRequest",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, errors.New("test"))

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Not Found err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.NotFound.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.InternalError.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", nil)

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
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/updateCreatorData", nil)

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
				r := httptest.NewRequest("GET", "/updateCreatorData", nil)

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
				r := httptest.NewRequest("PUT", "/updateCreatorData", nil)

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
			name: "Internal Err from Check User Version",
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
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", bytes.NewReader(bodyPrepare(models.BecameCreatorInfo{
					Name:        "testName",
					Description: "some test description",
				})))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)
				creatorClient.EXPECT().UpdateCreatorData(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, nil)

				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", bytes.NewReader(bodyPrepare(models.BecameCreatorInfo{
					Name:        "testName",
					Description: "some test description",
				})))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)
				creatorClient.EXPECT().UpdateCreatorData(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: ""}, errors.New("test"))

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "err from creator update creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/updateCreatorData", bytes.NewReader(bodyPrepare(models.BecameCreatorInfo{
					Name:        "testName",
					Description: "some test description",
				})))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)
				creatorClient.EXPECT().UpdateCreatorData(gomock.Any(), gomock.Any()).Return(&generatedCommon.Empty{Error: "test"}, nil)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.UpdateCreatorData(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestCreatorHandler_DeleteCoverPhoto(t *testing.T) {
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
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteCoverPhoto", nil)

				setJWTToken(r, "111")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "BadRequest",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteCoverPhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteCoverPhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, errors.New("test"))

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Not Found err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteCoverPhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.NotFound.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteCoverPhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.InternalError.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteCoverPhoto", nil)
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
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/DeleteCoverPhoto", nil)
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
				r := httptest.NewRequest("GET", "/DeleteCoverPhoto", nil)
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
				r := httptest.NewRequest("DELETE", "/DeleteCoverPhoto", nil)
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
			name: "Internal Err from Check User Version",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteCoverPhoto", nil)
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
			name: "wrong uuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteCoverPhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)
				r = mux.SetURLVars(r, map[string]string{
					"image-uuid": "1",
				})

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.DeleteCoverPhoto(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func TestCreatorHandler_DeleteProfilePhoto(t *testing.T) {
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
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteProfilePhoto", nil)

				setJWTToken(r, "111")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "BadRequest",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteProfilePhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from creator service",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteProfilePhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, errors.New("test"))

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Not Found err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteProfilePhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.NotFound.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal err from Check If Creator",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteProfilePhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: models.InternalError.Error(),
				}, nil)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Internal Err from auth service",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteProfilePhoto", nil)
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
			name: "Get CSRF with error",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/DeleteProfilePhoto", nil)
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
				r := httptest.NewRequest("GET", "/DeleteProfilePhoto", nil)
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
				r := httptest.NewRequest("DELETE", "/DeleteProfilePhoto", nil)
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
			name: "Internal Err from Check User Version",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteProfilePhoto", nil)
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
			name: "wrong uuid",
			mock: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/DeleteProfilePhoto", nil)
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				authClient.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(&generatedAuth.UserVersion{
					UserVersion: int64(1),
					Error:       "",
				}, nil)
				creatorClient.EXPECT().CheckIfCreator(gomock.Any(), gomock.Any()).Return(&generatedCommon.UUIDResponse{
					Value: uuid.Nil.String(),
					Error: "",
				}, nil)
				r = mux.SetURLVars(r, map[string]string{
					"image-uuid": "1",
				})

				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &CreatorHandler{
				creatorClient: creatorClient,
				authClient:    authClient,
				logger:        zapSugar,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.DeleteProfilePhoto(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}
