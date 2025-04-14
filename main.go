package main

import (
	passman "PAPS-LAB/passwordmanager"
	"PAPS-LAB/register"
	sessman "PAPS-LAB/sessionmanager"
	"PAPS-LAB/studiodb"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Handler struct {
	sm   *sessman.SessionManager
	pm   *passman.PasswordManager
	reg  *register.Register
	tmpl *template.Template
}

const (
	templateDir string = "templates"
)

func main() {
	files, err := filepath.Glob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		log.Fatalf("Ошибка при поиске файлов: %v", err)
	}

	newSM, err := sessman.NewSessionManager()
	if err != nil {
		log.Fatalf("Ошибка при подключении к redis: %v", err)
	}

	db, err := studiodb.NewDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	handlers := &Handler{
		sm:   newSM,
		pm:   passman.NewPasswordManager(db),
		reg:  register.NewRegister(db),
		tmpl: template.Must(template.ParseFiles(files...)),
	}

	r := mux.NewRouter()

	privilegedRouter := r.NewRoute().Subrouter()
	privilegedRouter.Use(handlers.CheckAuthMiddleWare)
	privilegedRouter.HandleFunc("/add", handlers.add).Methods("POST")
	privilegedRouter.HandleFunc("/delete", handlers.delete).Methods("POST")

	r.HandleFunc("/", handlers.mainPage).Methods("GET")
	r.HandleFunc("/login", handlers.loginPage).Methods("GET")
	r.HandleFunc("/login", handlers.login).Methods("POST")
	r.HandleFunc("/register", handlers.registerPage).Methods("GET")
	r.HandleFunc("/register", handlers.register).Methods("POST")
	r.HandleFunc("/logout", handlers.logout).Methods("POST")
	r.HandleFunc("/search", handlers.search).Methods("POST")
	r.HandleFunc("/return", handlers.returnToMainPage).Methods("POST")

	fmt.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
