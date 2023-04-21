package grpcAuth

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
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

func (h GrpcAuthHandler) Login(ctx context.Context, in *generated.LoginUser) (*generated.Token, error) {

	user := models.LoginUser{Login: in.Login, PasswordHash: in.PasswordHash}

	token, err := h.uc.SignIn(ctx, user)

	return &generated.Token{Cookie: token, Error: err.Error()}, nil
}

//TODO:функционал
