package main

import (
	"net/http"

	"github.com/s-hammon/my-web-agg/internal/database"
)

func (a *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondJSON(w, http.StatusOK, dbToUser(user))
}
