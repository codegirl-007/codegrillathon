package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) Hackathon(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I hit the hackathon route...")
	err := h.Template.ExecuteTemplate(w, "hackathon.html", nil)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
