package main

import (
	"net/http"
	"strconv"
	"time"

	passman "PAPS-LAB/passwordmanager"
	sessman "PAPS-LAB/sessionmanager"
	"PAPS-LAB/studiodb"
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

	errorID := h.pm.Check(&passman.User{
		Login:    login,
		Password: password,
	})

	switch errorID {
	case 1:
		errorMsg = "Неправильный пароль"
	case 2:
		errorMsg = "Такого пользователя не существует"
	}

	if errorMsg != "" {
		err := h.tmpl.ExecuteTemplate(w, "login.html", map[string]string{"ErrorMsg": errorMsg})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	sessionID, err := h.sm.Create(&sessman.Session{
		Login:     login,
		Useragent: r.UserAgent(),
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

	errorID := h.pm.CheckAvailableLogin(login)

	if errorID == 1 {
		errorMsg = "Пользователь с таким логином уже существует"
		err := h.tmpl.ExecuteTemplate(w, "register.html", map[string]string{"ErrorMsg": errorMsg})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err := h.pm.Create(&passman.User{
		Login:    login,
		Password: password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := h.sm.Create(&sessman.Session{
		Login:     login,
		Useragent: r.UserAgent(),
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
	err := h.db.Insert(
		studiodb.Item{
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

	err = h.db.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	items, err := h.db.SelectAny(
		studiodb.Item{
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
		Items []studiodb.Item
	}{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) returnToMainPage(w http.ResponseWriter, r *http.Request) {
	items, err := h.db.SelectAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []studiodb.Item
	}{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	session, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	items, err := h.db.SelectAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []studiodb.Item
	}{
		Items: items,
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
