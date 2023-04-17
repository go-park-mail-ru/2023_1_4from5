package http

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/subscription"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type SubscriptionHandler struct {
	usecase     subscription.SubscriptionUsecase
	authUsecase auth.AuthUsecase
	userUsecase user.UserUsecase
	logger      *zap.SugaredLogger
}

func NewSubscriptionHandler(uc subscription.SubscriptionUsecase, auc auth.AuthUsecase, uuc user.UserUsecase, logger *zap.SugaredLogger) *SubscriptionHandler {
	return &SubscriptionHandler{
		usecase:     uc,
		authUsecase: auc,
		userUsecase: uuc,
		logger:      logger,
	}
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {

	userDataJWT, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if _, err := h.authUsecase.CheckUserVersion(r.Context(), *userDataJWT); err != nil {
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
	if !subscriptionInfo.IsValid() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creatorID, isCreator, err := h.userUsecase.CheckIfCreator(r.Context(), userDataJWT.Id)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	} else if !isCreator {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	subscriptionInfo.Creator = creatorID
	subscriptionInfo.Id = uuid.New()

	if err = h.usecase.CreateSubscription(r.Context(), subscriptionInfo); err != nil {
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

	if _, err := h.authUsecase.CheckUserVersion(r.Context(), *userDataJWT); err != nil {
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

	creatorID, isCreator, err := h.userUsecase.CheckIfCreator(r.Context(), userDataJWT.Id)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	} else if !isCreator {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if err = h.usecase.DeleteSubscription(r.Context(), subscriptionID, creatorID); errors.Is(err, models.Forbbiden) {
		utils.Response(w, http.StatusForbidden, nil)
		return
	} else if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}
