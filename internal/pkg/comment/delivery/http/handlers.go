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

type CommentHandler struct {
	authClient    generatedAuth.AuthServiceClient
	userClient    generatedUser.UserServiceClient
	creatorClient generatedCreator.CreatorServiceClient
	logger        *zap.SugaredLogger
}

func NewCommentHandler(auc generatedAuth.AuthServiceClient, userClient generatedUser.UserServiceClient, creatorClient generatedCreator.CreatorServiceClient, logger *zap.SugaredLogger) *CommentHandler {
	return &CommentHandler{
		authClient:    auc,
		userClient:    userClient,
		creatorClient: creatorClient,
		logger:        logger,
	}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {

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

	commentInfo := models.Comment{}
	if err = easyjson.UnmarshalFromReader(r.Body, &commentInfo); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if !commentInfo.IsValid() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	out, err := h.creatorClient.IsPostAvailable(r.Context(), &generatedCreator.PostUserMessage{
		UserID: userDataJWT.Id.String(),
		PostID: commentInfo.PostID.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if out.Error != "" {
		utils.Response(w, http.StatusForbidden, out.Error)
		return
	}

	commentInfo.CommentID = uuid.New()

	out, err = h.creatorClient.CreateComment(r.Context(), &generatedCommon.Comment{
		Id:     commentInfo.CommentID.String(),
		PostID: commentInfo.PostID.String(),
		UserId: userDataJWT.Id.String(),
		Text:   commentInfo.Text,
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

	utils.Response(w, http.StatusOK, commentInfo)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
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

	postIDtmp, ok := mux.Vars(r)["post-uuid"]

	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	commentID, err := uuid.Parse(postIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	//out, err = h.creatorClient.DeleteComment(r.Context(), &generatedCommon.Comment{
	//	Id:     commentID.String(),
	//	UserId: userDataJWT.Id.String(),
	//})
	//
	//if err != nil {
	//	h.logger.Error(err)
	//	utils.Response(w, http.StatusInternalServerError, nil)
	//	return
	//}
	//
	//if out.Error == models.Forbbiden.Error() {
	//	utils.Response(w, http.StatusForbidden, nil)
	//	return
	//}

	utils.Response(w, http.StatusOK, commentID)

}
