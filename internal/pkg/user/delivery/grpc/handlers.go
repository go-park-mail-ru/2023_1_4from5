package grpcUser

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
	generatedUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	"github.com/google/uuid"
)

//go:generate mockgen -source=./generated/user_grpc.pb.go -destination=../../mocks/user_grpc.go -package=mock

type GrpcUserHandler struct {
	uc user.UserUsecase

	generatedUser.UserServiceServer
}

func NewGrpcUserHandler(uc user.UserUsecase) *GrpcUserHandler {
	return &GrpcUserHandler{
		uc: uc,
	}
}

func (h GrpcUserHandler) Follow(ctx context.Context, in *generatedUser.FollowMessage) (*generatedCommon.Empty, error) {
	userId, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	creatorId, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.uc.Follow(ctx, userId, creatorId)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcUserHandler) Unfollow(ctx context.Context, in *generatedUser.FollowMessage) (*generatedCommon.Empty, error) {
	userId, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	creatorId, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.uc.Unfollow(ctx, userId, creatorId)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcUserHandler) AddPaymentInfo(ctx context.Context, in *generatedUser.SubscriptionDetails) (*generatedCommon.Empty, error) {
	userId, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	creatorId, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	subId, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	paymentInfo, err := uuid.Parse(in.PaymentInfo)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.uc.AddPaymentInfo(ctx, models.SubscriptionDetails{
		Id:          subId,
		UserID:      userId,
		CreatorId:   creatorId,
		MonthCount:  in.MonthCount,
		PaymentInfo: paymentInfo,
	})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcUserHandler) Subscribe(ctx context.Context, in *generatedUser.PaymentInfo) (*generatedUser.SubscriptionName, error) {
	paymentInfo, err := uuid.Parse(in.PaymentID)
	if err != nil {
		return &generatedUser.SubscriptionName{Error: err.Error()}, nil
	}

	subNotification, err := h.uc.Subscribe(ctx, paymentInfo, in.Money)
	if err != nil {
		return &generatedUser.SubscriptionName{Error: err.Error()}, nil
	}
	return &generatedUser.SubscriptionName{Error: "", Name: subNotification.SubscriptionName, CreatorID: subNotification.CreatorID.String()}, nil
}

func (h GrpcUserHandler) GetProfile(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedUser.UserProfile, error) {
	userId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedUser.UserProfile{Error: err.Error()}, nil
	}
	profile, err := h.uc.GetProfile(ctx, userId)
	if err != nil {
		return &generatedUser.UserProfile{Error: err.Error()}, nil
	}
	return &generatedUser.UserProfile{
		Login:        profile.Login,
		Name:         profile.Name,
		ProfilePhoto: profile.ProfilePhoto.String(),
		Registration: profile.Registration.String(),
		IsCreator:    profile.IsCreator,
		CreatorID:    profile.CreatorId.String(),
		Error:        ""}, nil
}

func (h GrpcUserHandler) UpdatePhoto(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedUser.ImageID, error) {
	userId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedUser.ImageID{Error: err.Error()}, nil
	}
	imageId, err := h.uc.UpdatePhoto(ctx, userId)
	if err != nil {
		return &generatedUser.ImageID{Error: err.Error()}, nil
	}
	return &generatedUser.ImageID{
		Value: imageId.String(),
		Error: ""}, nil
}

func (h GrpcUserHandler) DeletePhoto(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCommon.Empty, error) {
	imageId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}

	err = h.uc.DeletePhoto(ctx, imageId)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{
		Error: ""}, nil
}

func (h GrpcUserHandler) UpdatePassword(ctx context.Context, in *generatedUser.UpdatePasswordMessage) (*generatedCommon.Empty, error) {
	userId, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.uc.UpdatePassword(ctx, userId, in.Password)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcUserHandler) UpdateProfileInfo(ctx context.Context, in *generatedUser.UpdateProfileInfoMessage) (*generatedCommon.Empty, error) {
	userId, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.uc.UpdateProfileInfo(ctx, models.UpdateProfileInfo{
		Login: in.Login,
		Name:  in.Name}, userId)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}

func (h GrpcUserHandler) Donate(ctx context.Context, in *generatedUser.DonateMessage) (*generatedUser.DonateResponse, error) {
	creatorId, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedUser.DonateResponse{Error: err.Error()}, nil
	}
	money, err := h.uc.Donate(ctx, models.Donate{
		CreatorID:  creatorId,
		MoneyCount: in.MoneyCount})
	if err != nil {
		return &generatedUser.DonateResponse{Error: err.Error()}, nil
	}
	return &generatedUser.DonateResponse{MoneyCount: money, Error: ""}, nil
}

func (h GrpcUserHandler) BecomeCreator(ctx context.Context, in *generatedUser.BecameCreatorInfoMessage) (*generatedCommon.UUIDResponse, error) {
	userId, err := uuid.Parse(in.UserID)
	if err != nil {
		return &generatedCommon.UUIDResponse{Error: err.Error()}, nil
	}
	creatorId, err := h.uc.BecomeCreator(ctx, models.BecameCreatorInfo{
		Name:        in.Name,
		Description: in.Description}, userId)
	if err != nil {
		return &generatedCommon.UUIDResponse{Error: err.Error()}, nil
	}
	return &generatedCommon.UUIDResponse{Value: creatorId.String()}, nil
}

func (h GrpcUserHandler) UserSubscriptions(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedUser.SubscriptionsMessage, error) {
	userId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedUser.SubscriptionsMessage{Error: err.Error()}, nil
	}
	subs, err := h.uc.UserSubscriptions(ctx, userId)
	if err != nil {
		return &generatedUser.SubscriptionsMessage{Error: err.Error()}, nil
	}
	var subsProto generatedUser.SubscriptionsMessage
	for _, v := range subs {
		subsProto.Subscriptions = append(subsProto.Subscriptions, &generatedCommon.Subscription{
			Id:           v.Id.String(),
			Creator:      v.Creator.String(),
			CreatorName:  v.CreatorName,
			CreatorPhoto: v.CreatorPhoto.String(),
			MonthCost:    v.MonthCost,
			Title:        v.Title,
			Description:  v.Description,
		})
	}
	subsProto.Error = ""
	return &subsProto, nil
}

func (h GrpcUserHandler) CheckIfCreator(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedUser.CheckCreatorMessage, error) {
	userId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedUser.CheckCreatorMessage{Error: err.Error()}, nil
	}
	creatorId, flag, err := h.uc.CheckIfCreator(ctx, userId)
	if err != nil {
		return &generatedUser.CheckCreatorMessage{Error: err.Error()}, nil
	}
	return &generatedUser.CheckCreatorMessage{
		ID:        creatorId.String(),
		IsCreator: flag,
		Error:     ""}, nil
}

func (h GrpcUserHandler) UserFollows(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedUser.FollowsMessage, error) {
	userId, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedUser.FollowsMessage{Error: err.Error()}, nil
	}
	follows, err := h.uc.UserFollows(ctx, userId)
	if err != nil {
		return &generatedUser.FollowsMessage{Error: err.Error()}, nil
	}
	var followsProto generatedUser.FollowsMessage
	for _, v := range follows {
		followsProto.Follows = append(followsProto.Follows, &generatedUser.Follow{
			Creator:      v.Creator.String(),
			CreatorName:  v.CreatorName,
			CreatorPhoto: v.CreatorPhoto.String(),
			Description:  v.Description,
		})
	}
	followsProto.Error = ""
	return &followsProto, nil
}
