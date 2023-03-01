package main

import (
	"database/sql"
	"fmt"
	authDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/http"
	authRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	creatorDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/http"
	creatorRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/repo"
	creatorUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	userDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/repo"
	userUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.Use(middleware.CORSMiddleware)

	srv := http.Server{Handler: r, Addr: fmt.Sprintf(":%s", "8000")}

	str, err := utils.GetConnectionString()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", str)
	if err != nil {
		return err
	}
	defer db.Close()

	tokenGenerator := authUsecase.NewTokenator()
	encryptor := authUsecase.NewEncryptor()

	authRepo := authRepository.NewAuthRepo(db)
	authUse := authUsecase.NewAuthUsecase(authRepo, tokenGenerator, encryptor)
	authHandler := authDelivery.NewAuthHandler(authUse)

	userRepo := userRepository.NewUserRepo(db)
	userUse := userUsecase.NewUserUsecase(userRepo)
	userHandler := userDelivery.NewUserHandler(userUse)

	creatorRepo := creatorRepository.NewCreatorRepo(db)
	creatorUse := creatorUsecase.NewCreatorUsecase(creatorRepo)
	creatorHandler := creatorDelivery.NewCreatorHandler(creatorUse)

	//TODO: придумать как отдавать статус авторства
	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/signUp", authHandler.SignUp).Methods(http.MethodPost)
		auth.HandleFunc("/signIn", authHandler.SignIn).Methods(http.MethodPost)
		auth.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodGet, http.MethodOptions)
	}

	user := r.PathPrefix("/user").Subrouter()
	{
		user.HandleFunc("/profile", userHandler.GetProfile).Methods(http.MethodGet)
	}

	creator := r.PathPrefix("/creator").Subrouter()
	{
		creator.HandleFunc("/myPage", creatorHandler.GetPage).Methods(http.MethodGet, http.MethodOptions)
	}
	http.Handle("/", r)
	return srv.ListenAndServe()
}
