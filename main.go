package main

import (
	"html/template"
	"log"
	"net/http"
	passman "papslab/passwordmanager"
	"papslab/register"
	sessman "papslab/sessionmanager"
	"papslab/studiodb"
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
	templateDir = "templates"
	staticDir   = "static"
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

	fs := http.FileServer(http.Dir(staticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	checkSessionRouter := r.NewRoute().Subrouter()
	checkSessionRouter.Use(handlers.CheckSessionMiddleWare)
	checkSessionRouter.HandleFunc("/login", handlers.loginPage).Methods("GET")
	checkSessionRouter.HandleFunc("/register", handlers.registerPage).Methods("GET")

	r.HandleFunc("/", handlers.mainPage).Methods("GET")
	r.HandleFunc("/login", handlers.login).Methods("POST")
	r.HandleFunc("/register", handlers.register).Methods("POST")
	r.HandleFunc("/logout", handlers.logout).Methods("POST")
	r.HandleFunc("/search", handlers.search).Methods("POST")
	r.HandleFunc("/add", handlers.add).Methods("POST")
	r.HandleFunc("/delete", handlers.delete).Methods("POST")
	r.HandleFunc("/return", handlers.returnToMainPage).Methods("POST")

	log.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
