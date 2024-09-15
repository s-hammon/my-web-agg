package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/s-hammon/my-web-agg/internal/database"
)

func (a *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID := r.PathValue("feedFollowID")
	id, err := uuid.Parse(feedFollowID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     id,
	}); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
