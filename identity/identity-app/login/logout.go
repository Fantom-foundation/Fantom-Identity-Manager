package login

import (
	"net/http"

	"github.com/volatiletech/authboss"
)

type logoutResponse struct {
}

func initiateLogout(challenge string, hydra *Hydra) logoutResponse {
	var res logoutResponse
	url := hydra.makeGetURL(logout, challenge)
	_ = hydra.getJSON(url, &res)

	return res
}

type acceptLogoutResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func acceptLogoutRequest(challenge string, hydra *Hydra) acceptLogoutResponse {
	var res acceptLogoutResponse
	url := hydra.makeAcceptURL(logout, challenge)
	_ = hydra.putJSON(url, nil, &res)

	return res
}

func LogoutMiddleware(ab *authboss.Authboss, hydra *Hydra) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/logout" && r.Method == http.MethodGet {
				if ch := r.URL.Query().Get(hydra.getChallengeName(logout)); ch != "" {
					initiateLogout(ch, hydra)
					res := acceptLogoutRequest(ch, hydra)
					ab.Paths.LogoutOK = res.RedirectTo
				}
			}

			handler.ServeHTTP(w, r)
		})
	}
}
