package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	client generatedAuth.AuthServiceClient
	logger *zap.SugaredLogger
}

func NewAuthHandler(cl generatedAuth.AuthServiceClient, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		client: cl,
		logger: logger,
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	user := models.LoginUser{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err = (models.User{Login: user.Login, PasswordHash: user.PasswordHash}).UserAuthIsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.client.SignIn(r.Context(), &generatedAuth.LoginUser{
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if token.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusUnauthorized, "no user with such login")
		return
	}
	if len(token.Error) != 0 {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	utils.Cookie(w, token.Cookie, "SSID")
	utils.Response(w, http.StatusOK, nil)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userData, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if uv, err := h.client.IncUserVersion(r.Context(), &generatedAuth.AccessDetails{
		Login:       userData.Login,
		Id:          userData.Id.String(),
		UserVersion: userData.UserVersion,
	}); err != nil || len(uv.Error) != 0 {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Cookie(w, "", "SSID")
	utils.Response(w, http.StatusOK, nil)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err = user.UserIsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.client.SignUp(r.Context(), &generatedAuth.User{
		Id:           user.Id.String(),
		Login:        user.Login,
		Name:         user.Name,
		ProfilePhoto: user.ProfilePhoto.String(),
		PasswordHash: user.PasswordHash,
		Registration: user.Registration.String(),
		UserVersion:  user.UserVersion,
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if len(token.Cookie) == 0 {
		if token.Error == models.WrongData.Error() {
			utils.Response(w, http.StatusConflict, "user with such login already exists")
			return
		}
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Cookie(w, token.Cookie, "SSID")
	utils.Response(w, http.StatusOK, nil)
}
