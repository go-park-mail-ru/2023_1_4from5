package grpcAuth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	mocks "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
)

func TestGrpcAuthHandler_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mocks.NewMockAuthUsecase(ctrl)
	handler := NewGrpcAuthHandler(usecase)

	client, closer := startGRPCServer(handler)
	defer closer()

	tests := []struct {
		name string
		in   *generated.LoginUser
		out  *generated.Token
		err  error
		mock func()
	}{
		{
			name: "OK",
			in: &generated.LoginUser{
				Login:        "test",
				PasswordHash: "test",
			},
			out: &generated.Token{Cookie: "test", Error: ""},
			err: nil,
			mock: func() {
				usecase.EXPECT().SignIn(gomock.Any(), gomock.Any()).Times(1).Return("test", nil)
			},
		},
		{
			name: "Error",
			in: &generated.LoginUser{
				Login:        "test",
				PasswordHash: "test",
			},
			out: &generated.Token{Cookie: "", Error: errors.New("test").Error()},
			err: nil,
			mock: func() {
				usecase.EXPECT().SignIn(gomock.Any(), gomock.Any()).Times(1).Return("", errors.New("test"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			token, err := client.SignIn(context.Background(), test.in)

			require.Equal(t, test.out, token, fmt.Errorf("%s :  expected %s, got %s",
				test.name, test.out, token))
			require.Equal(t, nil, err, fmt.Errorf("error wasnt expected, got %s",
				err))
		})
	}
}

func TestGrpcAuthHandler_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mocks.NewMockAuthUsecase(ctrl)
	handler := NewGrpcAuthHandler(usecase)

	client, closer := startGRPCServer(handler)
	defer closer()

	tests := []struct {
		name string
		in   *generated.User
		out  *generated.Token
		err  error
		mock func()
	}{
		{
			name: "OK",
			in: &generated.User{
				Login:        "test_login",
				Name:         "test_name",
				PasswordHash: "test_pass",
			},
			out: &generated.Token{Cookie: "test", Error: ""},
			err: nil,
			mock: func() {
				usecase.EXPECT().SignUp(gomock.Any(), gomock.Any()).Times(1).Return("test", nil)
			},
		},
		{
			name: "Error",
			in: &generated.User{
				Login:        "test_login",
				Name:         "test_name",
				PasswordHash: "test_pass",
			},
			out: &generated.Token{Cookie: "", Error: errors.New("test").Error()},
			err: nil,
			mock: func() {
				usecase.EXPECT().SignUp(gomock.Any(), gomock.Any()).Times(1).Return("", errors.New("test"))
			},
		},
		{
			name: "Error while parsing uuid",
			in: &generated.User{
				Login:        "test_login",
				Name:         "test_name",
				PasswordHash: "test_pass",
			},
			out: &generated.Token{Cookie: "", Error: errors.New("test").Error()},
			err: nil,
			mock: func() {
				usecase.EXPECT().SignUp(gomock.Any(), gomock.Any()).Times(1).Return("", errors.New("test"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			token, err := client.SignUp(context.Background(), test.in)

			require.Equal(t, test.out, token, fmt.Errorf("%s :  expected %s, got %s",
				test.name, test.out, token))
			require.Equal(t, nil, err, fmt.Errorf("error wasnt expected, got %s",
				err))
		})
	}
}

func TestGrpcAuthHandler_IncUserVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mocks.NewMockAuthUsecase(ctrl)
	handler := NewGrpcAuthHandler(usecase)

	client, closer := startGRPCServer(handler)
	defer closer()

	tests := []struct {
		name string
		in   *generated.AccessDetails
		out  *generatedCommon.Empty
		err  error
		mock func()
	}{
		{
			name: "OK",
			in: &generated.AccessDetails{
				Login:       "test_login",
				Id:          uuid.New().String(),
				UserVersion: 2,
			},
			out: &generatedCommon.Empty{Error: ""},
			err: nil,
			mock: func() {
				usecase.EXPECT().IncUserVersion(gomock.Any(), gomock.Any()).Times(1).Return(int64(3), nil)
			},
		},
		{
			name: "Error",
			in: &generated.AccessDetails{
				Login:       "test_login",
				Id:          uuid.New().String(),
				UserVersion: 2,
			},
			out: &generatedCommon.Empty{Error: errors.New("test").Error()},
			err: nil,
			mock: func() {
				usecase.EXPECT().IncUserVersion(gomock.Any(), gomock.Any()).Times(1).Return(int64(3), errors.New("test"))
			},
		},
		{
			name: "Error while parsing uuid",
			in: &generated.AccessDetails{
				Login:       "test_login",
				Id:          "test",
				UserVersion: 2,
			},
			out: &generatedCommon.Empty{Error: models.WrongData.Error()},
			err: nil,
			mock: func() {
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			token, err := client.IncUserVersion(context.Background(), test.in)

			require.Equal(t, test.out, token, fmt.Errorf("%s :  expected %s, got %s",
				test.name, test.out, token))
			require.Equal(t, nil, err, fmt.Errorf("error wasnt expected, got %s",
				err))
		})
	}
}

func TestGrpcAuthHandler_CheckUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mocks.NewMockAuthUsecase(ctrl)
	handler := NewGrpcAuthHandler(usecase)

	client, closer := startGRPCServer(handler)
	defer closer()

	tests := []struct {
		name string
		in   *generated.User
		out  *generated.User
		err  error
		mock func()
	}{
		{
			name: "OK",
			in: &generated.User{
				Login:        "test_login",
				PasswordHash: "test_pass",
				Id:           uuid.New().String(),
			},
			out: &generated.User{Error: "", Id: uuid.Nil.String(), ProfilePhoto: uuid.Nil.String(), Login: "test_login", PasswordHash: "test_pass", Registration: time.Time{}.String()},
			err: nil,
			mock: func() {
				usecase.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Times(1).Return(models.User{}, nil)
			},
		},
		{
			name: "Error",
			in: &generated.User{
				Login:        "test_login",
				PasswordHash: "test_pass",
				Id:           uuid.New().String(),
			},
			out: &generated.User{Error: errors.New("test").Error()},
			err: nil,
			mock: func() {
				usecase.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Times(1).Return(models.User{}, errors.New("test"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			token, err := client.CheckUser(context.Background(), test.in)

			require.Equal(t, test.out, token, fmt.Errorf("%s :  expected %s, got %s",
				test.name, test.out, token))
			require.Equal(t, nil, err, fmt.Errorf("error wasnt expected, got %s",
				err))
		})
	}
}

func TestGrpcAuthHandler_CheckUserVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mocks.NewMockAuthUsecase(ctrl)
	handler := NewGrpcAuthHandler(usecase)

	client, closer := startGRPCServer(handler)
	defer closer()
	_, err := uuid.Parse("uuid.New().String()")

	tests := []struct {
		name string
		in   *generated.AccessDetails
		out  *generated.UserVersion
		err  error
		mock func()
	}{
		{
			name: "OK",
			in: &generated.AccessDetails{
				Login:       "test_login",
				UserVersion: int64(2),
				Id:          uuid.New().String(),
			},
			out: &generated.UserVersion{Error: "", UserVersion: int64(2)},
			err: nil,
			mock: func() {
				usecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Times(1).Return(int64(2), nil)
			},
		},
		{
			name: "Error while parsing uuid",
			in: &generated.AccessDetails{
				Login:       "test_login",
				UserVersion: int64(2),
				Id:          "uuid.New().String()",
			},
			out: &generated.UserVersion{Error: err.Error()},
			err: nil,
			mock: func() {
			},
		},
		{
			name: "Error",
			in: &generated.AccessDetails{
				Login:       "test_login",
				UserVersion: int64(2),
				Id:          uuid.New().String(),
			},
			out: &generated.UserVersion{Error: errors.New("test").Error()},
			err: nil,
			mock: func() {
				usecase.EXPECT().CheckUserVersion(gomock.Any(), gomock.Any()).Times(1).Return(int64(2), errors.New("test"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			token, err := client.CheckUserVersion(context.Background(), test.in)

			require.Equal(t, test.out, token, fmt.Errorf("%s :  expected %s, got %s",
				test.name, test.out, token))
			require.Equal(t, nil, err, fmt.Errorf("error wasnt expected, got %s",
				err))
		})
	}
}

func startGRPCServer(impl generated.AuthServiceServer) (generated.AuthServiceClient, func()) {
	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	srv := grpc.NewServer()
	generated.RegisterAuthServiceServer(srv, impl)

	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err = listener.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		srv.Stop()
	}

	cl := generated.NewAuthServiceClient(conn)
	return cl, closer
}
