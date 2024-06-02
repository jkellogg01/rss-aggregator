package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jkellogg01/rss-aggregator/internal/database"
)

func (a *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	bodyDecoder := json.NewDecoder(r.Body)
	var body struct {
		Name string `json:"name"`
	}
	err := bodyDecoder.Decode(&body)
	if err != nil {
		log.Error("failed to decode JSON body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := time.Now()
	newId, err := uuid.NewRandom()
	if err != nil {
		log.Error("failed to generate UUID", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := a.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        newId,
		Name:      body.Name,
		CreatedAt: t,
		UpdatedAt: t,
	})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, 201, resp)
}

func (a *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, u database.User) {
	respondWithJSON(w, 200, u)
}
