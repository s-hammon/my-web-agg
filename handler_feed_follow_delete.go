package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (a *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowID := r.PathValue("feedFollowID")
	if feedFollowID == "" {
		respondError(w, http.StatusBadRequest, "must provide feedFollowID")
		return
	}
	id, err := uuid.Parse(feedFollowID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.DB.DeleteFeedFollow(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
