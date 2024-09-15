package main

import (
	"net/http"
	"strconv"

	"github.com/s-hammon/my-web-agg/internal/database"
)

func (a *apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limit := 15
	reqLimit := r.URL.Query().Get("limit")
	if reqLimit != "" {
		intLimit, err := strconv.Atoi(reqLimit)
		if err != nil {
			respondError(w, http.StatusBadRequest, "limit must be an integer")
			return
		}

		limit = intLimit
	}

	posts, err := a.DB.GetPostsByUserID(r.Context(), database.GetPostsByUserIDParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respPosts := make([]Post, len(posts))
	for i, post := range posts {
		respPosts[i] = dbToPost(post)
	}

	respondJSON(w, http.StatusOK, respPosts)
}
