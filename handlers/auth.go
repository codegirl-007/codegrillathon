package handlers

import (
	"codegrillathon/internals/database"
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
	session.Values["user_id"] = user.UserID
	session.Values["provider"] = user.Provider

	err = session.Save(r, w)
	if err != nil {
		fmt.Printf("error saving the session: %v", err)
	}

	dbClient, err := database.GetDbClientInstance()

	rows, err := dbClient.Query("SELECT COUNT(*) FROM users WHERE username = ? AND provider = ?", user.Name, "twitch")
	if err != nil {
		fmt.Printf("error checking user in users table: %v\n", err)
	}

	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			fmt.Printf("error with rows: %v", err)
		}
	}

	if count == 0 {
		_, err = dbClient.Exec(
			"INSERT INTO users (username, user_cap_id, provider, avatar_url, provider_id) VALUES (?, ?, ?, ?, ?)",
			user.Name, 1, "twitch", user.AvatarURL, user.UserID,
		)

		if err != nil {
			fmt.Printf("error saving user to database: %v", err)
		}
	}

	http.Redirect(w, r, "/welcome", http.StatusFound)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, "user-session")
	if err != nil {
		fmt.Printf("error retrieving session: %v", err)
		return
	}

	// Clear the session data
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1

	// Save the empty session
	err = session.Save(r, w)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := gothic.Store.Get(r, "user-session")
		userID, ok := session.Values["user_id"]
		if !ok || userID == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
