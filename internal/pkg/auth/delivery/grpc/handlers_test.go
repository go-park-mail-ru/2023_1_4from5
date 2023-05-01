package grpcAuth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	mocks "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
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

			token, err := client.SignIn(context.Background(), &generated.LoginUser{})

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
		}), grpc.WithInsecure())
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
