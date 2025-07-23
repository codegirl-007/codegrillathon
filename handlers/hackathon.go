package handlers

import (
	"codegrillathon/internals/database"
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
)

type HackathonForm struct {
	HackathonName string
	Description   string
}

type Hackathon struct {
	Id            int
	HackathonName string
	OwnerId       int
	StartDate     string
	EndDate       string
	Description   string
	Provider      string
}

func (h *Handler) Hackathon(w http.ResponseWriter, r *http.Request) {
	err := h.Template.ExecuteTemplate(w, "hackathon.html", nil)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (h *Handler) CreateHackathon(w http.ResponseWriter, r *http.Request) {
	err := h.Template.ExecuteTemplate(w, "create-hackathon.html", nil)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (h *Handler) ParseHackthonForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	hackathonName := r.FormValue("hackathonName")
	description := r.FormValue("description")
	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")

	if hackathonName == "" || description == "" || startDate == "" || endDate == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
	}

	session, err := gothic.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Error retrieving session for welcome page", http.StatusInternalServerError)
	}

	owner_id, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		return
	}

	provider, ok := session.Values["provider"].(string)
	if !ok {
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		return
	}

	dbClient, err := database.GetDbClientInstance()

	if err != nil {
		http.Error(w, "DB connection error", http.StatusInternalServerError)
		return
	}

	_, err = dbClient.Exec(
		`INSERT INTO hackathons (hackathon_name, owner_id, start_date, end_date, description, provider) VALUES (?, ?, ?, ?, ?, ?)`,
		hackathonName,
		owner_id,
		startDate,
		endDate,
		description,
		provider,
	)

	if err != nil {
		http.Error(w, "Database insert error", http.StatusInternalServerError)
		return
	}

	err = h.Template.ExecuteTemplate(w, "hackathon.html", nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func (h *Handler) ListHackathonsByProvider(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	dbClient, err := database.GetDbClientInstance()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
	}

	rows, err := dbClient.Query(
		`SELECT id, hackathon_name, owner_id, start_date, end_date, description, provider FROM hackathons WHERE provider = ?`,
		provider,
	)

	if err != nil {
		http.Error(w, "Error fetching from database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var hackathons []Hackathon

	for rows.Next() {
		var hackathon Hackathon
		err := rows.Scan(
			&hackathon.Id,
			&hackathon.HackathonName,
			&hackathon.OwnerId,
			&hackathon.StartDate,
			&hackathon.EndDate,
			&hackathon.Description,
			&hackathon.Provider,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}
		hackathons = append(hackathons, hackathon)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error iterating rows: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Provider   string
		Hackathons []Hackathon
	}{
		Provider:   provider,
		Hackathons: hackathons,
	}

	err = h.Template.ExecuteTemplate(w, "hackathons.html", data)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func (h *Handler) ListHackathonsByUser(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	username := r.PathValue("user")

	dbClient, err := database.GetDbClientInstance()
	if err != nil {
		fmt.Println("error using dbclient")
	}

	rows, err := dbClient.Query(
		"SELECT h.*, u.username FROM hackathon h INNER JOIN users u ON h.owner_id = u.id WHERE h.provider = ? AND h.owner_id = ? ",
		provider,
		username)

	// render the page
	if err != nil {
		http.Error(w, "Error fetching from database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var hackathons []Hackathon

	for rows.Next() {
		var hackathon Hackathon
		err := rows.Scan(
			&hackathon.Id,
			&hackathon.HackathonName,
			&hackathon.OwnerId,
			&hackathon.StartDate,
			&hackathon.EndDate,
			&hackathon.Description,
			&hackathon.Provider,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}
		hackathons = append(hackathons, hackathon)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error iterating rows: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Provider   string
		Hackathons []Hackathon
	}{
		Provider:   provider,
		Hackathons: hackathons,
	}

	err = h.Template.ExecuteTemplate(w, "hackathons.html", data)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}

}
