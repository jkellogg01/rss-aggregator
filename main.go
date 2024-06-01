package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Failed to convert port to an integer", "error", err)
	}

	mux := http.NewServeMux()
    mux.HandleFunc("GET /v1/healthz", func (w http.ResponseWriter, r *http.Request) {
        respondWithJSON(w, 200, map[string]string{
            "status": "ok",
        })
    })

    mux.HandleFunc("GET /v1/err", func (w http.ResponseWriter, r *http.Request) {
        respondWithError(w, 500, "Internal Server Error")  
    })

	app := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	err = app.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	return err
}

func respondWithError(w http.ResponseWriter, code int, message string) error {
    codeGroup := code / 100
    clientErrorGroup := 4
    serverErrorGroup := 5
    if codeGroup != clientErrorGroup && codeGroup != serverErrorGroup {
        log.Warn("responding with error using a non-error response code", "status", code)
    }
	return respondWithJSON(w, code, map[string]string{
		"error": message,
	})
}
