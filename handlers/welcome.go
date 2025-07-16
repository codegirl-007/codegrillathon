package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("we have reached the welcome route")

	err := h.Template.ExecuteTemplate(w, "welcome.html", nil)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
