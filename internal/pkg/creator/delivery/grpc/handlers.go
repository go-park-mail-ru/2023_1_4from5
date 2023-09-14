package grpcCreator

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	comment "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/comment"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/subscription"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"time"
)

//go:generate mockgen -source=./generated/creator_grpc.pb.go -destination=../../mocks/creator_grpc.go -package=mock

type GrpcCreatorHandler struct {
	uc  creator.CreatorUsecase
	puc post.PostUsecase
	auc attachment.AttachmentUsecase
	suc subscription.SubscriptionUsecase
	cuc comment.CommentUsecase
	generatedCreator.CreatorServiceServer
}

func NewGrpcCreatorHandler(uc creator.CreatorUsecase, puc post.PostUsecase, auc attachment.AttachmentUsecase, suc subscription.SubscriptionUsecase, cuc comment.CommentUsecase) *GrpcCreatorHandler {
	return &GrpcCreatorHandler{
		uc:  uc,
		puc: puc,
		auc: auc,
		suc: suc,
		cuc: cuc,
	}
}

func (h GrpcCreatorHandler) FindCreators(ctx context.Context, in *generatedCreator.KeywordMessage) (*generatedCreator.CreatorsMessage, error) {
	creators, err := h.uc.FindCreators(ctx, in.Keyword)
	if err != nil {
		return &generatedCreator.CreatorsMessage{Error: err.Error()}, nil
	}
	var creatorsMessage generatedCreator.CreatorsMessage
	for _, v := range creators {
		creatorsMessage.Creators = append(creatorsMessage.Creators, &generatedCreator.Creator{
			Id:             v.Id.String(),
			UserID:         v.UserId.String(),
			CreatorName:    v.Name,
			CreatorPhoto:   v.ProfilePhoto.String(),
			CoverPhoto:     v.CoverPhoto.String(),
			FollowersCount: v.FollowersCount,
			Description:    v.Description,
			PostsCount:     v.PostsCount,
		})
	}
	creatorsMessage.Error = ""

	return &creatorsMessage, nil
}

func (h GrpcCreatorHandler) GetAllCreators(ctx context.Context, in *generatedCommon.Empty) (*generatedCreator.CreatorsMessage, error) {
	creators, err := h.uc.GetAllCreators(ctx)
	if err != nil {
		return &generatedCreator.CreatorsMessage{Error: err.Error()}, nil
	}
	var creatorsMessage generatedCreator.CreatorsMessage
	for _, v := range creators {
		creatorsMessage.Creators = append(creatorsMessage.Creators, &generatedCreator.Creator{
			Id:             v.Id.String(),
			UserID:         v.UserId.String(),
			CreatorName:    v.Name,
			CreatorPhoto:   v.ProfilePhoto.String(),
			CoverPhoto:     v.CoverPhoto.String(),
			FollowersCount: v.FollowersCount,
			Description:    v.Description,
			PostsCount:     v.PostsCount,
		})
	}
	creatorsMessage.Error = ""

	return &creatorsMessage, nil
}

func (h GrpcCreatorHandler) GetFeed(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCreator.PostsMessage, error) {
	userID, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCreator.PostsMessage{Error: err.Error()}, nil
	}
	feed, err := h.uc.GetFeed(ctx, userID)
	if err != nil {
		return &generatedCreator.PostsMessage{Error: err.Error()}, nil
	}

	var postsProto generatedCreator.PostsMessage
	for i, post := range feed {
		postsProto.Posts = append(postsProto.Posts, &generatedCreator.Post{
			Id:            post.Id.String(),
			CreatorID:     post.Creator.String(),
			Creation:      post.Creation.Format(time.RFC3339),
			CreatorName:   post.CreatorName,
			LikesCount:    post.LikesCount,
			CommentsCount: post.CommentsCount,
			CreatorPhoto:  post.CreatorPhoto.String(),
			Title:         post.Title,
			Text:          post.Text,
			IsAvailable:   true,
			IsLiked:       post.IsLiked,
		})

		for _, attach := range post.Attachments {
			postsProto.Posts[i].PostAttachments = append(postsProto.Posts[i].PostAttachments, &generatedCreator.Attachment{
				ID:   attach.Id.String(),
				Type: attach.Type,
			})
		}
		postsProto.Posts[i].Subscriptions = nil

	}
	postsProto.Error = ""

	return &postsProto, nil
}

func (h GrpcCreatorHandler) GetPage(ctx context.Context, in *generatedCreator.UserCreatorMessage) (*generatedCreator.CreatorPage, error) {
	creatorID, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCreator.CreatorPage{Error: err.Error()}, nil
	}
	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCreator.CreatorPage{Error: err.Error()}, nil
	}

	page, err := h.uc.GetPage(ctx, userID, creatorID)

	if err != nil {
		return &generatedCreator.CreatorPage{Error: err.Error()}, nil
	}

	var creatorPage generatedCreator.CreatorPage
	creatorPage.AimInfo = &generatedCreator.Aim{
		Creator:     page.Aim.Creator.String(),
		Description: page.Aim.Description,
		MoneyNeeded: page.Aim.MoneyNeeded,
		MoneyGot:    page.Aim.MoneyGot,
	}
	creatorPage.Error = ""
	creatorPage.IsMyPage = page.IsMyPage
	creatorPage.Follows = page.Follows
	for _, sub := range page.Subscriptions {
		creatorPage.Subscriptions = append(creatorPage.Subscriptions, &generatedCommon.Subscription{
			Id:           sub.Id.String(),
			Creator:      sub.Creator.String(),
			CreatorName:  sub.CreatorName,
			CreatorPhoto: sub.CreatorPhoto.String(),
			MonthCost:    sub.MonthCost,
			Title:        sub.Title,
			Description:  sub.Description,
		})
	}

	creatorPage.CreatorInfo = &generatedCreator.Creator{
		Id:             page.CreatorInfo.Id.String(),
		UserID:         page.CreatorInfo.UserId.String(),
		CreatorName:    page.CreatorInfo.Name,
		CreatorPhoto:   page.CreatorInfo.ProfilePhoto.String(),
		CoverPhoto:     page.CreatorInfo.CoverPhoto.String(),
		FollowersCount: page.CreatorInfo.FollowersCount,
		Description:    page.CreatorInfo.Description,
		PostsCount:     page.CreatorInfo.PostsCount,
	}
	for i, post := range page.Posts {
		creatorPage.Posts = append(creatorPage.Posts, &generatedCreator.Post{
			Id:            post.Id.String(),
			CreatorID:     post.Creator.String(),
			Creation:      post.Creation.Format(time.RFC3339),
			CreatorName:   post.CreatorName,
			LikesCount:    post.LikesCount,
			CommentsCount: post.CommentsCount,
			CreatorPhoto:  post.CreatorPhoto.String(),
			Title:         post.Title,
			Text:          post.Text,
			IsAvailable:   post.IsAvailable,
			IsLiked:       post.IsLiked,
		})

		for _, attach := range post.Attachments {
			creatorPage.Posts[i].PostAttachments = append(creatorPage.Posts[i].PostAttachments, &generatedCreator.Attachment{
				ID:   attach.Id.String(),
				Type: attach.Type,
			})
		}

		for _, sub := range post.Subscriptions {
			creatorPage.Posts[i].Subscriptions = append(creatorPage.Posts[i].Subscriptions, &generatedCommon.Subscription{
				Id:           sub.Id.String(),
				Creator:      sub.Creator.String(),
				CreatorName:  sub.CreatorName,
				CreatorPhoto: sub.CreatorPhoto.String(),
				MonthCost:    sub.MonthCost,
				Title:        sub.Title,
				Description:  sub.Description,
			})
		}
	}

	return &creatorPage, nil
}

func (h GrpcCreatorHandler) UpdateCreatorData(ctx context.Context, in *generatedCreator.UpdateCreatorInfo) (*generatedCommon.Empty, error) {
	creatorID, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.uc.UpdateCreatorData(ctx, models.UpdateCreatorInfo{
		Description: in.Description,
		CreatorName: in.CreatorName,
		CreatorID:   creatorID})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) IsCreator(ctx context.Context, in *generatedCreator.UserCreatorMessage) (*generatedCreator.FlagMessage, error) {
	creatorID, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCreator.FlagMessage{Error: err.Error()}, nil
	}
	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCreator.FlagMessage{Error: err.Error()}, nil
	}

	isCreator, err := h.puc.IsCreator(ctx, userID, creatorID)
	if err != nil {
		return &generatedCreator.FlagMessage{Error: err.Error()}, nil
	}
	return &generatedCreator.FlagMessage{Error: "", Flag: isCreator}, nil
}

func (h GrpcCreatorHandler) CreateAim(ctx context.Context, in *generatedCreator.Aim) (*generatedCommon.Empty, error) {
	creatorID, err := uuid.Parse(in.Creator)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	err = h.uc.CreateAim(ctx, models.Aim{
		Creator:     creatorID,
		Description: in.Description,
		MoneyNeeded: in.MoneyNeeded,
		MoneyGot:    in.MoneyGot,
	})

	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) CheckIfCreator(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCommon.UUIDResponse, error) {
	userID, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCommon.UUIDResponse{Error: err.Error()}, nil
	}

	creatorID, err := h.uc.CheckIfCreator(ctx, userID)
	if err != nil {
		return &generatedCommon.UUIDResponse{Error: err.Error(), Value: creatorID.String()}, nil
	}
	return &generatedCommon.UUIDResponse{Error: "", Value: creatorID.String()}, nil
}

func (h GrpcCreatorHandler) CreatorNotificationInfo(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCreator.NotificationCreatorInfo, error) {
	creatorID, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCreator.NotificationCreatorInfo{Error: err.Error()}, nil
	}

	info, err := h.uc.CreatorNotificationInfo(ctx, creatorID)
	if err != nil {
		return &generatedCreator.NotificationCreatorInfo{Error: err.Error()}, nil
	}
	return &generatedCreator.NotificationCreatorInfo{Error: "", Name: info.Name, Photo: info.Photo.String()}, nil
}

func (h GrpcCreatorHandler) CreatePost(ctx context.Context, in *generatedCreator.PostCreationData) (*generatedCommon.Empty, error) {
	ID, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	creatorID, err := uuid.Parse(in.Creator)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	var attachs []models.AttachmentData
	for _, attach := range in.Attachments {
		attachID, err := uuid.Parse(attach.ID)
		if err != nil {
			return &generatedCommon.Empty{Error: err.Error()}, nil
		}
		attachs = append(attachs, models.AttachmentData{
			Id:   attachID,
			Data: nil,
			Type: attach.Type,
		})
	}

	var subs []uuid.UUID
	for _, sub := range in.AvailableSubscriptions {
		subID, err := uuid.Parse(sub)
		if err != nil {
			return &generatedCommon.Empty{Error: err.Error()}, nil
		}
		subs = append(subs, subID)
	}

	err = h.puc.CreatePost(ctx, models.PostCreationData{
		Id:                     ID,
		Creator:                creatorID,
		Title:                  in.Title,
		Text:                   in.Text,
		Attachments:            attachs,
		AvailableSubscriptions: subs,
	})

	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) DeleteAttachmentsFiles(ctx context.Context, in *generatedCreator.Attachments) (*generatedCommon.Empty, error) {
	var attachs []models.Attachment
	for _, attach := range in.Attachments {
		attachID, err := uuid.Parse(attach.ID)
		if err != nil {
			return &generatedCommon.Empty{Error: err.Error()}, nil
		}

		attachs = append(attachs, models.Attachment{
			Id:   attachID,
			Type: attach.Type,
		})
	}

	err := h.auc.DeleteAttachmentsFiles(ctx, attachs...)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) DeleteAttachmentsByPostID(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCommon.Empty, error) {
	postID, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	err = h.auc.DeleteAttachmentsByPostID(ctx, postID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) StatisticsFirstDate(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCreator.FirstDate, error) {
	creatorID, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCreator.FirstDate{Error: err.Error()}, nil
	}
	firstDate, err := h.uc.StatisticsFirstDate(ctx, creatorID)
	if err != nil {
		return &generatedCreator.FirstDate{Error: err.Error()}, nil
	}
	return &generatedCreator.FirstDate{Error: "", Date: firstDate}, nil
}

func (h GrpcCreatorHandler) DeleteAttachment(ctx context.Context, in *generatedCreator.PostAttachMessage) (*generatedCommon.Empty, error) {
	postID, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	attachID, err := uuid.Parse(in.Attachment.ID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	err = h.auc.DeleteAttachment(ctx, postID, models.Attachment{
		Id:   attachID,
		Type: in.Attachment.Type,
	})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) AddAttach(ctx context.Context, in *generatedCreator.PostAttachMessage) (*generatedCommon.Empty, error) {
	postID, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	attachID, err := uuid.Parse(in.Attachment.ID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	err = h.auc.AddAttach(ctx, postID, models.Attachment{
		Id:   attachID,
		Type: in.Attachment.Type,
	})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) GetFileExtension(ctx context.Context, in *generatedCreator.KeywordMessage) (*generatedCreator.Extension, error) {
	extension, flag := h.auc.GetFileExtension(ctx, in.Keyword)

	return &generatedCreator.Extension{
		Extension: extension,
		Flag:      flag,
	}, nil
}

func (h GrpcCreatorHandler) IsPostOwner(ctx context.Context, in *generatedCreator.PostUserMessage) (*generatedCreator.FlagMessage, error) {
	postID, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCreator.FlagMessage{Error: err.Error()}, nil
	}

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCreator.FlagMessage{Error: err.Error()}, nil
	}

	flag, err := h.puc.IsPostOwner(ctx, userID, postID)
	if err != nil {
		return &generatedCreator.FlagMessage{Error: err.Error()}, nil
	}
	return &generatedCreator.FlagMessage{Flag: flag, Error: ""}, nil
}

func (h GrpcCreatorHandler) IsPostAvailable(ctx context.Context, in *generatedCreator.PostUserMessage) (*generatedCommon.Empty, error) {
	postID, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	err = h.puc.IsPostAvailable(ctx, userID, postID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) DeletePost(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCommon.Empty, error) {
	postID, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	err = h.puc.DeletePost(ctx, postID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) GetCreatorBalance(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCreator.CreatorBalance, error) {
	creatorID, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCreator.CreatorBalance{Error: err.Error()}, nil
	}

	balance, err := h.uc.GetCreatorBalance(ctx, creatorID)
	if err != nil {
		return &generatedCreator.CreatorBalance{Error: err.Error()}, nil
	}
	return &generatedCreator.CreatorBalance{Error: "", Balance: balance}, nil
}

func (h GrpcCreatorHandler) UpdateBalance(ctx context.Context, in *generatedCreator.CreatorTransfer) (*generatedCreator.CreatorBalance, error) {
	creatorID, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCreator.CreatorBalance{Error: err.Error()}, nil
	}

	balance, err := h.uc.UpdateBalance(ctx, models.CreatorTransfer{
		Money:       in.Money,
		CreatorID:   creatorID,
		PhoneNumber: "",
	})
	if err != nil {
		return &generatedCreator.CreatorBalance{Error: err.Error()}, nil
	}
	return &generatedCreator.CreatorBalance{Error: "", Balance: balance}, nil
}

func (h GrpcCreatorHandler) GetPost(ctx context.Context, in *generatedCreator.PostUserMessage) (*generatedCreator.PostWithComments, error) {
	postID, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCreator.PostWithComments{Error: err.Error()}, nil
	}

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCreator.PostWithComments{Error: err.Error()}, nil
	}

	post, err := h.puc.GetPost(ctx, postID, userID)
	if err != nil {
		return &generatedCreator.PostWithComments{Error: err.Error()}, nil
	}

	var subs []*generatedCommon.Subscription

	for _, v := range post.Post.Subscriptions {
		subs = append(subs, &generatedCommon.Subscription{
			Id:           v.Id.String(),
			Creator:      v.Creator.String(),
			CreatorName:  v.CreatorName,
			CreatorPhoto: v.CreatorPhoto.String(),
			MonthCost:    v.MonthCost,
			Title:        v.Title,
			Description:  v.Description,
		})
	}

	var attachs []*generatedCreator.Attachment

	for _, v := range post.Post.Attachments {
		attachs = append(attachs, &generatedCreator.Attachment{ID: v.Id.String(), Type: v.Type})
	}

	var comments []*generatedCreator.Comment

	for _, v := range post.Comments {
		comments = append(comments, &generatedCreator.Comment{
			Id:         v.CommentID.String(),
			UserId:     v.UserID.String(),
			Username:   v.Username,
			UserPhoto:  v.UserPhoto.String(),
			PostID:     v.PostID.String(),
			Text:       v.Text,
			Creation:   v.Creation.Format(time.RFC3339),
			LikesCount: v.LikesCount,
			IsOwner:    v.IsOwner,
			IsLiked:    v.IsLiked,
		})
	}

	return &generatedCreator.PostWithComments{Error: "", Post: &generatedCreator.Post{
		Id:              post.Post.Id.String(),
		CreatorID:       post.Post.Creator.String(),
		Creation:        post.Post.Creation.Format(time.RFC3339),
		CreatorName:     post.Post.CreatorName,
		LikesCount:      post.Post.LikesCount,
		CreatorPhoto:    post.Post.CreatorPhoto.String(),
		Title:           post.Post.Title,
		Text:            post.Post.Text,
		IsAvailable:     post.Post.IsAvailable,
		PostAttachments: attachs,
		IsLiked:         post.Post.IsLiked,
		Subscriptions:   subs,
	}, Comments: comments}, nil
}

func (h GrpcCreatorHandler) EditPost(ctx context.Context, in *generatedCreator.PostEditData) (*generatedCommon.Empty, error) {
	postID, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	var subs []uuid.UUID
	for _, sub := range in.AvailableSubscriptions {
		subID, err := uuid.Parse(sub)
		if err != nil {
			return &generatedCommon.Empty{Error: err.Error()}, nil
		}
		subs = append(subs, subID)
	}
	err = h.puc.EditPost(ctx, models.PostEditData{
		Id:                     postID,
		Title:                  in.Title,
		Text:                   in.Text,
		AvailableSubscriptions: subs,
	})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) RemoveLike(ctx context.Context, in *generatedCreator.PostUserMessage) (*generatedCreator.Like, error) {
	postID, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCreator.Like{Error: err.Error()}, nil
	}

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCreator.Like{Error: err.Error()}, nil
	}
	like, err := h.puc.RemoveLike(ctx, userID, postID)
	if err != nil {
		return &generatedCreator.Like{Error: err.Error()}, nil
	}

	return &generatedCreator.Like{
		LikesCount: like.LikesCount,
		PostID:     like.PostID.String(),
		Error:      "",
	}, nil
}

func (h GrpcCreatorHandler) AddLike(ctx context.Context, in *generatedCreator.PostUserMessage) (*generatedCreator.Like, error) {
	postID, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCreator.Like{Error: err.Error()}, nil
	}

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCreator.Like{Error: err.Error()}, nil
	}
	like, err := h.puc.AddLike(ctx, userID, postID)
	if err != nil {
		return &generatedCreator.Like{Error: err.Error()}, nil
	}

	return &generatedCreator.Like{
		LikesCount: like.LikesCount,
		PostID:     like.PostID.String(),
		Error:      "",
	}, nil
}

func (h GrpcCreatorHandler) UpdateProfilePhoto(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCommon.UUIDResponse, error) {
	creatorId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCommon.UUIDResponse{Error: err.Error()}, nil
	}
	imageId, err := h.uc.UpdateProfilePhoto(ctx, creatorId)
	if err != nil {
		return &generatedCommon.UUIDResponse{Error: err.Error()}, nil
	}
	return &generatedCommon.UUIDResponse{
		Value: imageId.String(),
		Error: ""}, nil
}

func (h GrpcCreatorHandler) DeleteProfilePhoto(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCommon.Empty, error) {
	creatorId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.uc.DeleteProfilePhoto(ctx, creatorId)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{
		Error: ""}, nil
}

func (h GrpcCreatorHandler) UpdateCoverPhoto(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCommon.UUIDResponse, error) {
	creatorId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCommon.UUIDResponse{Error: err.Error()}, nil
	}
	imageId, err := h.uc.UpdateCoverPhoto(ctx, creatorId)
	if err != nil {
		return &generatedCommon.UUIDResponse{Error: err.Error()}, nil
	}
	return &generatedCommon.UUIDResponse{
		Value: imageId.String(),
		Error: ""}, nil
}

func (h GrpcCreatorHandler) DeleteCoverPhoto(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCommon.Empty, error) {
	creatorId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	err = h.uc.DeleteCoverPhoto(ctx, creatorId)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{
		Error: ""}, nil
}

func (h GrpcCreatorHandler) CreateSubscription(ctx context.Context, in *generatedCommon.Subscription) (*generatedCommon.Empty, error) {
	subId, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	creatorId, err := uuid.Parse(in.Creator)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.suc.CreateSubscription(ctx, models.Subscription{
		Id:          subId,
		Creator:     creatorId,
		CreatorName: in.CreatorName,
		MonthCost:   in.MonthCost,
		Title:       in.Title,
		Description: in.Description,
	})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) EditSubscription(ctx context.Context, in *generatedCommon.Subscription) (*generatedCommon.Empty, error) {
	subId, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	creatorId, err := uuid.Parse(in.Creator)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.suc.EditSubscription(ctx, models.Subscription{
		Id:          subId,
		Creator:     creatorId,
		MonthCost:   in.MonthCost,
		Title:       in.Title,
		Description: in.Description,
	})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) DeleteSubscription(ctx context.Context, in *generatedCreator.SubscriptionCreatorMessage) (*generatedCommon.Empty, error) {
	subId, err := uuid.Parse(in.SubscriptionID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	creatorId, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.suc.DeleteSubscription(ctx, subId, creatorId)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) CreateComment(ctx context.Context, in *generatedCreator.Comment) (*generatedCommon.Empty, error) {
	commentId, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCommon.Empty{Error: models.WrongData.Error()}, nil
	}
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		return &generatedCommon.Empty{Error: models.WrongData.Error()}, nil
	}
	postId, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCommon.Empty{Error: models.WrongData.Error()}, nil
	}
	err = h.cuc.CreateComment(ctx, models.Comment{
		CommentID: commentId,
		UserID:    userId,
		PostID:    postId,
		Text:      in.Text,
	})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) DeleteComment(ctx context.Context, in *generatedCreator.Comment) (*generatedCommon.Empty, error) {
	commentId, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCommon.Empty{Error: models.WrongData.Error()}, nil
	}
	postId, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCommon.Empty{Error: models.WrongData.Error()}, nil
	}
	err = h.cuc.DeleteComment(ctx, models.Comment{
		CommentID: commentId,
		PostID:    postId,
	})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) EditComment(ctx context.Context, in *generatedCreator.Comment) (*generatedCommon.Empty, error) {
	commentId, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCommon.Empty{Error: models.WrongData.Error()}, nil
	}

	err = h.cuc.EditComment(ctx, models.Comment{
		CommentID: commentId,
		Text:      in.Text,
	})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcCreatorHandler) AddLikeComment(ctx context.Context, in *generatedCreator.Comment) (*generatedCreator.Like, error) {
	commentId, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCreator.Like{Error: models.WrongData.Error()}, nil
	}
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		return &generatedCreator.Like{Error: models.WrongData.Error()}, nil
	}
	postId, err := uuid.Parse(in.PostID)
	if err != nil {
		return &generatedCreator.Like{Error: models.WrongData.Error()}, nil
	}

	likesCount, err := h.cuc.AddLike(ctx, models.Comment{
		CommentID: commentId,
		UserID:    userId,
		PostID:    postId,
	})

	if err != nil {
		return &generatedCreator.Like{Error: err.Error()}, nil
	}
	return &generatedCreator.Like{Error: "", LikesCount: likesCount}, nil
}

func (h GrpcCreatorHandler) RemoveLikeComment(ctx context.Context, in *generatedCreator.Comment) (*generatedCreator.Like, error) {
	commentId, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCreator.Like{Error: models.WrongData.Error()}, nil
	}
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		return &generatedCreator.Like{Error: models.WrongData.Error()}, nil
	}

	likesCount, err := h.cuc.RemoveLike(ctx, models.Comment{
		CommentID: commentId,
		UserID:    userId,
	})

	if err != nil {
		return &generatedCreator.Like{Error: err.Error()}, nil
	}
	return &generatedCreator.Like{Error: "", LikesCount: likesCount}, nil
}

func (h GrpcCreatorHandler) IsCommentOwner(ctx context.Context, in *generatedCreator.Comment) (*generatedCreator.FlagMessage, error) {
	commentId, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCreator.FlagMessage{
			Flag:  false,
			Error: models.WrongData.Error(),
		}, nil
	}
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		return &generatedCreator.FlagMessage{
			Flag:  false,
			Error: models.WrongData.Error(),
		}, nil
	}

	flag, err := h.cuc.IsCommentOwner(ctx, models.Comment{
		CommentID: commentId,
		UserID:    userId,
	})
	if err != nil {
		return &generatedCreator.FlagMessage{
			Flag:  flag,
			Error: err.Error(),
		}, nil
	}
	return &generatedCreator.FlagMessage{
		Flag:  flag,
		Error: "",
	}, nil
}

func (h GrpcCreatorHandler) Statistics(ctx context.Context, in *generatedCreator.StatisticsInput) (*generatedCreator.Stat, error) {
	creatorId, err := uuid.Parse(in.CreatorId)
	if err != nil {
		return &generatedCreator.Stat{Error: models.WrongData.Error()}, nil
	}

	firstDate, err := time.Parse(time.RFC3339, in.FirstDate)
	if err != nil {
		return &generatedCreator.Stat{Error: models.WrongData.Error()}, nil
	}
	secondDate, err := time.Parse(time.RFC3339, in.SecondDate)
	if err != nil {
		return &generatedCreator.Stat{Error: models.WrongData.Error()}, nil
	}

	stat, err := h.uc.Statistics(ctx, models.StatisticsDates{
		CreatorId:   creatorId,
		FirstMonth:  firstDate,
		SecondMonth: secondDate,
	})

	if err != nil {
		return &generatedCreator.Stat{Error: err.Error()}, nil
	}
	return &generatedCreator.Stat{
		CreatorId:              stat.CreatorId.String(),
		PostsPerMonth:          stat.PostsPerMonth,
		SubscriptionsBought:    stat.SubscriptionsBought,
		DonationsCount:         stat.DonationsCount,
		MoneyFromDonations:     stat.MoneyFromDonations,
		MoneyFromSubscriptions: stat.MoneyFromSubscriptions,
		NewFollowers:           stat.NewFollowers,
		LikesCount:             stat.LikesCount,
		Error:                  "",
	}, nil
}
