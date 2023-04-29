package grpcCreator

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type GrpcCreatorHandler struct {
	uc  creator.CreatorUsecase
	puc post.PostUsecase
	generatedCreator.CreatorServiceServer
}

func NewGrpcCreatorHandler(uc creator.CreatorUsecase, puc post.PostUsecase) *GrpcCreatorHandler {
	return &GrpcCreatorHandler{
		uc:  uc,
		puc: puc,
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
			Id:           post.Id.String(),
			CreatorID:    post.Creator.String(),
			Creation:     post.Creation.String(),
			CreatorName:  post.CreatorName,
			LikesCount:   post.LikesCount,
			CreatorPhoto: post.CreatorPhoto.String(),
			Title:        post.Title,
			Text:         post.Text,
			IsAvailable:  true,
			IsLiked:      post.IsLiked,
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
	fmt.Println(page.IsMyPage)

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
			Id:           post.Id.String(),
			CreatorID:    post.Creator.String(),
			Creation:     post.Creation.String(),
			CreatorName:  post.CreatorName,
			LikesCount:   post.LikesCount,
			CreatorPhoto: post.CreatorPhoto.String(),
			Title:        post.Title,
			Text:         post.Text,
			IsAvailable:  post.IsAvailable,
			IsLiked:      post.IsLiked,
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
