package auth

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrNoAuthHeader = errors.New("no authorization header found")
    ErrBadAuthHeader = errors.New("malformed authorization header")
)

func GetAPIKey(headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return "", ErrNoAuthHeader
    }
    key, found := strings.CutPrefix(authHeader, "ApiKey ")
    if !found {
        return "", ErrBadAuthHeader
    }

    return key, nil
}
