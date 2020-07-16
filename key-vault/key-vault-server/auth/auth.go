package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

var (
	userCtxKey     = &contextKey{"user"}
	malformedToken = errors.New("malformed token")
)

type contextKey struct {
	name string
}

// User context information
type User struct {
	Uuid string
}

type Middleware func(http.Handler) http.Handler

func extractToken(r *http.Request) (*string, error) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return nil, malformedToken
	} else {
		token := authHeader[1]
		return &token, nil
	}
}

func extractUserInfo(_ *string) (User, error) {
	return User{}, nil
}

func ValidateTokenActive(_ *string) error {
	return nil
}

func AuthenticationMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/query" {
				token, err := extractToken(r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				user, err := extractUserInfo(token)
				if err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				err = ValidateTokenActive(token)
				if err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}

				// Store extracted values to context
				ctx := context.WithValue(r.Context(), userCtxKey, user)
				// Next call
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}
