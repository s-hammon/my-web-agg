package main

import "net/http"

func (a *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := a.DB.GetFeeds(r.Context())
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respFeeds := make([]Feed, len(feeds))
	for i, feed := range feeds {
		respFeeds[i] = dbToFeed(feed)
	}

	respondJSON(w, http.StatusOK, respFeeds)
}
