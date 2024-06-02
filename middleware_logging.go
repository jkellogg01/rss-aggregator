package main

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

func middlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%7s @ %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("finished in %v", time.Since(start))
	})
}
