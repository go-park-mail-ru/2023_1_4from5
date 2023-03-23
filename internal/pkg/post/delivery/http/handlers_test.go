package http

import (
	"bytes"
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
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func bodyPrepare(like models.Like) []byte {
	userjson, err := json.Marshal(&like)
	if err != nil {
		return nil
	}
	return userjson
}

var testUser = models.User{
	Login:        "Dasha2003!",
	PasswordHash: "Dasha2003!",
	Name:         "Дарья Такташова",
}

type fields struct {
	usecase           post.PostUsecase
	authUsecase       auth.AuthUsecase
	attachmentUsecase attachment.AttachmentUsecase
}

func TestNewPostHandler(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)
	mockPostUsecase := mockPost.NewMockPostUsecase(ctl)
	mockAttachUsecase := mockAttach.NewMockAttachmentUsecase(ctl)
	testHandler := NewPostHandler(mockPostUsecase, mockAuthUsecase, mockAttachUsecase)
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

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(models.User{Login: testUser.Login, Id: uuid.New()})

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
	}

	for i := 0; i < len(tests); i++ {
		value := token
		switch tests[i].name {
		case "OK":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().AddLike(gomock.Any(), gomock.Any()).Return(models.Like{}, nil)
		case "Unauthorized":
			value = "body"
		case "InternalServerError":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().AddLike(gomock.Any(), gomock.Any()).Return(models.Like{}, models.InternalError)
		case "BadRequest1":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, nil)
		case "BadRequest2":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().AddLike(gomock.Any(), gomock.Any()).Return(models.Like{}, models.WrongData)
		case "Forbidden":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, errors.New("test err"))
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

	os.Setenv("TOKEN_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	token, _ := tkn.GetJWTToken(models.User{Login: testUser.Login, Id: uuid.New()})

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
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().RemoveLike(gomock.Any(), gomock.Any()).Return(models.Like{}, nil)
		case "Unauthorized":
			value = "body"
		case "InternalServerError":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().RemoveLike(gomock.Any(), gomock.Any()).Return(models.Like{}, models.InternalError)
		case "BadRequest1":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, nil)
		case "BadRequest2":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, nil)
			mockPostUsecase.EXPECT().RemoveLike(gomock.Any(), gomock.Any()).Return(models.Like{}, models.WrongData)
		case "Forbidden":
			mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any()).Return(0, errors.New("test err"))
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
			}
			w := httptest.NewRecorder()

			h.RemoveLike(w, test.args.r)
			require.Equal(t, test.args.expectedResponse, w.Code, fmt.Errorf("%s :  expected %d, got %d",
				test.name, test.args.expectedResponse, w.Code))
		})
	}
}
