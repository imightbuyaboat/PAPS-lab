package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	sessman "papslab/sessionmanager"
	"papslab/studiodb"
	passman "papslab/studiodb/passwordmanager"
	"papslab/studiodb/register"
)

const (
	templateDir = "templates"
)

type Handler struct {
	sm   *sessman.SessionManager
	pm   *passman.PasswordManager
	reg  *register.Register
	tmpl *template.Template
}

func NewHandler() (*Handler, error) {
	files, err := filepath.Glob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске файлов: %v", err)
	}

	newSM, err := sessman.NewSessionManager()
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к redis: %v", err)
	}

	db, err := studiodb.NewDB()
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %v", err)
	}

	return &Handler{
		sm:   newSM,
		pm:   passman.NewPasswordManager(db),
		reg:  register.NewRegister(db),
		tmpl: template.Must(template.ParseFiles(files...)),
	}, nil
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	exists, isPriv, err := h.pm.Check(&passman.User{
		Login:    login,
		Password: password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		errorMsg := "Некорректные логин или пароль"
		err := h.tmpl.ExecuteTemplate(w, "login.html", map[string]string{"ErrorMsg": errorMsg})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	sessionID, err := h.sm.Create(&sessman.Session{
		Login:      login,
		Useragent:  r.UserAgent(),
		Priveleged: isPriv,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionID.String(),
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	var errorMsg string

	if login == "" || password == "" {
		errorMsg = "Некорректные логин или пароль"
		err := h.tmpl.ExecuteTemplate(w, "register.html", map[string]string{"ErrorMsg": errorMsg})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	exists, err := h.pm.IsLoginAvailable(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		errorMsg = "Пользователь с таким логином уже существует"
		err := h.tmpl.ExecuteTemplate(w, "register.html", map[string]string{"ErrorMsg": errorMsg})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = h.pm.Insert(&passman.User{
		Login:    login,
		Password: password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := h.sm.Create(&sessman.Session{
		Login:      login,
		Useragent:  r.UserAgent(),
		Priveleged: false,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionID.String(),
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := sessman.ParseSessionID(session.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.sm.Delete(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	session, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else if !session.Priveleged {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = h.reg.Insert(
		register.Item{
			Id:           0,
			Organization: r.FormValue("organization"),
			City:         r.FormValue("city"),
			Phone:        r.FormValue("phone"),
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	session, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else if !session.Priveleged {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.reg.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	var isPriv bool

	session, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else {
		isPriv = session.Priveleged
	}

	items, err := h.reg.SelectAny(
		register.Item{
			Id:           0,
			Organization: r.FormValue("organization"),
			City:         r.FormValue("city"),
			Phone:        r.FormValue("phone"),
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "index.html", struct {
		Items       []register.Item
		ShowButtons bool
	}{
		Items:       items,
		ShowButtons: isPriv,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ReturnToMainPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	var isPriv bool

	session, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else {
		isPriv = session.Priveleged
	}

	items, err := h.reg.SelectAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "index.html", struct {
		Items       []register.Item
		ShowButtons bool
	}{
		Items:       items,
		ShowButtons: isPriv,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) checkSession(r *http.Request) (*sessman.Session, error) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sessionID, err := sessman.ParseSessionID(cookieSessionID.Value)
	if err != nil {
		return nil, err
	}

	return h.sm.Check(sessionID)
}

func (h *Handler) CheckSessionMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.checkSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if session != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
