package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/notification"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CreatorHandler struct {
	creatorClient   generatedCreator.CreatorServiceClient
	authClient      generatedAuth.AuthServiceClient
	notificationApp notification.NotificationApp
	logger          *zap.SugaredLogger
}

func NewCreatorHandler(creatorClient generatedCreator.CreatorServiceClient, authClient generatedAuth.AuthServiceClient, na notification.NotificationApp, logger *zap.SugaredLogger) *CreatorHandler {
	return &CreatorHandler{
		creatorClient:   creatorClient,
		authClient:      authClient,
		notificationApp: na,
		logger:          logger,
	}
}

func (h *CreatorHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	creatorID, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorID.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorID.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	balance, err := h.creatorClient.GetCreatorBalance(r.Context(), &generatedCommon.UUIDMessage{Value: creatorID.Value})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if balance.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, balance.Balance)
}

func (h *CreatorHandler) TransferMoney(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	creatorID, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorID.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorID.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	transfer := models.CreatorTransfer{}
	err = easyjson.UnmarshalFromReader(r.Body, &transfer)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	balance, err := h.creatorClient.GetCreatorBalance(r.Context(), &generatedCommon.UUIDMessage{Value: creatorID.Value})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if balance.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if balance.Balance < transfer.Money {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	reqID, err := h.requestPayment(transfer)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if err = h.processPayment(reqID); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	balance, err = h.creatorClient.UpdateBalance(r.Context(), &generatedCreator.CreatorTransfer{
		CreatorID: creatorID.Value,
		Money:     transfer.Money,
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if balance.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, balance)
}

func (h *CreatorHandler) requestPayment(transfer models.CreatorTransfer) (models.PaymentResponse, error) {
	paymentToken, flag := os.LookupEnv("PAYMENT_TOKEN")
	if !flag {
		return models.PaymentResponse{}, errors.New("no payment token")
	}

	method := "POST"
	payload := strings.NewReader("pattern_id=p2p&to=" + transfer.PhoneNumber + "&identifier_type=phone" + "&amount=" + fmt.Sprintf("%f", transfer.Money) + "&comment=Payment%20to%20authorID%3D%3CUUID%3E&message=Payment%20from%20SubMe")

	client := &http.Client{}
	req, err := http.NewRequest(method, models.RequestPaymentURL, payload)

	if err != nil {
		h.logger.Error(err)
		return models.PaymentResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+paymentToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		h.logger.Error(err)
		return models.PaymentResponse{}, err
	}

	var reqID models.PaymentResponse
	err = easyjson.UnmarshalFromReader(res.Body, &reqID)
	if err != nil {
		h.logger.Error(err)
		return models.PaymentResponse{}, err
	}
	return reqID, nil
}

func (h *CreatorHandler) processPayment(reqID models.PaymentResponse) error {
	method := "POST"
	paymentToken, flag := os.LookupEnv("PAYMENT_TOKEN")
	if !flag {
		return errors.New("no payment token")
	}

	payload2 := strings.NewReader(fmt.Sprintf("request_id=%s", reqID.RequestID))

	req, err := http.NewRequest(method, models.ProcessPaymentURL, payload2)
	if err != nil {
		h.logger.Error(err)
		return err
	}
	req.Header.Add("Authorization", "Bearer "+paymentToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		h.logger.Error(err)
		return err
	}
	defer res.Body.Close()
	return nil
}

func (h *CreatorHandler) SubscribeCreatorToNotifications(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	token := models.NotificationToken{}
	err = easyjson.UnmarshalFromReader(r.Body, &token)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creatorID, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorID.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorID.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	err = h.notificationApp.AddUserToNotificationTopic(fmt.Sprintf("%s-%s", creatorID.Value, "creator"), token, context.Background())
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *CreatorHandler) UnsubscribeCreatorNotifications(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	token := models.NotificationToken{}
	err = easyjson.UnmarshalFromReader(r.Body, &token)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creatorID, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorID.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorID.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	err = h.notificationApp.RemoveUserFromNotificationTopic(fmt.Sprintf("%s-%s", creatorID.Value, "creator"), token, context.Background())
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *CreatorHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
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

	out, err := h.creatorClient.GetFeed(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	feed := make([]models.Post, len(out.Posts))

	for i, post := range out.Posts {
		err = feed[i].PostToModel(post)
		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		feed[i].Sanitize()
	}

	utils.Response(w, http.StatusOK, feed)
}

func (h *CreatorHandler) GetAllCreators(w http.ResponseWriter, r *http.Request) {
	out, err := h.creatorClient.GetAllCreators(r.Context(), &generatedCommon.Empty{Error: ""})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	creators := make([]models.Creator, len(out.Creators))

	for i, creator := range out.Creators {
		err = creators[i].CreatorToModel(creator)
		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creators[i].Sanitize()
	}
	utils.Response(w, http.StatusOK, creators)
}

func (h *CreatorHandler) UpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
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

	creatorId, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorId.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorId.Error != "" {
		h.logger.Error(creatorId.Error)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorId.Value == uuid.Nil.String() {
		utils.Response(w, http.StatusBadRequest, nil)
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
		err = os.Remove(filepath.Join(models.FolderPath, fmt.Sprintf("%s.jpg", oldName.String())))
		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}
	}

	name, err := h.creatorClient.UpdateProfilePhoto(r.Context(), &generatedCommon.UUIDMessage{
		Value: creatorId.Value})
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

func (h *CreatorHandler) UpdateCoverPhoto(w http.ResponseWriter, r *http.Request) {
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

	creatorId, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorId.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorId.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorId.Value == uuid.Nil.String() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	err = r.ParseMultipartForm(models.MaxFileSize)
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
		err = os.Remove(models.FolderPath + fmt.Sprintf("%s.jpg", oldName.String()))
		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}
	}

	name, err := h.creatorClient.UpdateCoverPhoto(r.Context(), &generatedCommon.UUIDMessage{
		Value: creatorId.Value})
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

func (h *CreatorHandler) FindCreator(w http.ResponseWriter, r *http.Request) {
	keyword, ok := mux.Vars(r)["keyword"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	out, err := h.creatorClient.FindCreators(r.Context(), &generatedCreator.KeywordMessage{Keyword: keyword})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	creators := make([]models.Creator, len(out.Creators))

	for i, creator := range out.Creators {
		err = creators[i].CreatorToModel(creator)
		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creators[i].Sanitize()
	}

	utils.Response(w, http.StatusOK, creators)
}

func (h *CreatorHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	creatorUUID, ok := mux.Vars(r)["creator-uuid"]
	if !ok {
		fmt.Println(creatorUUID)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	userInfo := models.AccessDetails{}
	tmpUserInfo, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		tmpUserInfo = &userInfo
	}

	_, err = uuid.Parse(creatorUUID)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creatorPage, err := h.creatorClient.GetPage(r.Context(), &generatedCreator.UserCreatorMessage{
		UserID:    tmpUserInfo.Id.String(),
		CreatorID: creatorUUID,
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorPage.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if creatorPage.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	var page models.CreatorPage
	page, err = h.CreatorPageToModel(creatorPage)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	page.Sanitize()

	utils.Response(w, http.StatusOK, page)
}

func (h *CreatorHandler) CreateAim(w http.ResponseWriter, r *http.Request) {
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
	aimInfo := models.Aim{}
	err = easyjson.UnmarshalFromReader(r.Body, &aimInfo)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err = aimInfo.IsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	isCreator, err := h.creatorClient.IsCreator(r.Context(), &generatedCreator.UserCreatorMessage{
		UserID:    userDataJWT.Id.String(),
		CreatorID: aimInfo.Creator.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if isCreator.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if isCreator.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !isCreator.Flag {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	out, err := h.creatorClient.CreateAim(r.Context(), &generatedCreator.Aim{
		Creator:     aimInfo.Creator.String(),
		Description: aimInfo.Description,
		MoneyNeeded: aimInfo.MoneyNeeded,
		MoneyGot:    aimInfo.MoneyGot,
	})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if out.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}

func (h *CreatorHandler) UpdateCreatorData(w http.ResponseWriter, r *http.Request) {
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

	updCreator := models.BecameCreatorInfo{}

	creatorID, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorID.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorID.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	err = easyjson.UnmarshalFromReader(r.Body, &updCreator)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err = updCreator.IsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	out, err := h.creatorClient.UpdateCreatorData(r.Context(), &generatedCreator.UpdateCreatorInfo{
		CreatorName: updCreator.Name,
		Description: updCreator.Description,
		CreatorID:   creatorID.Value,
	})
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

func (h *CreatorHandler) CreatorPageToModel(creatorPage *generatedCreator.CreatorPage) (models.CreatorPage, error) {
	var page models.CreatorPage

	err := page.CreatorInfo.CreatorToModel(creatorPage.CreatorInfo)
	if err != nil {
		h.logger.Error(err)
		return models.CreatorPage{}, models.InternalError
	}

	for _, sub := range creatorPage.Subscriptions {
		var subscription models.Subscription
		err = subscription.ProtoSubscriptionToModel(sub)
		if err != nil {
			h.logger.Error(err)
			return models.CreatorPage{}, models.InternalError
		}
		page.Subscriptions = append(page.Subscriptions, subscription)
	}

	page.IsMyPage = creatorPage.IsMyPage
	page.Follows = creatorPage.Follows
	err = page.Aim.AimToModel(creatorPage.AimInfo)
	if err != nil {
		h.logger.Error(err)
		return models.CreatorPage{}, models.InternalError
	}

	page.Posts = make([]models.Post, len(creatorPage.Posts))
	for i, post := range creatorPage.Posts {
		err = page.Posts[i].PostToModel(post)
		if err != nil {
			h.logger.Error(err)
			return models.CreatorPage{}, models.InternalError
		}
	}
	return page, nil
}

func (h *CreatorHandler) DeleteCoverPhoto(w http.ResponseWriter, r *http.Request) {
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

	creatorId, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorId.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorId.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
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

	out, err := h.creatorClient.DeleteCoverPhoto(r.Context(), &generatedCommon.UUIDMessage{
		Value: creatorId.Value})
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

func (h *CreatorHandler) DeleteProfilePhoto(w http.ResponseWriter, r *http.Request) {
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

	creatorId, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorId.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorId.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
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

	out, err := h.creatorClient.DeleteProfilePhoto(r.Context(), &generatedCommon.UUIDMessage{
		Value: creatorId.Value})
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

func (h *CreatorHandler) StatisticsFirstDate(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	creatorID, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorID.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorID.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	date, err := h.creatorClient.StatisticsFirstDate(r.Context(), &generatedCommon.UUIDMessage{Value: creatorID.Value})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if date.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, date.Date)

}

func (h *CreatorHandler) Statistics(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	monthGap := models.StatisticsDates{}
	err = easyjson.UnmarshalFromReader(r.Body, &monthGap)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if monthGap.FirstMonth.Unix() > monthGap.SecondMonth.Unix() {
		utils.Response(w, http.StatusBadRequest, "First month can't be bigger than second")
		return
	}

	creatorID, err := h.creatorClient.CheckIfCreator(r.Context(), &generatedCommon.UUIDMessage{Value: userDataJWT.Id.String()})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if creatorID.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if creatorID.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	stat, err := h.creatorClient.Statistics(r.Context(), &generatedCreator.StatisticsInput{
		CreatorId:  creatorID.Value,
		FirstDate:  monthGap.FirstMonth.Format(time.RFC3339),
		SecondDate: monthGap.SecondMonth.Format(time.RFC3339),
	})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if stat.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	var statistics models.Statistics

	if err = statistics.StatToModel(stat); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	statistics.CreatorId, _ = uuid.Parse(creatorID.Value)

	utils.Response(w, http.StatusOK, statistics)
}
