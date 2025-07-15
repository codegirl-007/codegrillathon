package handlers

import (
	"html/template"
	"net/http"
)

type Handler struct {
	Template template.Template
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	err := h.Template.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
