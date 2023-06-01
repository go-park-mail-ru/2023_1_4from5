package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	generatedUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type SubscriptionHandler struct {
	authClient    generatedAuth.AuthServiceClient
	creatorClient generatedCreator.CreatorServiceClient
	userClient    generatedUser.UserServiceClient
	logger        *zap.SugaredLogger
}

func NewSubscriptionHandler(auc generatedAuth.AuthServiceClient, creatorClient generatedCreator.CreatorServiceClient, userClient generatedUser.UserServiceClient, logger *zap.SugaredLogger) *SubscriptionHandler {
	return &SubscriptionHandler{
		authClient:    auc,
		logger:        logger,
		creatorClient: creatorClient,
		userClient:    userClient,
	}
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {

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

	subscriptionInfo := models.Subscription{}
	if err = easyjson.UnmarshalFromReader(r.Body, &subscriptionInfo); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err = subscriptionInfo.IsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
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

	subscriptionInfo.Id = uuid.New()

	out, err := h.creatorClient.CreateSubscription(r.Context(), &generatedCommon.Subscription{
		Id:          subscriptionInfo.Id.String(),
		Creator:     creatorId.Value,
		MonthCost:   subscriptionInfo.MonthCost,
		Title:       subscriptionInfo.Title,
		Description: subscriptionInfo.Description,
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

	utils.Response(w, http.StatusOK, subscriptionInfo)
}

func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionIDTmp, ok := mux.Vars(r)["sub-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	subscriptionID, err := uuid.Parse(subscriptionIDTmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

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

	out, err := h.creatorClient.DeleteSubscription(r.Context(), &generatedCreator.SubscriptionCreatorMessage{
		SubscriptionID: subscriptionID.String(),
		CreatorID:      creatorId.Value,
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error == models.NotFound.Error() {
		utils.Response(w, http.StatusForbidden, nil)
		return
	} else if out.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *SubscriptionHandler) EditSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionID, ok := mux.Vars(r)["sub-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	subscriptionInfo := models.Subscription{}
	_, err := uuid.Parse(subscriptionID)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

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

	if err = easyjson.UnmarshalFromReader(r.Body, &subscriptionInfo); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err = subscriptionInfo.IsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	if subscriptionInfo.Creator == uuid.Nil {
		utils.Response(w, http.StatusBadRequest, nil)
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

	if creatorId.Value != subscriptionInfo.Creator.String() {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	out, err := h.creatorClient.EditSubscription(r.Context(), &generatedCommon.Subscription{
		Id:          subscriptionID,
		Creator:     creatorId.Value,
		MonthCost:   subscriptionInfo.MonthCost,
		Title:       subscriptionInfo.Title,
		Description: subscriptionInfo.Description,
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
