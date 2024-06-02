package main

import (
	"context"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/jkellogg01/rss-aggregator/internal/auth"
	"github.com/jkellogg01/rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
    ctx := context.Background()
    return func (w http.ResponseWriter, r *http.Request) {
        key, err := auth.GetAPIKey(r.Header)
        if err != nil {
            log.Error("failed to get API key", "error", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        user, err := cfg.DB.GetUser(ctx, key)
        if err != nil {
            log.Error("no user with this API key", "error", err)
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        handler(w, r, user)
    }
}
