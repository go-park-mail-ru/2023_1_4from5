package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CreatorHandler struct {
	creatorClient generatedCreator.CreatorServiceClient
	authClient    generatedAuth.AuthServiceClient
	logger        *zap.SugaredLogger
}

func NewCreatorHandler(creatorClient generatedCreator.CreatorServiceClient, authClient generatedAuth.AuthServiceClient, logger *zap.SugaredLogger) *CreatorHandler {
	return &CreatorHandler{
		creatorClient: creatorClient,
		authClient:    authClient,
		logger:        logger,
	}
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
		postID, err := uuid.Parse(post.Id)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creatorID, err := uuid.Parse(post.CreatorID)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creatorPhoto, err := uuid.Parse(post.CreatorPhoto)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}

		reg, err := time.Parse("2006-01-02 15:04:05 -0700 -0700", post.Creation)

		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}

		feed[i] = models.Post{
			Id:            postID,
			Creator:       creatorID,
			CreatorPhoto:  creatorPhoto,
			CreatorName:   post.CreatorName,
			Creation:      reg,
			LikesCount:    post.LikesCount,
			Title:         post.Title,
			Text:          post.Text,
			IsAvailable:   post.IsAvailable,
			IsLiked:       post.IsLiked,
			Subscriptions: nil,
		}

		for _, attach := range post.PostAttachments {
			attachID, err := uuid.Parse(attach.ID)
			if err != nil {
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}
			feed[i].Attachments = append(feed[i].Attachments, models.Attachment{
				Id:   attachID,
				Type: attach.Type,
			})
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

	for i, v := range out.Creators {
		creatorId, err := uuid.Parse(v.Id)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		userId, err := uuid.Parse(v.UserID)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creatorPhoto, err := uuid.Parse(v.CreatorPhoto)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		coverPhoto, err := uuid.Parse(v.CoverPhoto)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creators[i] = models.Creator{
			Id:             creatorId,
			UserId:         userId,
			Name:           v.CreatorName,
			CoverPhoto:     coverPhoto,
			ProfilePhoto:   creatorPhoto,
			FollowersCount: v.FollowersCount,
			Description:    v.Description,
			PostsCount:     v.PostsCount,
		}

		creators[i].Sanitize()
	}
	utils.Response(w, http.StatusOK, creators)
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

	for i, v := range out.Creators {
		creatorId, err := uuid.Parse(v.Id)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		userId, err := uuid.Parse(v.UserID)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creatorPhoto, err := uuid.Parse(v.CreatorPhoto)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		coverPhoto, err := uuid.Parse(v.CoverPhoto)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creators[i] = models.Creator{
			Id:             creatorId,
			UserId:         userId,
			Name:           v.CreatorName,
			CoverPhoto:     coverPhoto,
			ProfilePhoto:   creatorPhoto,
			FollowersCount: v.FollowersCount,
			Description:    v.Description,
			PostsCount:     v.PostsCount,
		}

		creators[i].Sanitize()
	}

	utils.Response(w, http.StatusOK, creators)
}

func (h *CreatorHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	creatorUUID, ok := mux.Vars(r)["creator-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	userInfo := models.AccessDetails{}
	tmpUserInfo, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		tmpUserInfo = &userInfo
	}

	creatorID, err := uuid.Parse(creatorUUID)
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

	creatorPhoto, err := uuid.Parse(creatorPage.CreatorInfo.CreatorPhoto)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	coverPhoto, err := uuid.Parse(creatorPage.CreatorInfo.CoverPhoto)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	for _, sub := range creatorPage.Subscriptions {
		subID, err := uuid.Parse(sub.Id)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		page.Subscriptions = append(page.Subscriptions, models.Subscription{
			Id:           subID,
			Creator:      creatorID,
			CreatorName:  creatorPage.CreatorInfo.CreatorName,
			CreatorPhoto: creatorPhoto,
			MonthCost:    sub.MonthCost,
			Title:        sub.Title,
			Description:  sub.Description,
		})
	}

	page.IsMyPage = creatorPage.IsMyPage
	page.Follows = creatorPage.Follows
	page.Aim = models.Aim{
		Creator:     creatorID,
		Description: creatorPage.AimInfo.Description,
		MoneyNeeded: creatorPage.AimInfo.MoneyNeeded,
		MoneyGot:    creatorPage.AimInfo.MoneyGot,
	}
	page.CreatorInfo = models.Creator{
		Id:             creatorID,
		UserId:         userInfo.Id,
		Name:           creatorPage.CreatorInfo.CreatorName,
		CoverPhoto:     coverPhoto,
		ProfilePhoto:   creatorPhoto,
		FollowersCount: creatorPage.CreatorInfo.FollowersCount,
		Description:    creatorPage.CreatorInfo.Description,
		PostsCount:     creatorPage.CreatorInfo.PostsCount,
	}
	page.Posts = make([]models.Post, len(creatorPage.Posts))
	for i, post := range creatorPage.Posts {

		postID, err := uuid.Parse(post.Id)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}

		reg, err := time.Parse("2006-01-02 15:04:05 -0700 -0700", post.Creation)

		if err != nil {
			h.logger.Error(err)
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}

		page.Posts[i] = models.Post{
			Id:           postID,
			Creator:      creatorID,
			CreatorPhoto: creatorPhoto,
			CreatorName:  post.CreatorName,
			Creation:     reg,
			LikesCount:   post.LikesCount,
			Title:        post.Title,
			Text:         post.Text,
			IsAvailable:  post.IsAvailable,
			IsLiked:      post.IsLiked,
		}

		for _, attach := range post.PostAttachments {
			attachID, err := uuid.Parse(attach.ID)
			if err != nil {
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}
			page.Posts[i].Attachments = append(page.Posts[i].Attachments, models.Attachment{
				Id:   attachID,
				Type: attach.Type,
			})
		}

		for _, sub := range post.Subscriptions {
			subID, err := uuid.Parse(sub.Id)
			if err != nil {
				utils.Response(w, http.StatusInternalServerError, nil)
				return
			}
			page.Posts[i].Subscriptions = append(page.Posts[i].Subscriptions, models.Subscription{
				Id:           subID,
				Creator:      creatorID,
				CreatorName:  creatorPage.CreatorInfo.CreatorName,
				CreatorPhoto: creatorPhoto,
				MonthCost:    sub.MonthCost,
				Title:        sub.Title,
				Description:  sub.Description,
			})
		}
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
	aimInfo := models.Aim{}
	err = easyjson.UnmarshalFromReader(r.Body, &aimInfo)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if len(aimInfo.Description) > 100 {
		utils.Response(w, http.StatusBadRequest, nil)
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
	if err != nil || !(updCreator.IsValid()) {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	out, err := h.creatorClient.UpdateCreatorData(r.Context(), &generatedCreator.UpdateCreatorInfo{
		CreatorName: updCreator.Name,
		Description: updCreator.Description,
		CreatorID:   creatorID.String(),
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
