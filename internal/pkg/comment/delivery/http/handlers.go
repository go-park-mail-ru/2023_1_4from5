package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
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

	out, err = h.creatorClient.CreateComment(r.Context(), &generatedCreator.Comment{
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

	commentInfo := models.Comment{}
	if err = easyjson.UnmarshalFromReader(r.Body, &commentInfo); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	commentIDtmp, ok := mux.Vars(r)["comment-uuid"]

	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	commentInfo.CommentID, err = uuid.Parse(commentIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	isCommentOwner, err := h.creatorClient.IsCommentOwner(r.Context(), &generatedCreator.Comment{
		Id:     commentInfo.CommentID.String(),
		UserId: userDataJWT.Id.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if isCommentOwner.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if isCommentOwner.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !isCommentOwner.Flag {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	out, err := h.creatorClient.DeleteComment(r.Context(), &generatedCreator.Comment{
		Id:     commentInfo.CommentID.String(),
		PostID: commentInfo.PostID.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if out.Error != "" {
		utils.Response(w, http.StatusInternalServerError, out.Error)
		return
	}

	utils.Response(w, http.StatusOK, nil)

}

func (h *CommentHandler) EditComment(w http.ResponseWriter, r *http.Request) {

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

	commentIDtmp, ok := mux.Vars(r)["comment-uuid"]

	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	commentInfo.CommentID, err = uuid.Parse(commentIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	isCommentOwner, err := h.creatorClient.IsCommentOwner(r.Context(), &generatedCreator.Comment{
		Id:     commentInfo.CommentID.String(),
		UserId: userDataJWT.Id.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if isCommentOwner.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if isCommentOwner.Error != "" {
		utils.Response(w, http.StatusInternalServerError, isCommentOwner.Error)
		return
	}

	if !isCommentOwner.Flag {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	out, err := h.creatorClient.EditComment(r.Context(), &generatedCreator.Comment{
		Id:   commentInfo.CommentID.String(),
		Text: commentInfo.Text,
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, err)
		return
	}

	if out.Error != "" {
		utils.Response(w, http.StatusInternalServerError, err)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *CommentHandler) AddLike(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	commentInfo := models.Comment{}
	if err = easyjson.UnmarshalFromReader(r.Body, &commentInfo); err != nil {
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

	commentIDtmp, ok := mux.Vars(r)["comment-uuid"]

	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	commentInfo.CommentID, err = uuid.Parse(commentIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	outLike, err := h.creatorClient.AddLikeComment(r.Context(), &generatedCreator.Comment{
		Id:     commentInfo.CommentID.String(),
		PostID: commentInfo.PostID.String(),
		UserId: userDataJWT.Id.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if outLike.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, outLike.Error)
		return
	}

	if outLike.Error != "" {
		utils.Response(w, http.StatusInternalServerError, outLike.Error)
		return
	}

	utils.Response(w, http.StatusOK, outLike.LikesCount)
}

func (h *CommentHandler) RemoveLike(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	commentInfo := models.Comment{}

	commentIDtmp, ok := mux.Vars(r)["comment-uuid"]

	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	commentInfo.CommentID, err = uuid.Parse(commentIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	outLike, err := h.creatorClient.RemoveLikeComment(r.Context(), &generatedCreator.Comment{
		Id:     commentInfo.CommentID.String(),
		UserId: userDataJWT.Id.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if outLike.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, outLike.Error)
		return
	}

	if outLike.Error != "" {
		utils.Response(w, http.StatusInternalServerError, outLike.Error)
		return
	}

	utils.Response(w, http.StatusOK, outLike.LikesCount)
}
