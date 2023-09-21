package main

import (
	"database/sql"
	"fmt"
	grpcAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	authRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
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

	logger, err := utils.FileLogger("/var/log/auth_app.log")
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

	tokenGenerator := authUsecase.NewTokenator()
	encryptor, err := authUsecase.NewEncryptor()
	if err != nil {
		return err
	}

	authRepo := authRepository.NewAuthRepo(db, zapSugar)
	authUse := authUsecase.NewAuthUsecase(authRepo, tokenGenerator, encryptor, zapSugar)
	service := grpcAuth.NewGrpcAuthHandler(authUse)

	srv, ok := net.Listen("tcp", ":8010")
	if ok != nil {
		log.Fatalln("can't listen port", err)
	}

	metricsMw := middleware.NewMetricsMiddleware()
	metricsMw.Register(middleware.ServiceAuthName)

	server := grpc.NewServer(grpc.UnaryInterceptor(metricsMw.ServerMetricsInterceptor))

	generatedAuth.RegisterAuthServiceServer(server, service)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: ":8011"}

	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil {
			fmt.Print(err)
		}
	}()

	fmt.Print("auth running on: ", srv.Addr())
	return server.Serve(srv)
}
