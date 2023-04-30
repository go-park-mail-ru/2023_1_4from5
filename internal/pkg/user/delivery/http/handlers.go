package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	generatedUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UserHandler struct {
	userClient generatedUser.UserServiceClient
	authClient generatedAuth.AuthServiceClient
	logger     *zap.SugaredLogger
}

func NewUserHandler(userClient generatedUser.UserServiceClient, auc generatedAuth.AuthServiceClient, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		userClient: userClient,
		authClient: auc,
		logger:     logger,
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
		h.logger.Error(err)
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

	out, err := h.userClient.Unfollow(r.Context(), &generatedUser.FollowMessage{
		UserID:    userInfo.Id.String(),
		CreatorID: creatorId.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if out.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
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
	userProfile, err := h.userClient.GetProfile(r.Context(),
		&generatedCommon.UUIDMessage{Value: userInfo.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if userProfile.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if userProfile.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	reg, err := time.Parse("2006-01-02 15:04:05 -0700 -0700", userProfile.Registration)

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	photo, err := uuid.Parse(userProfile.ProfilePhoto)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	creatorID, err := uuid.Parse(userProfile.CreatorID)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	profile := models.UserProfile{
		Login:        userProfile.Login,
		Registration: reg,
		ProfilePhoto: photo,
		CreatorId:    creatorID,
		IsCreator:    userProfile.IsCreator,
		Name:         userProfile.Name,
	}

	profile.Sanitize()

	utils.Response(w, http.StatusOK, profile)
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
		UserVersion: userDataJWT.UserVersion,
	})
	if err != nil {
		h.logger.Error(err)
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

	err = r.ParseMultipartForm(models.MaxFileSize) // maxMemory
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
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if http.DetectContentType(buf) != "image/jpeg" && http.DetectContentType(buf) != "image/png" && http.DetectContentType(buf) != "image/jpg" {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	defer file.Close()

	var oldName uuid.UUID
	oldName, err = uuid.Parse(r.PostFormValue("path"))
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if oldName != uuid.Nil {
		err = os.Remove(filepath.Join(models.FolderPath, fmt.Sprintf("%s.jpg", oldName.String())))
		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}
	}

	name, err := h.userClient.UpdatePhoto(r.Context(), &generatedCommon.UUIDMessage{
		Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if name.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	f, err := os.Create(fmt.Sprintf("%s%s.jpg", models.FolderPath, name.Value))
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, file); err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, name.Value)
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
		UserVersion: userDataJWT.UserVersion,
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
	out, err := h.userClient.UpdatePassword(r.Context(), &generatedUser.UpdatePasswordMessage{
		Password: encryptedPwd.Password,
		UserID:   userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if out.Error != "" {
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
		UserVersion: userDataJWT.UserVersion,
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

	newMoneyCount, err := h.userClient.Donate(r.Context(), &generatedUser.DonateMessage{MoneyCount: donateInfo.MoneyCount,
		CreatorID: donateInfo.CreatorID.String(),
		UserID:    userDataJWT.Id.String()})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if newMoneyCount.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if newMoneyCount.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, newMoneyCount.MoneyCount)
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
		UserVersion: userDataJWT.UserVersion,
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

	out, err := h.userClient.UpdateProfileInfo(r.Context(), &generatedUser.UpdateProfileInfoMessage{Login: updProfile.Login,
		Name: updProfile.Name, UserID: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error != "" {
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
		UserVersion: userDataJWT.UserVersion,
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

	out, err := h.userClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.IsCreator {
		utils.Response(w, http.StatusConflict, nil)
		return
	}

	authorInfo := models.BecameCreatorInfo{}

	err = easyjson.UnmarshalFromReader(r.Body, &authorInfo)
	if err != nil || !(authorInfo.IsValid()) {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creatorId, err := h.userClient.BecomeCreator(r.Context(), &generatedUser.BecameCreatorInfoMessage{Name: authorInfo.Name,
		Description: authorInfo.Description, UserID: userDataJWT.Id.String()})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorId.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, creatorId.Value)
}

func (h *UserHandler) UserSubscriptions(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	out, err := h.userClient.UserSubscriptions(r.Context(), &generatedCommon.UUIDMessage{
		Value: userDataJWT.Id.String()})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	subs := make([]models.Subscription, len(out.Subscriptions))

	for i, v := range out.Subscriptions {
		subId, err := uuid.Parse(v.Id)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creatorId, err := uuid.Parse(v.Creator)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creatorPhoto, err := uuid.Parse(v.CreatorPhoto)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		subs[i] = models.Subscription{
			Id:           subId,
			Creator:      creatorId,
			CreatorPhoto: creatorPhoto,
			CreatorName:  v.CreatorName,
			MonthCost:    v.MonthCost,
			Title:        v.Title,
			Description:  v.Description,
		}

		subs[i].Sanitize()
	}

	utils.Response(w, http.StatusOK, subs)
}

func (h *UserHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	//TODO: проверить соответствие количества денег(и вообще в идеале не класть его и считать из month_count, и вообще должен лететь токен киви какой-нибудь)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	uv, err := h.authClient.CheckUserVersion(r.Context(), &generatedAuth.AccessDetails{
		Login:       userDataJWT.Login,
		Id:          userDataJWT.Id.String(),
		UserVersion: userDataJWT.UserVersion,
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

	_, err = uuid.Parse(subUUID)
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

	out, err := h.userClient.Subscribe(r.Context(), &generatedUser.SubscriptionDetails{
		CreatorID:  subscription.CreatorId.String(),
		UserID:     userDataJWT.Id.String(),
		Id:         subUUID,
		MonthCount: subscription.MonthCount,
		Money:      subscription.Money})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}
