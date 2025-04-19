package main

import (
	"net/http"
	"strconv"
	"time"

	bt "papslab/basic_types"
	sessman "papslab/sessionmanager"
)

func (h *Handler) loginPage(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	exists, isPriv, err := h.pm.Check(&bt.User{
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

	sessionID, err := h.sm.Create(&bt.Session{
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

func (h *Handler) registerPage(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
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

	err = h.pm.Insert(&bt.User{
		Login:    login,
		Password: password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := h.sm.Create(&bt.Session{
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

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	err := h.reg.Insert(
		bt.Item{
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

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
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
		bt.Item{
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
		Items       []bt.Item
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

func (h *Handler) returnToMainPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
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
		Items       []bt.Item
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

func (h *Handler) checkSession(r *http.Request) (*bt.Session, error) {
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

func (h *Handler) CheckAuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		next.ServeHTTP(w, r)
	})
}
