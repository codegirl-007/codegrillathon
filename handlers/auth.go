package handlers

import (
	"fmt"
	"net/http"

	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "twitch")
	r.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(w, r)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := os.Getenv("SESSION_SECRET") // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30               // 30 days
	isProd := false                    // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
	session, _ := gothic.Store.Get(r, "user-session")
	session.Values["user_name"] = user.Name
	session.Values["avatar_url"] = user.AvatarURL

	err = session.Save(r, w)
	if err != nil {
		fmt.Printf("error saving the session: %v", err)
	}

	http.Redirect(w, r, "/welcome", http.StatusFound)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, "user-session")
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

	http.Redirect(w, r, "/", http.StatusFound)
}
