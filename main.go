package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/s-hammon/my-web-agg/internal/auth"
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

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

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
