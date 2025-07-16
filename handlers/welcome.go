package handlers

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
)

type PageData struct {
	Username string
}

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Error retrieving session for welcome page", http.StatusInternalServerError)
	}
	username, ok := session.Values["user_name"].(string)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	fmt.Printf("username: %s", username)
	err = h.Template.ExecuteTemplate(w, "welcome.html", &PageData{
		Username: username,
	})
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
