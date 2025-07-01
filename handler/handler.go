package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"papslab/item"
	"papslab/session"
	"papslab/user"
	"path/filepath"
	"strconv"
	"time"
)

const (
	templateDir = "templates"
)

type Handler struct {
	sm   SessionManager
	strg Storage
	tmpl *template.Template
}

func NewHandler(sm SessionManager, strg Storage) (*Handler, error) {
	files, err := filepath.Glob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске файлов: %v", err)
	}

	return &Handler{
		sm:   sm,
		strg: strg,
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

	exists, isPriv, err := h.strg.CheckUser(&user.User{
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

	sessionID, err := h.sm.Create(&session.Session{
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

	exists, err := h.strg.IsLoginAvailable(login)
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

	err = h.strg.InsertUser(&user.User{
		Login:    login,
		Password: password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := h.sm.Create(&session.Session{
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
	sess, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := session.ParseSessionID(sess.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.sm.Delete(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sess.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sess)

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	sess, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sess == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else if !sess.Priveleged {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = h.strg.InsertItem(
		item.Item{
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
	sess, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sess == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else if !sess.Priveleged {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.strg.DeleteItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	var isPriv bool

	sess, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sess == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else {
		isPriv = sess.Priveleged
	}

	items, err := h.strg.SelectAnyItems(
		item.Item{
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
		Items       []item.Item
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

	sess, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sess == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	} else {
		isPriv = sess.Priveleged
	}

	items, err := h.strg.SelectAllItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "index.html", struct {
		Items       []item.Item
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

func (h *Handler) checkSession(r *http.Request) (*session.Session, error) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sessionID, err := session.ParseSessionID(cookieSessionID.Value)
	if err != nil {
		return nil, err
	}

	return h.sm.Check(sessionID)
}

func (h *Handler) CheckSessionMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := h.checkSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if sess != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
