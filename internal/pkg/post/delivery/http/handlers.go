package http

import (
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
)

type PostHandler struct {
	authClient      generatedAuth.AuthServiceClient
	creatorClient   generatedCreator.CreatorServiceClient
	logger          *zap.SugaredLogger
	notificationApp notification.NotificationApp
}

func NewPostHandler(auc generatedAuth.AuthServiceClient, csc generatedCreator.CreatorServiceClient, logger *zap.SugaredLogger, app notification.NotificationApp) *PostHandler {
	return &PostHandler{
		authClient:      auc,
		creatorClient:   csc,
		logger:          logger,
		notificationApp: app,
	}
}

// nolint:gocognit
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
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

	out, err := h.creatorClient.IsCreator(r.Context(), &generatedCreator.UserCreatorMessage{
		UserID:    userDataJWT.Id.String(),
		CreatorID: postData.Creator.String(),
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
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !out.Flag {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	_, ok = postValues["text"]
	if ok {
		postData.Text = postValues["text"][0]
	}
	_, ok = postValues["title"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	postData.Title = postValues["title"][0]
	if err = postData.IsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, ok := r.MultipartForm.File["attachments"]; ok {
		h.logger.Info("Got attachs")
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
		h.logger.Info(postFilesTmp)

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

			attachmentType, err := h.creatorClient.GetFileExtension(r.Context(), &generatedCreator.KeywordMessage{Keyword: postData.Attachments[i].Type})

			if err != nil {
				h.logger.Error(err)
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}

			if !attachmentType.Flag {
				utils.Response(w, http.StatusUnsupportedMediaType, nil)
				return

			}
			postData.Attachments[i].Id = uuid.New()
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
	var attachProto []*generatedCreator.Attachment
	for _, attach := range postData.Attachments {
		attachProto = append(attachProto, &generatedCreator.Attachment{
			ID:   attach.Id.String(),
			Type: attach.Type,
		})
	}

	var subsProto []string
	for _, sub := range postData.AvailableSubscriptions {
		subsProto = append(subsProto, sub.String())
	}

	errMessage, err := h.creatorClient.CreatePost(r.Context(), &generatedCreator.PostCreationData{
		Id:                     postData.Id.String(),
		Creator:                postData.Creator.String(),
		Title:                  postData.Title,
		Text:                   postData.Text,
		Attachments:            attachProto,
		AvailableSubscriptions: subsProto,
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if errMessage.Error != "" {
		_, _ = h.creatorClient.DeleteAttachmentsFiles(r.Context(), &generatedCreator.Attachments{Attachments: attachProto})
		h.logger.Error(errMessage.Error)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	for i, attach := range postData.Attachments {
		attachmentType, err := h.creatorClient.GetFileExtension(r.Context(), &generatedCreator.KeywordMessage{Keyword: attach.Type})

		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}

		if !attachmentType.Flag {
			errMessage, err = h.creatorClient.DeleteAttachmentsFiles(r.Context(), &generatedCreator.Attachments{Attachments: attachProto[:i]})

			if err != nil {
				h.logger.Error(err)
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}

			if errMessage.Error != "" {
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}
			errMessage, err = h.creatorClient.DeletePost(r.Context(), &generatedCommon.UUIDMessage{Value: postData.Id.String()})

			if err != nil {
				h.logger.Error(err)
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}

			if errMessage.Error != "" {
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}

			utils.Response(w, http.StatusUnsupportedMediaType, nil)
			return
		}
		f, err := os.Create(fmt.Sprintf("%s.%s", filepath.Join(models.FolderPath, attach.Id.String()), attachmentType.Extension))
		if err != nil {
			h.logger.Error(err)
			errMessage, err = h.creatorClient.DeleteAttachmentsFiles(r.Context(), &generatedCreator.Attachments{Attachments: attachProto[:i]})

			if err != nil {
				h.logger.Error(err)
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}

			if errMessage.Error != "" {
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}
			_, err = h.creatorClient.DeletePost(r.Context(), &generatedCommon.UUIDMessage{Value: postData.Id.String()})
			if err != nil {
				h.logger.Error(err)
			}
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}

		if _, err := io.Copy(f, attach.Data); err != nil {
			h.logger.Error(err)
			errMessage, err = h.creatorClient.DeleteAttachmentsFiles(r.Context(), &generatedCreator.Attachments{Attachments: attachProto[:i]})

			if err != nil {
				h.logger.Error(err)
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}

			if errMessage.Error != "" {
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}

			_, err = h.creatorClient.DeletePost(r.Context(), &generatedCommon.UUIDMessage{Value: postData.Id.String()})
			if err != nil {
				h.logger.Error(err)
			}
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		_ = f.Close()
	}

	creatorInfo, err := h.creatorClient.CreatorNotificationInfo(r.Context(), &generatedCommon.UUIDMessage{Value: postData.Creator.String()})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if len(creatorInfo.Error) != 0 {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	notification := models.Notification{
		Topic: fmt.Sprintf("%s-%s", postData.Creator, "user"),
		Title: "Новый пост",
		Body:  fmt.Sprintf("У автора %s вышел новый пост \"%s\"", creatorInfo.Name, postData.Title),
		Photo: fmt.Sprintf("%s%s.jpg", models.PhotoURL, creatorInfo.Photo),
	}

	err = h.notificationApp.SendUserNotification(notification, r.Context())
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func (h *PostHandler) AddLike(w http.ResponseWriter, r *http.Request) {
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

	var like models.Like
	err = easyjson.UnmarshalFromReader(r.Body, &like)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	isPostOwner, err := h.creatorClient.IsPostOwner(r.Context(), &generatedCreator.PostUserMessage{
		UserID: (*userDataJWT).Id.String(),
		PostID: like.PostID.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if isPostOwner.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if isPostOwner.Flag {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	likeProto, err := h.creatorClient.AddLike(r.Context(), &generatedCreator.PostUserMessage{
		UserID: userDataJWT.Id.String(),
		PostID: like.PostID.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if likeProto.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if likeProto.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	like.LikesCount = likeProto.LikesCount
	utils.Response(w, http.StatusOK, like)
}

func (h *PostHandler) RemoveLike(w http.ResponseWriter, r *http.Request) {
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

	var like models.Like
	err = easyjson.UnmarshalFromReader(r.Body, &like)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	likeProto, err := h.creatorClient.RemoveLike(r.Context(), &generatedCreator.PostUserMessage{
		UserID: userDataJWT.Id.String(),
		PostID: like.PostID.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if likeProto.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if likeProto.Error == models.InternalError.Error() {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	like.LikesCount = likeProto.LikesCount
	utils.Response(w, http.StatusOK, like)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
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

	isPostOwner, err := h.creatorClient.IsPostOwner(r.Context(), &generatedCreator.PostUserMessage{
		UserID: (*userDataJWT).Id.String(),
		PostID: postID.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if isPostOwner.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if isPostOwner.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !isPostOwner.Flag {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	out, err := h.creatorClient.DeleteAttachmentsByPostID(r.Context(), &generatedCommon.UUIDMessage{Value: postID.String()})

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
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	out, err = h.creatorClient.DeletePost(r.Context(), &generatedCommon.UUIDMessage{Value: postID.String()})
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

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	var userID uuid.UUID
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)
	if err == nil {
		userID = userDataJWT.Id
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

	postProto, err := h.creatorClient.GetPost(r.Context(), &generatedCreator.PostUserMessage{
		UserID: userID.String(),
		PostID: postID.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if postProto.Error == models.Forbbiden.Error() {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if postProto.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if postProto.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	var post models.PostWithComments
	err = post.PostWithCommentsToModel(postProto)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	post.Post.Sanitize()

	utils.Response(w, http.StatusOK, post)
}

func (h *PostHandler) EditPost(w http.ResponseWriter, r *http.Request) {
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

	postEditData := models.PostEditData{}
	postEditData.Id, err = uuid.Parse(postIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	isPostOwner, err := h.creatorClient.IsPostOwner(r.Context(), &generatedCreator.PostUserMessage{
		UserID: (*userDataJWT).Id.String(),
		PostID: postEditData.Id.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if isPostOwner.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if isPostOwner.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !isPostOwner.Flag {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	if err = easyjson.UnmarshalFromReader(r.Body, &postEditData); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err = (models.PostCreationData{Title: postEditData.Title, Text: postEditData.Text}).IsValid(); err != nil {
		utils.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	var subs []string
	for _, v := range postEditData.AvailableSubscriptions {
		subs = append(subs, v.String())
	}
	out, err := h.creatorClient.EditPost(r.Context(), &generatedCreator.PostEditData{
		Id:                     postEditData.Id.String(),
		Title:                  postEditData.Title,
		Text:                   postEditData.Text,
		AvailableSubscriptions: subs,
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

func (h *PostHandler) AddAttach(w http.ResponseWriter, r *http.Request) {
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
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	isPostOwner, err := h.creatorClient.IsPostOwner(r.Context(), &generatedCreator.PostUserMessage{
		UserID: (*userDataJWT).Id.String(),
		PostID: postID.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if isPostOwner.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if isPostOwner.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !isPostOwner.Flag {
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

	attachmentType, err := h.creatorClient.GetFileExtension(r.Context(), &generatedCreator.KeywordMessage{Keyword: attach.Type})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if !attachmentType.Flag {
		utils.Response(w, http.StatusUnsupportedMediaType, nil)
		return
	}
	f, err := os.Create(fmt.Sprintf("%s.%s", filepath.Join(models.FolderPath, attach.Id.String()), attachmentType.Extension))
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if _, err := io.Copy(f, attach.Data); err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	_ = f.Close()

	out, err := h.creatorClient.AddAttach(r.Context(), &generatedCreator.PostAttachMessage{
		PostID: postID.String(),
		Attachment: &generatedCreator.Attachment{
			ID:   attach.Id.String(),
			Type: attach.Type,
		},
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error != "" {
		var attachProto = []*generatedCreator.Attachment{{ID: attach.Id.String(), Type: attach.Type}}
		_, err = h.creatorClient.DeleteAttachmentsFiles(r.Context(), &generatedCreator.Attachments{Attachments: attachProto})

		if err != nil {
			h.logger.Error(err)
		}

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
	isPostOwner, err := h.creatorClient.IsPostOwner(r.Context(), &generatedCreator.PostUserMessage{
		UserID: (*userDataJWT).Id.String(),
		PostID: postID.String(),
	})

	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if isPostOwner.Error == models.WrongData.Error() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if isPostOwner.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !isPostOwner.Flag {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	out, err := h.creatorClient.DeleteAttachment(r.Context(), &generatedCreator.PostAttachMessage{
		PostID: postID.String(),
		Attachment: &generatedCreator.Attachment{
			ID:   attachInfo.Id.String(),
			Type: attachInfo.Type,
		},
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
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}
