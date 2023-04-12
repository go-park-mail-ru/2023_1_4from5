package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	mockAttach "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	mockAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	mockPost "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post/mocks"
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

type fields struct {
	usecase           post.PostUsecase
	authUsecase       auth.AuthUsecase
	attachmentUsecase attachment.AttachmentUsecase
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

	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)
	mockAttachUsecase := mockAttach.NewMockAttachmentUsecase(ctl)

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

	testHandler := NewPostHandler(mockPostUsecase, mockAuthUsecase, mockAttachUsecase, zapSugar)
	if testHandler.usecase != mockPostUsecase {
		t.Error("bad constructor")
	}
	if testHandler.authUsecase != mockAuthUsecase {
		t.Error("bad constructor")
	}
	if testHandler.attachmentUsecase != mockAttachUsecase {
		t.Error("bad constructor")
	}
}

type args struct {
	r                *http.Request
	expectedResponse int
}

func TestPostHandler_AddLike(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)
	mockAttachUsecase := mockAttach.NewMockAttachmentUsecase(ctl)
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
		name   string
		fields fields
		args   args
	}{
		{
			name:   "OK",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()}))),
				expectedResponse: http.StatusOK,
			},
		},
		{
			name:   "Unauthorized",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()}))),
				expectedResponse: http.StatusUnauthorized,
			},
		},
		{
			name:   "BadRequest1",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader([]byte("Trying to signIn"))),
				expectedResponse: http.StatusBadRequest,
			},
		},
		{
			name:   "BadRequest2",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{}))),
				expectedResponse: http.StatusBadRequest,
			},
		},
		{
			name:   "BadRequest3",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{}))),
				expectedResponse: http.StatusBadRequest,
			},
		},
		{
			name:   "Forbidden",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{}))),
				expectedResponse: http.StatusForbidden,
			},
		},
		{
			name:   "InternalServerError",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{}))),
				expectedResponse: http.StatusInternalServerError,
			},
		},
		{
			name:   "InternalServerError2",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/addLike",
					bytes.NewReader(bodyPrepare(models.Like{}))),
				expectedResponse: http.StatusInternalServerError,
			},
		},
	}

	for i := 0; i < len(tests); i++ {
		value := token
		switch tests[i].name {
		case "OK":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			mockPostUsecase.EXPECT().AddLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, nil)
		case "Unauthorized":
			value = "body"
		case "InternalServerError":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			mockPostUsecase.EXPECT().AddLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, models.InternalError)
		case "InternalServerError2":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, models.InternalError)
		case "BadRequest1":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
		case "BadRequest2":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			mockPostUsecase.EXPECT().AddLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, models.WrongData)
		case "BadRequest3":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
		case "Forbidden":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, errors.New("test err"))
		}
		tests[i].args.r.AddCookie(&http.Cookie{
			Name:     "SSID",
			Value:    value,
			Expires:  time.Time{},
			HttpOnly: true,
		})
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &PostHandler{
				usecase:           test.fields.usecase,
				authUsecase:       test.fields.authUsecase,
				attachmentUsecase: test.fields.attachmentUsecase,
				logger:            zapSugar,
			}
			w := httptest.NewRecorder()

			h.AddLike(w, test.args.r)
			require.Equal(t, test.args.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d",
				test.name, test.args.expectedResponse, w.Code))
		})
	}
}

func TestPostHandler_RemoveLike(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)
	mockAttachUsecase := mockAttach.NewMockAttachmentUsecase(ctl)

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
		name   string
		fields fields
		args   args
	}{
		{
			name:   "OK",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()}))),
				expectedResponse: http.StatusOK,
			},
		},
		{
			name:   "Unauthorized",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{LikesCount: 0, PostID: uuid.New()}))),
				expectedResponse: http.StatusUnauthorized,
			},
		},
		{
			name:   "BadRequest1",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader([]byte("Trying to signIn"))),
				expectedResponse: http.StatusBadRequest,
			},
		},
		{
			name:   "BadRequest2",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{}))),
				expectedResponse: http.StatusBadRequest,
			},
		},
		{
			name:   "Forbidden",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{}))),
				expectedResponse: http.StatusForbidden,
			},
		},
		{
			name:   "InternalServerError",
			fields: fields{mockPostUsecase, mockAuthUsecase, mockAttachUsecase},
			args: args{
				r: httptest.NewRequest("PUT", "/api/post/removeLike",
					bytes.NewReader(bodyPrepare(models.Like{}))),
				expectedResponse: http.StatusInternalServerError,
			},
		},
	}

	for i := 0; i < len(tests); i++ {
		value := token
		switch tests[i].name {
		case "OK":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().RemoveLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, nil)
		case "Unauthorized":
			value = "body"
		case "InternalServerError":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().RemoveLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, models.InternalError)
		case "BadRequest1":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
		case "BadRequest2":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().RemoveLike(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Like{}, models.WrongData)
		case "Forbidden":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, errors.New("test err"))
		}
		tests[i].args.r.AddCookie(&http.Cookie{
			Name:     "SSID",
			Value:    value,
			Expires:  time.Time{},
			HttpOnly: true,
		})
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &PostHandler{
				usecase:           test.fields.usecase,
				authUsecase:       test.fields.authUsecase,
				attachmentUsecase: test.fields.attachmentUsecase,
				logger:            zapSugar,
			}
			w := httptest.NewRecorder()

			h.RemoveLike(w, test.args.r)
			require.Equal(t, test.args.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d",
				test.name, test.args.expectedResponse, w.Code))
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

	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)
	mockAttachUsecase := mockAttach.NewMockAttachmentUsecase(ctl)

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
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, errors.New("test"))

				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/post/create",
					nil)

				setJWTToken(r, bdy)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)

				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/create",
					nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, "111")

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)

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

				os.Setenv("CSRF_SECRET", "TEST")
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, models.WrongData)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Is creator internal",
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, models.InternalError)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Is creator false",
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsCreator(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
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
			h.CreatePost(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

var newPostErr = models.PostEditData{Title: "testavbjkwkjebojweabkvsn;awlvmnbjerkvjawlvnkaoeibr aelsvjoerbjvkas,zjfonwileabuv", Text: "testtest"}
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

	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)
	mockAttachUsecase := mockAttach.NewMockAttachmentUsecase(ctl)

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
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "//post/edit/{post-uuid}",
					nil)

				setJWTToken(r, bdy)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, errors.New("test"))

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
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)

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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)

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
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Is post owner wrong data",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/edit/{post-uuid}",
					bytes.NewReader(bodyPrepare(newPost)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, models.WrongData)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Is post owner internal error",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/edit/{post-uuid}",
					bytes.NewReader(bodyPrepare(newPost)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, models.InternalError)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Is post owner forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/edit/{post-uuid}",
					bytes.NewReader(bodyPrepare(newPost)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Wrong data in body",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/edit/{post-uuid}",
					bytes.NewReader([]byte("111")))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Wrong length for text or title",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/edit/{post-uuid}",
					bytes.NewReader(bodyPrepare(newPostErr)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Wrong length for text or title",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/edit/{post-uuid}",
					bytes.NewReader(bodyPrepare(newPost)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockPostUsecase.EXPECT().EditPost(gomock.Any(), gomock.Any()).Return(models.InternalError)

				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Wrong length for text or title",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/post/edit/{post-uuid}",
					bytes.NewReader(bodyPrepare(newPost)))

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				r = mux.SetURLVars(r, map[string]string{
					"post-uuid": uuid.NewString(),
				})

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockPostUsecase.EXPECT().IsPostOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockPostUsecase.EXPECT().EditPost(gomock.Any(), gomock.Any()).Return(nil)

				return r
			},
			expectedStatus: http.StatusOK,
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
			h.EditPost(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

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
}
