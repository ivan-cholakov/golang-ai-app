package handler

import (
	"dreampicai/db"
	"dreampicai/view/home"
	"net/http"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	account, err := db.GetAccountByUserID(user.ID)
	if err != nil {
		return err
	}

	return home.Index().Render(r.Context(), w)
}
