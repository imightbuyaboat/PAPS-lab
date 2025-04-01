package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	sessionID, err := h.sm.Create(&sessman.Session{
		Login:     login,
		Useragent: r.UserAgent(),
	})
	if err != nil {
		fmt.Println("hiiiii")
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

func (h *Handler) listPage(w http.ResponseWriter, r *http.Request) {
	items, err := h.db.SelectAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "list.html", struct {
		Items []studiodb.Item
	}{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) addPage(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "add.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	err := h.db.Insert(
		studiodb.Item{
			Id:           0,
			Organization: r.FormValue("Organization"),
			Phone:        r.FormValue("Phone"),
		})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) deletePage(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "delete.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	session, err := h.checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	err = h.tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
