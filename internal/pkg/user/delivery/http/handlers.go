package http

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/notification"
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
	"strconv"
	"strings"
	"time"
)

type UserHandler struct {
	userClient      generatedUser.UserServiceClient
	authClient      generatedAuth.AuthServiceClient
	notificationApp notification.NotificationApp
	logger          *zap.SugaredLogger
}

func NewUserHandler(userClient generatedUser.UserServiceClient, auc generatedAuth.AuthServiceClient, na notification.NotificationApp, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		userClient:      userClient,
		authClient:      auc,
		notificationApp: na,
		logger:          logger,
	}
}

func (h *UserHandler) SubscribeUserToNotifications(w http.ResponseWriter, r *http.Request) {
	creatorID, ok := mux.Vars(r)["creator-uuid"]
	if !ok {
		fmt.Println("wrong")
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	token := models.NotificationToken{}
	err := easyjson.UnmarshalFromReader(r.Body, &token)
	if err != nil {
		fmt.Println(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	err = h.notificationApp.AddUserToNotificationTopic(fmt.Sprintf("%s-%s", creatorID, "user"), token, context.Background())
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *UserHandler) UnsubscribeUserNotifications(w http.ResponseWriter, r *http.Request) {
	creatorID, ok := mux.Vars(r)["creator-uuid"]
	if !ok {
		fmt.Println("wrong")
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	token := models.NotificationToken{}
	err := easyjson.UnmarshalFromReader(r.Body, &token)
	if err != nil {
		fmt.Println(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	err = h.notificationApp.RemoveUserFromNotificationTopic(fmt.Sprintf("%s-%s", creatorID, "user"), token, context.Background())
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, nil)
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

	err = r.ParseMultipartForm(models.MaxFileSize)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	file, fileTmp, err := r.FormFile("upload")
	if err != nil {
		h.logger.Error(err)
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
		err = os.Remove(models.FolderPath + fmt.Sprintf("%s.jpg", oldName.String()))
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
	if err != nil || (models.User{PasswordHash: updPwd.NewPassword}).UserPasswordIsValid() != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	if updPwd.NewPassword == updPwd.OldPassword {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	out1, err := h.authClient.CheckUser(r.Context(), &generatedAuth.User{
		Login:        userDataJWT.Login,
		PasswordHash: updPwd.OldPassword,
	})

	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out1.Error == models.WrongPassword.Error() {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	if out1.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
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
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err = (models.User{Name: updProfile.Name}).UserNameIsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = (models.User{Login: updProfile.Login}.UserLoginIsValid()); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
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
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err = authorInfo.IsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
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

func GetPaymentStringMap() map[string]string {
	result := map[string]string{
		"notification_type": "",
		"operation_id":      "",
		"amount":            "",
		"currency":          "",
		"datetime":          "",
		"sender":            "",
		"codepro":           "",
		"label":             "",
	}
	return result
}

func (h *UserHandler) Payment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	//check sha-1
	paymentStringMap := GetPaymentStringMap()

	sha1Hash := ""
	for key, value := range r.Form {
		if key == "sha1_hash" {
			sha1Hash = value[0]
		}
		if _, ok := paymentStringMap[key]; ok {
			paymentStringMap[key] = value[0]
		}
	}

	paymentSecret, flag := os.LookupEnv("PAYMENT_SECRET")
	if !flag {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	paymentString := strings.Join([]string{paymentStringMap["notification_type"],
		paymentStringMap["operation_id"], paymentStringMap["amount"], paymentStringMap["currency"], paymentStringMap["datetime"], paymentStringMap["sender"],
		paymentStringMap["codepro"], paymentSecret, paymentStringMap["label"]}, "&")

	hash := sha1.New()
	hash.Write([]byte(paymentString))
	paymentStringSHA := hash.Sum(nil)
	if fmt.Sprintf("%x", paymentStringSHA) != sha1Hash {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	paymentInfo := models.PaymentDetails{}
	str := strings.Split(paymentStringMap["label"], ";")
	paymentInfo.Operation = str[0]
	paymentInfo.CreatorId, err = uuid.Parse(str[1])
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if tmp, err := strconv.ParseFloat(paymentStringMap["amount"], 32); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	} else {
		paymentInfo.Money = float32(tmp)
	}

	if paymentInfo.Operation == "subscribe" {

		out, err := h.userClient.Subscribe(r.Context(), &generatedUser.PaymentInfo{PaymentID: paymentInfo.CreatorId.String(),
			Money: paymentInfo.Money})

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

		_ = h.notificationApp.SendUserNotification(models.Notification{
			Topic: fmt.Sprintf("%s-%s", out.CreatorID, "creator"),
			Title: "Новая подписка",
			Body:  fmt.Sprintf("На вас была оформлена подписка %s", out.Name),
		}, r.Context())

		utils.Response(w, http.StatusOK, nil)
	} else if paymentInfo.Operation == "donate" {

		newMoneyCount, err := h.userClient.Donate(r.Context(), &generatedUser.DonateMessage{
			MoneyCount: paymentInfo.Money,
			CreatorID:  paymentInfo.CreatorId.String()})

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

		_ = h.notificationApp.SendUserNotification(models.Notification{
			Topic: fmt.Sprintf("%s-%s", paymentInfo.CreatorId.String(), "creator"),
			Title: "Новый донат",
			Body:  fmt.Sprintf("Вам пришёл новый донат на сумму %f", paymentInfo.Money),
		}, r.Context())

		utils.Response(w, http.StatusOK, nil)
	} else {
		utils.Response(w, http.StatusBadRequest, nil)
	}
}

func (h *UserHandler) AddPaymentInfo(w http.ResponseWriter, r *http.Request) {
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

	subscription.PaymentInfo = uuid.New()
	fmt.Println(subscription)
	out, err := h.userClient.AddPaymentInfo(r.Context(), &generatedUser.SubscriptionDetails{
		CreatorID:   subscription.CreatorId.String(),
		UserID:      userDataJWT.Id.String(),
		Id:          subUUID,
		MonthCount:  subscription.MonthCount,
		PaymentInfo: subscription.PaymentInfo.String()})

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

	utils.Response(w, http.StatusOK, subscription.PaymentInfo)
}

func (h *UserHandler) DeleteProfilePhoto(w http.ResponseWriter, r *http.Request) {
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
	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil || *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	var oldName uuid.UUID
	imageID, ok := mux.Vars(r)["image-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	oldName, err = uuid.Parse(imageID)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if oldName != uuid.Nil {
		err = os.Remove(models.FolderPath + fmt.Sprintf("%s.jpg", oldName.String()))
		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}
	}

	name, err := h.userClient.DeletePhoto(r.Context(), &generatedCommon.UUIDMessage{
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

	utils.Response(w, http.StatusOK, nil)
}

func (h *UserHandler) UserFollows(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	out, err := h.userClient.UserFollows(r.Context(), &generatedCommon.UUIDMessage{
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

	follows := make([]models.Follow, len(out.Follows))

	for i, v := range out.Follows {
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
		follows[i] = models.Follow{
			Creator:      creatorId,
			CreatorPhoto: creatorPhoto,
			CreatorName:  v.CreatorName,
			Description:  v.Description,
		}

		follows[i].Sanitize()
	}

	utils.Response(w, http.StatusOK, follows)
}
