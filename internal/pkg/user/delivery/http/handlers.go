package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
	generatedUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type UserHandler struct {
	usecase    user.UserUsecase
	userClient generatedUser.UserServiceClient
	authClient generatedAuth.AuthServiceClient
}

func NewUserHandler(uc user.UserUsecase, userClient generatedUser.UserServiceClient, auc generatedAuth.AuthServiceClient) *UserHandler {
	return &UserHandler{
		usecase:    uc,
		userClient: userClient,
		authClient: auc,
	}
}

func (h *UserHandler) Follow(w http.ResponseWriter, r *http.Request) {
	userInfo, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	creatorUUID, ok := mux.Vars(r)["creator-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creatorId, err := uuid.Parse(creatorUUID)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	answer, err := h.userClient.Follow(r.Context(), &generatedUser.FollowMessage{
		UserID:    userInfo.Id.String(),
		CreatorID: creatorId.String(),
	})
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if answer.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if answer.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *UserHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	userInfo, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	creatorUUID, ok := mux.Vars(r)["creator-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creatorId, err := uuid.Parse(creatorUUID)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	err = h.usecase.Unfollow(r.Context(), userInfo.Id, creatorId)
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *UserHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	uv, err := h.authClient.CheckUserVersion(r.Context(), &generatedAuth.AccessDetails{
		Login:       userDataJWT.Login,
		Id:          userDataJWT.Id.String(),
		UserVersion: int64(userDataJWT.UserVersion),
	})
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if len(uv.Error) != 0 {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if r.Method == http.MethodGet {
		tokenCSRF, err := token.GetCSRFToken(models.User{Login: userDataJWT.Login, Id: userDataJWT.Id, UserVersion: userDataJWT.UserVersion})
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, nil)
			return
		}
		utils.ResponseWithCSRF(w, tokenCSRF)
		return
	}
	// check CSRF token
	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil || *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	subUUID, ok := mux.Vars(r)["sub-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	subId, err := uuid.Parse(subUUID)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	subscription := models.SubscriptionDetails{}

	err = easyjson.UnmarshalFromReader(r.Body, &subscription)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if subscription.MonthCount <= 0 {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	subscription.UserID = userDataJWT.Id
	subscription.Id = subId

	err = h.usecase.Subscribe(r.Context(), subscription)
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userInfo, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	userProfile, err := h.usecase.GetProfile(r.Context(), *userInfo)
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err == models.NotFound {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	userProfile.Sanitize()

	utils.Response(w, http.StatusOK, userProfile)
}

func (h *UserHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	userInfo, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	homePage, err := h.usecase.GetHomePage(r.Context(), *userInfo)
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err == models.NotFound {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	homePage.Sanitize()
	utils.Response(w, http.StatusOK, homePage)
}

func (h *UserHandler) UpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	uv, err := h.authClient.CheckUserVersion(r.Context(), &generatedAuth.AccessDetails{
		Login:       userDataJWT.Login,
		Id:          userDataJWT.Id.String(),
		UserVersion: int64(userDataJWT.UserVersion),
	})
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if len(uv.Error) != 0 {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if r.Method == http.MethodGet {
		tokenCSRF, err := token.GetCSRFToken(models.User{Login: userDataJWT.Login, Id: userDataJWT.Id, UserVersion: userDataJWT.UserVersion})
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, nil)
			return
		}
		utils.ResponseWithCSRF(w, tokenCSRF)
		return
	}
	// check CSRF token
	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil || *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	err = r.ParseMultipartForm(4 << 20) // maxMemory
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	file, fileTmp, err := r.FormFile("upload")
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	// проверка типа файла
	buf, _ := io.ReadAll(file)
	file.Close()
	if file, err = fileTmp.Open(); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if http.DetectContentType(buf) != "image/jpeg" && http.DetectContentType(buf) != "image/png" {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	defer file.Close()

	var oldName uuid.UUID
	oldName, err = uuid.Parse(r.PostFormValue("path"))
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if oldName != uuid.Nil {
		err = os.Remove(filepath.Join(models.FolderPath, fmt.Sprintf("%s.jpg", oldName.String())))
		if err != nil {
			utils.Response(w, http.StatusBadRequest, nil)
		}
	}

	name, err := h.usecase.UpdatePhoto(r.Context(), *userDataJWT)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
	}

	f, err := os.Create(fmt.Sprintf("%s%s.jpg", models.FolderPath, name.String()))
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, file); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
	}

	utils.Response(w, http.StatusOK, name)
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	uv, err := h.authClient.CheckUserVersion(r.Context(), &generatedAuth.AccessDetails{
		Login:       userDataJWT.Login,
		Id:          userDataJWT.Id.String(),
		UserVersion: int64(userDataJWT.UserVersion),
	})
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if len(uv.Error) != 0 {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if r.Method == http.MethodGet {
		tokenCSRF, err := token.GetCSRFToken(models.User{Login: userDataJWT.Login, Id: userDataJWT.Id, UserVersion: userDataJWT.UserVersion})
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, nil)
			return
		}
		utils.ResponseWithCSRF(w, tokenCSRF)
		return
	}
	// check CSRF token
	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil || *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	updPwd := models.UpdatePasswordInfo{}

	err = easyjson.UnmarshalFromReader(r.Body, &updPwd)
	if err != nil || !(models.User{PasswordHash: updPwd.NewPassword}.UserPasswordIsValid()) {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if updPwd.NewPassword == updPwd.OldPassword {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	_, err = h.authClient.CheckUser(r.Context(), &generatedAuth.User{
		Login:        userDataJWT.Login,
		PasswordHash: updPwd.OldPassword,
	})
	if err == models.WrongPassword {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	encryptedPwd, err := h.authClient.EncryptPwd(r.Context(), &generatedAuth.EncryptPwdMg{Password: updPwd.NewPassword})
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err = h.usecase.UpdatePassword(r.Context(), userDataJWT.Id, encryptedPwd.Password); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	tokenJWT, err := h.authClient.SignIn(r.Context(), &generatedAuth.LoginUser{
		Login:        userDataJWT.Login,
		PasswordHash: updPwd.NewPassword,
	})
	if len(tokenJWT.Error) != 0 {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Cookie(w, tokenJWT.Cookie, "SSID")
	utils.Response(w, http.StatusOK, nil)
}

func (h *UserHandler) Donate(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	uv, err := h.authClient.CheckUserVersion(r.Context(), &generatedAuth.AccessDetails{
		Login:       userDataJWT.Login,
		Id:          userDataJWT.Id.String(),
		UserVersion: int64(userDataJWT.UserVersion),
	})
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if len(uv.Error) != 0 {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	if r.Method == http.MethodGet {
		tokenCSRF, err := token.GetCSRFToken(models.User{Login: userDataJWT.Login, Id: userDataJWT.Id, UserVersion: userDataJWT.UserVersion})
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, nil)
			return
		}
		utils.ResponseWithCSRF(w, tokenCSRF)
		utils.Response(w, http.StatusOK, nil)
		return
	}
	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil || *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	donateInfo := models.Donate{}

	err = easyjson.UnmarshalFromReader(r.Body, &donateInfo)
	if err != nil || donateInfo.MoneyCount < 1 {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	newMoneyCount, err := h.usecase.Donate(r.Context(), donateInfo, userDataJWT.Id)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, newMoneyCount)
}

func (h *UserHandler) UpdateData(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	uv, err := h.authClient.CheckUserVersion(r.Context(), &generatedAuth.AccessDetails{
		Login:       userDataJWT.Login,
		Id:          userDataJWT.Id.String(),
		UserVersion: int64(userDataJWT.UserVersion),
	})
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if len(uv.Error) != 0 {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if r.Method == http.MethodGet {
		tokenCSRF, err := token.GetCSRFToken(models.User{Login: userDataJWT.Login, Id: userDataJWT.Id, UserVersion: userDataJWT.UserVersion})
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, nil)
			return
		}
		utils.ResponseWithCSRF(w, tokenCSRF)
		return
	}

	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil || *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	updProfile := models.UpdateProfileInfo{}

	err = easyjson.UnmarshalFromReader(r.Body, &updProfile)
	if err != nil || !(models.User{Name: updProfile.Name}.UserNameIsValid() && models.User{Login: updProfile.Login}.UserLoginIsValid()) {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err = h.usecase.UpdateProfileInfo(r.Context(), updProfile, userDataJWT.Id); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *UserHandler) BecomeCreator(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	uv, err := h.authClient.CheckUserVersion(r.Context(), &generatedAuth.AccessDetails{
		Login:       userDataJWT.Login,
		Id:          userDataJWT.Id.String(),
		UserVersion: int64(userDataJWT.UserVersion),
	})
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if len(uv.Error) != 0 {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	if r.Method == http.MethodGet {
		tokenCSRF, err := token.GetCSRFToken(models.User{Login: userDataJWT.Login, Id: userDataJWT.Id, UserVersion: userDataJWT.UserVersion})
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, nil)
			return
		}
		utils.ResponseWithCSRF(w, tokenCSRF)
		return
	}

	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil || *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	if _, isAlsoCreator, err := h.usecase.CheckIfCreator(r.Context(), userDataJWT.Id); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	} else if isAlsoCreator {
		utils.Response(w, http.StatusConflict, nil)
		return
	}

	authorInfo := models.BecameCreatorInfo{}

	err = easyjson.UnmarshalFromReader(r.Body, &authorInfo)
	if err != nil || !(authorInfo.IsValid()) {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creatorId, err := h.usecase.BecomeCreator(r.Context(), authorInfo, userDataJWT.Id)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, creatorId)
}
