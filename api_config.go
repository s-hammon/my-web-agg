package main

import (
	"net/http"

	"github.com/s-hammon/my-web-agg/internal/auth"
	"github.com/s-hammon/my-web-agg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (a *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetToken("ApiKey", r.Header)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := a.DB.GetUserByAPIKey(r.Context(), token)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}

		handler(w, r, user)
	}
}
