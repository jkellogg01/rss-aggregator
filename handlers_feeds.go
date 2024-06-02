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

func (a *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, u database.User) {
    ctx := context.Background()
	bodyDecoder := json.NewDecoder(r.Body)
	var body struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
    bodyDecoder.Decode(&body)
    newID, err := uuid.NewRandom()
    if err != nil {
        log.Error("failed to generate uuid", "error", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    t := time.Now()
    result, err := a.DB.CreateFeed(ctx, database.CreateFeedParams{
        ID: newID,
        CreatedAt: t,
        UpdatedAt: t,
        Name: body.Name,
        Url: body.URL,
        UserID: u.ID,
    })
    if err != nil {
        log.Error("failed to get feeds from database", "error", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    respondWithJSON(w, 201, result)
}

func (a *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    results, err := a.DB.GetFeeds(ctx)
    if err != nil {
        log.Error("failed to get feeds from database", "error", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    respondWithJSON(w, 200, results)
}
