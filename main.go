package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/rss-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Failed to convert port to an integer", "error", err)
	}
	db, err := sql.Open("postgres", os.Getenv("PG_DBSTRING"))
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	dbQueries := database.New(db)
	cfg := apiConfig{DB: dbQueries}

	mux := http.NewServeMux()
    
	mux.HandleFunc("GET /v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, 200, map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("GET /v1/err", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, 500, "Internal Server Error")
	})

	mux.HandleFunc("POST /v1/users", cfg.handleCreateUser)
    mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.handleGetUser))

	app := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	err = app.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	codeGroup := code / 100
	clientErrorGroup := 4
	serverErrorGroup := 5
	if codeGroup != clientErrorGroup && codeGroup != serverErrorGroup {
		log.Warn("responding with error using a non-error response code", "status", code)
	}
	respondWithJSON(w, code, map[string]string{
		"error": message,
	})
}
