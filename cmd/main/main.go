package main

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	"github.com/gorilla/mux"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
func run() error {
	r := mux.NewRouter()

	r.Use(middleware.CORSMiddleware)

	//srv := http.Server{Handler: r, Addr: fmt.Sprintf(":%s", "8000")}

	//conn, err := config.GetConnectionString() //вытаскивать зашитые данные придумать!!!
	//if err != nil {
	//	return err
	//}
	//
	//pool, err := pgxpool.Connect(context.Background(), conn)
	//if err != nil {
	//	return err
	//}
	//
	//tokenGenerator := authUsecase.NewTokenator()
	//onlineRepo := authRepository.NewOnlineRepo(pool)
	//authRepo := authRepository.NewAuthRepo(pool)
	//authUse := authUsecase.NewAuthUsecase(authRepo, tokenGenerator)
	//authHandler := authDelivery.NewAuthHandler(authUse, onlineRepo)
	//
	//auth := r.PathPrefix("/user").Subrouter()
	//{
	//	auth.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	//	auth.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost, http.MethodOptions)
	//	auth.HandleFunc("/signup", authHandler.SignUp).Methods(http.MethodPost)
	//	auth.HandleFunc("/auth", authHandler.AuthStatus).Methods(http.MethodGet)
	//}
	//
	//http.Handle("/", r)
	//log.Print("main running on: ", srv.Addr)
	//return srv.ListenAndServe()
	return nil
}
