package main

import (
	"net/http"

	"github.com/s-hammon/my-web-agg/internal/database"
)

func (a *apiConfig) handlerGetFeedFollowsByUserID(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := a.DB.GetFeedFollowsByID(r.Context(), user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respFeedFollows := []FeedFollow{}
	for _, ff := range feedFollows {
		respFeedFollows = append(respFeedFollows, dbToFeedFollow(ff))
	}

	respondJSON(w, http.StatusOK, respFeedFollows)
}
