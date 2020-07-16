package auth

import (
	"context"
	"errors"
	"key-vault-server/config"
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

func AuthenticationMiddleware(cfg *config.Config) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == cfg.GraphqlEntrypoint {
				token, err := extractToken(r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
				userId, err := getUserId(cfg, token)
				if err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}

				// Store extracted values to context
				user := User{}
				user.Uuid = *userId
				ctx := context.WithValue(r.Context(), userCtxKey, user)
				// Next call
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUser(ctx context.Context) *User {
	user, _ := ctx.Value(userCtxKey).(User)
	return &user
}
