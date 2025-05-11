package main

import (
	"log"
	"net/http"
	"papslab/handler"

	"github.com/gorilla/mux"
)

const (
	staticDir = "static"
)

func main() {
	handlers, err := handler.NewHandler()
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir(staticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	checkSessionRouter := r.NewRoute().Subrouter()
	checkSessionRouter.Use(handlers.CheckSessionMiddleWare)
	checkSessionRouter.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	checkSessionRouter.HandleFunc("/register", handlers.RegisterPage).Methods("GET")

	r.HandleFunc("/", handlers.MainPage).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/logout", handlers.Logout).Methods("POST")
	r.HandleFunc("/search", handlers.Search).Methods("POST")
	r.HandleFunc("/add", handlers.Add).Methods("POST")
	r.HandleFunc("/delete", handlers.Delete).Methods("POST")
	r.HandleFunc("/return", handlers.ReturnToMainPage).Methods("POST")

	log.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
