package handler

import (
	"dreampicai/view/home"
	"net/http"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	// user := getAuthenticatedUser(r)

	return home.Index().Render(r.Context(), w)
}
