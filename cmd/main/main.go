package main

import (
	"database/sql"
	"fmt"
	authDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/http"
	authRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
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

	str, err := middleware.GetConnectionString()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", str)
	if err != nil {
		return err
	}
	defer db.Close()

	tokenGenerator := authUsecase.NewTokenator()
	authRepo := authRepository.NewAuthRepo(db)
	authUse := authUsecase.NewAuthUsecase(authRepo, tokenGenerator)
	authHandler := authDelivery.NewAuthHandler(authUse)

	auth := r.PathPrefix("/user").Subrouter()
	{
		auth.HandleFunc("/signUp", authHandler.SignUp).Methods(http.MethodPost)
		auth.HandleFunc("/signIn", authHandler.SignIn).Methods(http.MethodGet)
		auth.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost, http.MethodOptions)
	}

	http.Handle("/", r)
	return srv.ListenAndServe()
}
