package handlers

import (
	"github.com/markbates/goth/gothic"
	"net/http"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	err := h.Template.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, "_gothic-session")
	if err != nil {
		return
	}

	// Clear the session data
	session.Values = make(map[interface{}]interface{})

	// Save the empty session
	err = session.Save(r, w)
	if err != nil {
		return
	}

	//w.Redirect(http.StatusTemporaryRedirect, "/")
}
