package main

import (
	passman "PAPS-LAB/passwordmanager"
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
	db   *studiodb.DB
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

	db, err := studiodb.NewDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	handlers := &Handler{
		sm:   sessman.NewSessionManager(),
		pm:   passman.NewPasswordManager(),
		db:   db,
		tmpl: template.Must(template.ParseFiles(files...)),
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.mainPage).Methods("GET")
	r.HandleFunc("/login", handlers.loginPage).Methods("GET")
	r.HandleFunc("/login", handlers.login).Methods("POST")
	r.HandleFunc("/register", handlers.registerPage).Methods("GET")
	r.HandleFunc("/register", handlers.register).Methods("POST")
	r.HandleFunc("/logout", handlers.logout).Methods("POST")
	r.HandleFunc("/list", handlers.listPage).Methods("GET")
	r.HandleFunc("/add", handlers.addPage).Methods("GET")
	r.HandleFunc("/add", handlers.add).Methods("POST")
	r.HandleFunc("/delete", handlers.deletePage).Methods("GET")
	r.HandleFunc("/delete", handlers.delete).Methods("POST")

	fmt.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
