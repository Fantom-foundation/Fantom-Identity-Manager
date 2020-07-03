package login

import (
	"net/http"

	"github.com/volatiletech/authboss"
)

type logoutResponse struct {
}

func initiateLogout(challenge string) logoutResponse {
	var res logoutResponse
	url := makeGetURL(logout, challenge)
	_ = getJSON(url, &res)

	return res
}

type acceptLogoutResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func acceptLogoutRequest(challenge string) acceptLogoutResponse {
	var res acceptLogoutResponse
	url := makeAcceptURL(logout, challenge)
	_ = putJSON(url, nil, &res)

	return res
}

func LogoutMiddleware(ab *authboss.Authboss) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/logout" && r.Method == http.MethodGet {
				if ch := r.URL.Query().Get(getChallengeName(logout)); ch != "" {
					initiateLogout(ch)
					res := acceptLogoutRequest(ch)
					ab.Paths.LogoutOK = res.RedirectTo
				}
			}

			handler.ServeHTTP(w, r)
		})
	}
}
