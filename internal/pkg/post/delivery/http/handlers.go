package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type PostHandler struct {
	usecase           post.PostUsecase
	authUsecase       auth.AuthUsecase
	attachmentUsecase attachment.AttachmentUsecase
	logger            *zap.SugaredLogger
}

func NewPostHandler(uc post.PostUsecase, auc auth.AuthUsecase, attuc attachment.AttachmentUsecase, logger *zap.SugaredLogger) *PostHandler {
	return &PostHandler{
		usecase:           uc,
		authUsecase:       auc,
		attachmentUsecase: attuc,
		logger:            logger,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {

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
		utils.Cookie(w, tokenCSRF, "X-CSRF-Token")
		return
	}

	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil || *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(models.MaxFormSize))
	err = r.ParseMultipartForm(models.MaxFormSize)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	postValues := r.MultipartForm.Value
	var postData models.PostCreationData

	_, ok := postValues["creator"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if postData.Creator, err = uuid.Parse(postValues["creator"][0]); err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	ok, err = h.usecase.IsCreator(r.Context(), userDataJWT.Id, postData.Creator)

	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !ok {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	_, ok = postValues["text"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	postData.Text = postValues["text"][0]
	if len(postData.Text) > 4000 {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	_, ok = postValues["title"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	postData.Title = postValues["title"][0]
	if len(postData.Title) > 40 {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if _, ok := r.MultipartForm.File["attachments"]; ok {
		if len(r.MultipartForm.File["attachments"]) > models.MaxFiles {
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}
		for _, headers := range r.MultipartForm.File {
			for _, header := range headers {
				if header.Size > int64(models.MaxFileSize) {
					utils.Response(w, http.StatusBadRequest, nil)
					return
				}
			}
		}

		postFilesTmp := r.MultipartForm.File["attachments"]

		postData.Attachments = make([]models.AttachmentData, len(postFilesTmp))
		for i, file := range postFilesTmp {
			tmpFile, err := file.Open()
			if err != nil {
				utils.Response(w, http.StatusBadRequest, nil)
				return
			}
			buf, _ := io.ReadAll(tmpFile)

			if err = tmpFile.Close(); err != nil {
				h.logger.Error(err)
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}

			if postData.Attachments[i].Data, err = file.Open(); err != nil {
				utils.Response(w, http.StatusBadRequest, nil)
				return
			}
			postData.Attachments[i].Type = http.DetectContentType(buf)
			postData.Attachments[i].Id = uuid.New()
		}

		if err = h.attachmentUsecase.CreateAttachments(r.Context(), postData.Attachments...); err == models.Unsupported {
			utils.Response(w, http.StatusUnsupportedMediaType, nil)
			return
		} else if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}

	}

	tmpSubs, ok := postValues["subscriptions"]
	if ok {
		postData.AvailableSubscriptions = make([]uuid.UUID, len(tmpSubs))
		for i, sub := range tmpSubs {
			if postData.AvailableSubscriptions[i], err = uuid.Parse(sub); err != nil {
				utils.Response(w, http.StatusBadRequest, nil)
				return
			}
		}
	}

	postData.Id = uuid.New()
	if err := h.usecase.CreatePost(r.Context(), postData); err != nil {
		_ = h.attachmentUsecase.DeleteAttachments(postData.Attachments...)
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *PostHandler) AddLike(w http.ResponseWriter, r *http.Request) {
	userData, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	if _, err := h.authUsecase.CheckUserVersion(r.Context(), *userData); err != nil {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	var like models.Like
	err = easyjson.UnmarshalFromReader(r.Body, &like)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	postOwner, err := h.usecase.IsPostOwner(r.Context(), userData.Id, like.PostID)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if postOwner {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	like, err = h.usecase.AddLike(r.Context(), userData.Id, like.PostID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, like)
}

func (h *PostHandler) RemoveLike(w http.ResponseWriter, r *http.Request) {
	userData, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if _, err := h.authUsecase.CheckUserVersion(r.Context(), *userData); err != nil {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	var like models.Like
	err = easyjson.UnmarshalFromReader(r.Body, &like)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	like, err = h.usecase.RemoveLike(r.Context(), userData.Id, like.PostID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err == models.InternalError {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, like)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
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
		utils.Cookie(w, tokenCSRF, "X-CSRF-Token")
		return
	}

	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	postIDtmp, ok := mux.Vars(r)["post-uuid"]

	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	postID, err := uuid.Parse(postIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	ok, err = h.usecase.IsPostOwner(r.Context(), (*userDataJWT).Id, postID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !ok {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	err = h.attachmentUsecase.DeleteAttachmentsByPostID(r.Context(), postID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err = h.usecase.DeletePost(r.Context(), postID); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	postIDtmp, ok := mux.Vars(r)["post-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	postID, err := uuid.Parse(postIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	post, err := h.usecase.GetPost(r.Context(), postID, userDataJWT.Id)

	if err == models.Forbbiden {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	post.Sanitize()

	utils.Response(w, http.StatusOK, post)
}

func (h *PostHandler) EditPost(w http.ResponseWriter, r *http.Request) {
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
		utils.Cookie(w, tokenCSRF, "X-CSRF-Token")
		return
	}

	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	postIDtmp, ok := mux.Vars(r)["post-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	postEditData := models.PostEditData{}
	postEditData.Id, err = uuid.Parse(postIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	ok, err = h.usecase.IsPostOwner(r.Context(), (*userDataJWT).Id, postEditData.Id)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !ok {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	if err = easyjson.UnmarshalFromReader(r.Body, &postEditData); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if len(postEditData.Title) > 40 || len(postEditData.Text) > 4000 {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err = h.usecase.EditPost(r.Context(), postEditData); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *PostHandler) AddAttach(w http.ResponseWriter, r *http.Request) {
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
		utils.Cookie(w, tokenCSRF, "X-CSRF-Token")
		return
	}

	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	postIDtmp, ok := mux.Vars(r)["post-uuid"]
	postID, err := uuid.Parse(postIDtmp)

	ok, err = h.usecase.IsPostOwner(r.Context(), (*userDataJWT).Id, postID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !ok {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(models.MaxFileSize))
	err = r.ParseMultipartForm(models.MaxFormSize)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if len(r.MultipartForm.File["attachment"]) > 1 {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if r.MultipartForm.File["attachment"][0].Size > int64(models.MaxFileSize) {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	postFilesTmp := r.MultipartForm.File["attachment"]
	attach := models.AttachmentData{}
	tmpFile, err := postFilesTmp[0].Open()
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	buf, _ := io.ReadAll(tmpFile)

	if err = tmpFile.Close(); err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if attach.Data, err = postFilesTmp[0].Open(); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	attach.Type = http.DetectContentType(buf)
	attach.Id = uuid.New()

	if err = h.attachmentUsecase.CreateAttachments(r.Context(), attach); err == models.Unsupported {
		utils.Response(w, http.StatusUnsupportedMediaType, nil)
		return
	} else if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *PostHandler) DeleteAttach(w http.ResponseWriter, r *http.Request) {
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
		utils.Cookie(w, tokenCSRF, "X-CSRF-Token")
		return
	}

	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if *userDataCSRF != *userDataJWT {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	postIdTmp, ok := mux.Vars(r)["post-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	postID, err := uuid.Parse(postIdTmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	attachInfo := models.Attachment{}
	err = easyjson.UnmarshalFromReader(r.Body, &attachInfo)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	ok, err = h.usecase.IsPostOwner(r.Context(), (*userDataJWT).Id, postID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !ok {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	err = h.attachmentUsecase.DeleteAttachment(r.Context(), attachInfo.Id, postID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	err = h.attachmentUsecase.DeleteAttachments(models.AttachmentData{Id: attachInfo.Id, Type: attachInfo.Type})
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}
