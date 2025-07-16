package handlers

import (
	"net/http"
)

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) {
	err := h.Template.ExecuteTemplate(w, "welcome.html", nil)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
