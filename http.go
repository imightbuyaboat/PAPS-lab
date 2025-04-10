package main

import (
	"net/http"
	"strconv"
	"time"

	passman "PAPS-LAB/passwordmanager"
	"PAPS-LAB/register"
	sessman "PAPS-LAB/sessionmanager"
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

	var errorMsg string

	exists, isPriv, err := h.pm.Check(&passman.User{
		Login:    login,
		Password: password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		errorMsg = "Некорректные логин или пароль"
	}

	if errorMsg != "" {
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
		Value:   sessionID.ID,
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

	exists, err := h.pm.CheckAvailableLogin(login)
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
		Value:   sessionID.ID,
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.sm.Delete(&sessman.SessionID{
		ID: session.Value,
	})

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	session, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if !session.Priveleged {
		http.Redirect(w, r, "/", http.StatusFound)
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

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	session, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if !session.Priveleged {
		http.Redirect(w, r, "/", http.StatusFound)
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

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	var isPriv bool

	session, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
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

	session := h.sm.Check(&sessman.SessionID{ID: cookieSessionID.Value})
	return session, nil
}
