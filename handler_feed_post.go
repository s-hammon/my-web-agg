package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/my-web-agg/internal/database"
)

func (a *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	type response struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feed, err := a.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Name:      params.Name,
		Url:       params.URL,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feedFollow, err := a.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    feed.UserID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response{
		Feed:       dbToFeed(feed),
		FeedFollow: dbToFeedFollow(feedFollow),
	})
}
