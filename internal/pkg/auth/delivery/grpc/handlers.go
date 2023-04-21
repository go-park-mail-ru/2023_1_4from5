package grpcAuth

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	"github.com/google/uuid"
	"time"
)

//go:generate mockgen -source=auth_grpc.pb.go -destination=auth_grpc.go -package=grpc

type GrpcAuthHandler struct {
	uc auth.AuthUsecase
	generated.AuthServiceServer
}

func NewGrpcAuthHandler(uc auth.AuthUsecase) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		uc: uc,
	}
}

func (h GrpcAuthHandler) SignIn(ctx context.Context, in *generated.LoginUser) (*generated.Token, error) {
	user := models.LoginUser{Login: in.Login, PasswordHash: in.PasswordHash}

	token, err := h.uc.SignIn(ctx, user)
	if err == nil {
		return &generated.Token{Cookie: token, Error: ""}, nil
	}
	return &generated.Token{Cookie: token, Error: err.Error()}, nil
}

func (h GrpcAuthHandler) IncUserVersion(ctx context.Context, in *generated.AccessDetails) (*generated.Empty, error) {
	idTmp, err := uuid.Parse(in.Id)
	if err != nil {
		return &generated.Empty{Error: models.WrongData.Error()}, nil
	}
	user := models.AccessDetails{
		Login:       in.Login,
		Id:          idTmp,
		UserVersion: int(in.UserVersion), //TODO: перегнать все int поля в int64
	}

	_, err = h.uc.IncUserVersion(ctx, user)
	if err == nil {
		return &generated.Empty{Error: ""}, nil
	}
	return &generated.Empty{Error: err.Error()}, nil
}

func (h GrpcAuthHandler) SignUp(ctx context.Context, in *generated.User) (*generated.Token, error) {
	idTmp, err := uuid.Parse(in.Id)
	idTmpPh, err := uuid.Parse(in.ProfilePhoto)
	if err != nil {
		return &generated.Token{Error: models.WrongData.Error()}, nil
	}
	user := models.User{
		Id:           idTmp,
		Login:        in.Login,
		Name:         in.Name,
		ProfilePhoto: idTmpPh,
		PasswordHash: in.PasswordHash,
		Registration: time.Time{},
		UserVersion:  0,
	}
	token, err := h.uc.SignUp(ctx, user)
	if err == nil {
		return &generated.Token{
			Cookie: token,
			Error:  "",
		}, nil
	}
	return &generated.Token{Error: err.Error()}, nil

}
