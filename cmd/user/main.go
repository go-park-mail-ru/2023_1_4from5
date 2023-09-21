package main

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	grpcUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc"
	generatedUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	userRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/repo"
	userUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
func run() error {

	logger, err := utils.FileLogger("/var/log/user_app.log")
	if err != nil {
		return err
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Print(err)
		}
	}(logger)

	zapSugar := logger.Sugar()
	utils.Init()
	str, err := utils.GetConnectionString()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", str)
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	userRepo := userRepository.NewUserRepo(db, zapSugar)
	userUse := userUsecase.NewUserUsecase(userRepo, zapSugar)
	service := grpcUser.NewGrpcUserHandler(userUse)

	srv, ok := net.Listen("tcp", ":8020")
	if ok != nil {
		log.Fatalln("can't listen port", err)
	}

	metricsMw := middleware.NewMetricsMiddleware()
	metricsMw.Register(middleware.ServiceUserName)

	server := grpc.NewServer(grpc.UnaryInterceptor(metricsMw.ServerMetricsInterceptor))

	generatedUser.RegisterUserServiceServer(server, service)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: ":8021"}

	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil {
			fmt.Print(err)
		}
	}()

	fmt.Print("user running on: ", srv.Addr())
	return server.Serve(srv)
}
