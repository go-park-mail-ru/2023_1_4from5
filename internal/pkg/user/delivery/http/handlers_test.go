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
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	mock "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/mocks"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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
	const body = "body"
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

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
			usecaseMock.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(models.UserProfile{}, nil)
			status = http.StatusOK
		case 1:
			value = body
			status = http.StatusUnauthorized
		case 2:
			usecaseMock.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(models.UserProfile{}, models.InternalError)
			status = http.StatusInternalServerError
		case 3:
			usecaseMock.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(models.UserProfile{}, models.NotFound)
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
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

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
			usecaseMock.EXPECT().GetHomePage(gomock.Any(), gomock.Any()).Return(models.UserHomePage{}, nil)
			status = http.StatusOK
		case 1:
			value = "body"
			status = http.StatusUnauthorized
		case 2:
			usecaseMock.EXPECT().GetHomePage(gomock.Any(), gomock.Any()).Return(models.UserHomePage{}, models.InternalError)
			status = http.StatusInternalServerError
		case 3:
			usecaseMock.EXPECT().GetHomePage(gomock.Any(), gomock.Any()).Return(models.UserHomePage{}, models.NotFound)
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

func TestUpdatePassword(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("SECRET", "TEST")
	os.Setenv("CSRF_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: uuid.New()})

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
			usecaseMock.EXPECT().GetHomePage(gomock.Any(), gomock.Any()).Return(models.UserHomePage{}, nil)
			status = http.StatusOK
		case 1:
			value = "body"
			status = http.StatusUnauthorized
		case 2:
			usecaseMock.EXPECT().GetHomePage(gomock.Any(), gomock.Any()).Return(models.UserHomePage{}, models.InternalError)
			status = http.StatusInternalServerError
		case 3:
			usecaseMock.EXPECT().GetHomePage(gomock.Any(), gomock.Any()).Return(models.UserHomePage{}, models.NotFound)
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

func TestUserHandler_UpdateProfilePhoto(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	os.Setenv("TOKEN_SECRET", "TEST")
	os.Setenv("CSRF_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	id := uuid.New()
	name := uuid.New()
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})

	usecaseMock := mock.NewMockUserUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)

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

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				usecaseMock.EXPECT().UpdatePhoto(gomock.Any(), gomock.Any()).Return(name, models.InternalError)
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

				setJWTToken(r, bdy)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, errors.New("test"))

				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, bdy)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)

				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Unauthorized error while creating CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, bdy)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				os.Unsetenv("CSRF_SECRET")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
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
				r := httptest.NewRequest("POST", "/user/updateProfilePhoto",
					nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				os.Setenv("CSRF_SECRET", "TEST")
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
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

				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "OK",
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

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				r.Header.Add("Content-Type", writer.FormDataContentType())
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &UserHandler{
				usecase:     usecaseMock,
				authUsecase: mockAuthUsecase,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.UpdateProfilePhoto(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

func bodyPrepare(body interface{}) []byte {
	bodyJSON, err := json.Marshal(&body)
	if err != nil {
		return nil
	}
	return bodyJSON
}

func responseBodyPrepare(status, moneyGot int) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	utils.Response(w, status, moneyGot)
	return w
}

func responseBodyPrepareErr(status int, moneyGot interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	utils.Response(w, status, moneyGot)
	return w
}

var testDonation = models.Donate{CreatorID: uuid.New(), MoneyCount: 100}
var testDonationErr = models.Donate{CreatorID: uuid.New(), MoneyCount: -1}

func TestUserHandler_Donate(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("TOKEN_SECRET", "TEST")
	os.Setenv("CSRF_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	id := uuid.New()
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})

	usecaseMock := mock.NewMockUserUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
		expectedBody   int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/donate", bytes.NewReader(bodyPrepare(testDonation)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				usecaseMock.EXPECT().Donate(gomock.Any(), gomock.Any(), gomock.Any()).Return(100, nil)
				return r
			},
			expectedStatus: http.StatusOK,
			expectedBody:   testDonation.MoneyCount,
		},
		{
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/donate", bytes.NewReader(bodyPrepare(testDonation)))
				setJWTToken(r, "1111")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/donate", bytes.NewReader(bodyPrepare(testDonation)))
				setJWTToken(r, bdy)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/user/donate", bytes.NewReader(nil))
				setJWTToken(r, bdy)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Get CSRF but no secret key",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/user/donate", bytes.NewReader(nil))
				setJWTToken(r, bdy)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				os.Unsetenv("CSRF_SECRET")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/donate",
					nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, "111")

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)

				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "BadRequest",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/donate", bytes.NewReader(bodyPrepare(testDonationErr)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   testDonation.MoneyCount,
		},
		{
			name: "Internal Error",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/donate", bytes.NewReader(bodyPrepare(testDonation)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				usecaseMock.EXPECT().Donate(gomock.Any(), gomock.Any(), gomock.Any()).Return(100, models.InternalError)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   testDonation.MoneyCount,
		},
		{
			name: "Bad Request2",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/donate", bytes.NewReader(bodyPrepare(testDonation)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				usecaseMock.EXPECT().Donate(gomock.Any(), gomock.Any(), gomock.Any()).Return(100, models.WrongData)
				return r
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   testDonation.MoneyCount,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv("CSRF_SECRET", "TEST")
			h := &UserHandler{
				usecase:     usecaseMock,
				authUsecase: mockAuthUsecase,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.Donate(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
			if test.expectedStatus != http.StatusOK || test.name == "Get CSRF" {
				require.Equal(t, responseBodyPrepareErr(test.expectedStatus, nil).Body.String(), w.Body.String(), fmt.Errorf("Wrong body"))
			} else {
				require.Equal(t, responseBodyPrepare(test.expectedStatus, test.expectedBody).Body.String(), w.Body.String(), fmt.Errorf("Wrong body"))
			}
		})
	}
}

var testUpdateProfileInfo = models.UpdateProfileInfo{Login: "Dasha2003", Name: "Daria Taktashova"}
var testUpdateProfileInfoErr = models.UpdateProfileInfo{Login: "D", Name: "Daria Taktashova"}

func TestUserHandler_UpdateData(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("TOKEN_SECRET", "TEST")
	os.Setenv("CSRF_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	id := uuid.New()
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})

	usecaseMock := mock.NewMockUserUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
		expectedBody   int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updateData", bytes.NewReader(bodyPrepare(testUpdateProfileInfo)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				usecaseMock.EXPECT().UpdateProfileInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updateData", bytes.NewReader(bodyPrepare(testDonation)))
				setJWTToken(r, "1111")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updateData", bytes.NewReader(bodyPrepare(testDonation)))
				setJWTToken(r, bdy)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/user/updateData", bytes.NewReader(nil))
				setJWTToken(r, bdy)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/updateData",
					nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, "tokenCSRF")

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)

				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Get CSRF but no secret key",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/user/updateData", bytes.NewReader(nil))
				setJWTToken(r, bdy)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				os.Unsetenv("CSRF_SECRET")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Wrong Data",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updateData", bytes.NewReader(bodyPrepare(testUpdateProfileInfoErr)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal Error",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updateData", bytes.NewReader(bodyPrepare(testUpdateProfileInfo)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				usecaseMock.EXPECT().UpdateProfileInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.InternalError)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv("CSRF_SECRET", "TEST")
			h := &UserHandler{
				usecase:     usecaseMock,
				authUsecase: mockAuthUsecase,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.UpdateData(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}

var testUpdatePasswordInfo = models.UpdatePasswordInfo{NewPassword: "12345678ab", OldPassword: "12345678aa"}
var testUpdatePasswordInfoSame = models.UpdatePasswordInfo{NewPassword: "12345678ab", OldPassword: "12345678ab"}
var testUpdatePasswordInfoErr = models.UpdatePasswordInfo{NewPassword: "1234", OldPassword: "12345678ab"}

func TestUserHandler_UpdatePassword(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	os.Setenv("TOKEN_SECRET", "TEST")
	os.Setenv("CSRF_SECRET", "TEST")
	tkn := &usecase.Tokenator{}
	id := uuid.New()
	bdy, _ := tkn.GetJWTToken(context.Background(), models.User{Login: testUser.Login, Id: id})
	tokenCSRF, _ := token.GetCSRFToken(models.User{Login: testUser.Login, Id: id})

	usecaseMock := mock.NewMockUserUsecase(ctl)
	mockAuthUsecase := mockAuth.NewMockAuthUsecase(ctl)

	tests := []struct {
		name           string
		mock           func() *http.Request
		expectedStatus int
	}{
		{
			name: "OK",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updatePassword", bytes.NewReader(bodyPrepare(testUpdatePasswordInfo)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockAuthUsecase.EXPECT().EncryptPwd(gomock.Any(), gomock.Any()).Return("cbkjqvaiuwklaNCVBhjaewskl")
				mockAuthUsecase.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, nil)
				usecaseMock.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockAuthUsecase.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return("token", nil)

				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Get CSRF but no secret key",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/user/updatePassword", bytes.NewReader(nil))
				setJWTToken(r, bdy)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				os.Unsetenv("CSRF_SECRET")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("POST", "/user/updatePassword",
					nil)

				setJWTToken(r, bdy)
				setCSRFToken(r, "tokenCSRF")

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)

				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Unauthorized",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updatePassword", bytes.NewReader(bodyPrepare(testDonation)))
				setJWTToken(r, "bdy")
				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Forbidden",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updatePassword", bytes.NewReader(bodyPrepare(testDonation)))
				setJWTToken(r, bdy)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Get CSRF",
			mock: func() *http.Request {
				r := httptest.NewRequest("GET", "/user/updatePassword", bytes.NewReader(nil))
				setJWTToken(r, bdy)

				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Same Password",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updatePassword", bytes.NewReader(bodyPrepare(testUpdatePasswordInfoSame)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Password is not valid",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updatePassword", bytes.NewReader(bodyPrepare(testUpdatePasswordInfoErr)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Password is not valid",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updatePassword", bytes.NewReader(bodyPrepare(testUpdatePasswordInfo)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockAuthUsecase.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, models.InternalError)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Password is not valid",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updatePassword", bytes.NewReader(bodyPrepare(testUpdatePasswordInfo)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockAuthUsecase.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, errors.New("test"))
				return r
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal Error in UpdatePassword",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updatePassword", bytes.NewReader(bodyPrepare(testUpdatePasswordInfo)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockAuthUsecase.EXPECT().EncryptPwd(gomock.Any(), gomock.Any()).Return("cbkjqvaiuwklaNCVBhjaewskl")
				mockAuthUsecase.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, nil)
				usecaseMock.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.InternalError)
				return r
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Error in SignIn",
			mock: func() *http.Request {
				r := httptest.NewRequest("PUT", "/user/updatePassword", bytes.NewReader(bodyPrepare(testUpdatePasswordInfo)))
				setJWTToken(r, bdy)
				setCSRFToken(r, tokenCSRF)
				mockAuthUsecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Return(0, nil)
				mockAuthUsecase.EXPECT().EncryptPwd(gomock.Any(), gomock.Any()).Return("cbkjqvaiuwklaNCVBhjaewskl")
				mockAuthUsecase.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(models.User{}, nil)
				usecaseMock.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockAuthUsecase.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return("token", errors.New("test"))

				return r
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv("CSRF_SECRET", "TEST")
			h := &UserHandler{
				usecase:     usecaseMock,
				authUsecase: mockAuthUsecase,
			}
			w := httptest.NewRecorder()
			r := test.mock()
			h.UpdatePassword(w, r)
			require.Equal(t, test.expectedStatus, w.Code, fmt.Errorf("%s :  expected %d, got %d,"+
				" for test:%s", test.name, test.expectedStatus, w.Code, test.name))
		})
	}
}
