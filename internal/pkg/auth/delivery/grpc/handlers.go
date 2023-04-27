package grpcAuth

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	"github.com/google/uuid"
	"time"
)

//go:generate mockgen -source=auth_grpc.pb.go -destination=auth_grpc.go -package=grpc

type GrpcAuthHandler struct {
	uc auth.AuthUsecase

	generatedAuth.AuthServiceServer
}

func NewGrpcAuthHandler(uc auth.AuthUsecase) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		uc: uc,
	}
}

func (h GrpcAuthHandler) SignIn(ctx context.Context, in *generatedAuth.LoginUser) (*generatedAuth.Token, error) {
	user := models.LoginUser{Login: in.Login, PasswordHash: in.PasswordHash}

	token, err := h.uc.SignIn(ctx, user)
	if err == nil {
		return &generatedAuth.Token{Cookie: token, Error: ""}, nil
	}
	return &generatedAuth.Token{Cookie: token, Error: err.Error()}, nil
}

func (h GrpcAuthHandler) IncUserVersion(ctx context.Context, in *generatedAuth.AccessDetails) (*generatedCommon.Empty, error) {
	idTmp, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedCommon.Empty{Error: models.WrongData.Error()}, nil
	}
	user := models.AccessDetails{
		Login:       in.Login,
		Id:          idTmp,
		UserVersion: int(in.UserVersion), //TODO: перегнать все int поля в int64
	}

	_, err = h.uc.IncUserVersion(ctx, user)
	if err == nil {
		return &generatedCommon.Empty{Error: ""}, nil
	}
	return &generatedCommon.Empty{Error: err.Error()}, nil
}

func (h GrpcAuthHandler) SignUp(ctx context.Context, in *generatedAuth.User) (*generatedAuth.Token, error) {
	idTmp, err := uuid.Parse(in.Id)
	idTmpPh, err := uuid.Parse(in.ProfilePhoto)
	if err != nil {
		return &generatedAuth.Token{Error: models.WrongData.Error()}, nil
	}
	user := models.User{
		Id:           idTmp,
		Login:        in.Login, //TODO: наверное здесь можно заполнять не все поля
		Name:         in.Name,
		ProfilePhoto: idTmpPh,
		PasswordHash: in.PasswordHash,
		Registration: time.Time{},
		UserVersion:  0,
	}
	token, err := h.uc.SignUp(ctx, user)
	if err == nil {
		return &generatedAuth.Token{
			Cookie: token,
			Error:  "",
		}, nil
	}
	return &generatedAuth.Token{Error: err.Error()}, nil

}

func (h GrpcAuthHandler) CheckUserVersion(ctx context.Context, in *generatedAuth.AccessDetails) (*generatedAuth.UserVersion, error) {
	idTmp, err := uuid.Parse(in.Id)
	if err != nil {
		return &generatedAuth.UserVersion{Error: err.Error()}, nil
	}
	user := models.AccessDetails{
		Login:       in.Login,
		Id:          idTmp,
		UserVersion: int(in.UserVersion),
	}

	uv, err := h.uc.CheckUserVersion(ctx, user)
	if err == nil {
		return &generatedAuth.UserVersion{
			UserVersion: int64(uv),
			Error:       "",
		}, nil
	}
	return &generatedAuth.UserVersion{Error: err.Error()}, nil
}

func (h GrpcAuthHandler) CheckUser(ctx context.Context, in *generatedAuth.User) (*generatedAuth.User, error) {
	idTmp, err := uuid.Parse(in.Id)
	idTmpPh, err := uuid.Parse(in.ProfilePhoto)
	if err != nil {
		return &generatedAuth.User{Error: models.WrongData.Error()}, nil
	}
	user := models.User{
		Id:           idTmp,
		Login:        in.Login, //TODO: наверное здесь можно заполнять не все поля
		Name:         in.Name,
		ProfilePhoto: idTmpPh,
		PasswordHash: in.PasswordHash,
		Registration: time.Time{},
		UserVersion:  0,
	}
	checkedUser, err := h.uc.CheckUser(ctx, user)
	if err == nil {
		return &generatedAuth.User{
			Id:           checkedUser.Id.String(),
			Login:        in.Login, //TODO: наверное здесь можно заполнять не все поля
			Name:         in.Name,
			ProfilePhoto: checkedUser.ProfilePhoto.String(),
			PasswordHash: in.PasswordHash,
			Registration: checkedUser.Registration.String(),
			UserVersion:  0,
			Error:        "",
		}, nil
	}
	return &generatedAuth.User{Error: err.Error()}, nil

}

func (h GrpcAuthHandler) EncryptPwd(ctx context.Context, in *generatedAuth.EncryptPwdMg) (*generatedAuth.EncryptPwdMg, error) {
	return &generatedAuth.EncryptPwdMg{Password: h.uc.EncryptPwd(ctx, in.Password)}, nil
}
