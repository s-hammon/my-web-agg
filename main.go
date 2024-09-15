package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/s-hammon/my-web-agg/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("couldn't find env var: PORT")
	}
	dbURL := os.Getenv("CONN_STRING")
	if dbURL == "" {
		log.Fatal("couldn't find env var: CONN_STRING")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	cfg := &apiConfig{DB: dbQueries}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerError)

	mux.HandleFunc("POST /v1/users", cfg.handlerCeateUser)
	mux.HandleFunc("GET /v1/users/", cfg.middlewareAuth(cfg.handlerGetUserByAPIKey))

	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.handlerCreateFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.handlerGetFeeds)

	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.handlerCreateFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.handlerGetFeedFollowsByUserID))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.handlerDeleteFeedFollow))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	const requestConcurrency = 10
	const requestInterval = time.Minute
	go scrapeWorker(dbQueries, requestConcurrency, requestInterval)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
