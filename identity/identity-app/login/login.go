package login

import (
	"context"
	"net/http"

	"github.com/volatiletech/authboss"

	"identity-app/model"
)

const (
	CTXKeyChallenge string = "challenge"
)

type Response struct {
	Skip    bool   `json:"skip"`
	Subject string `json:"subject"`
}

func initiateLogin(challenge string, hydra *Hydra) Response {
	var res Response
	url := hydra.makeGetURL(login, challenge)
	_ = hydra.getJSON(url, &res)

	return res
}

type acceptLoginResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func acceptLoginRequest(w http.ResponseWriter, r *http.Request, challenge string, subject string, hydra *Hydra) {
	var res acceptLoginResponse
	url := hydra.makeAcceptURL(login, challenge)
	body := map[string]interface{}{
		"subject": subject,
	}
	_ = hydra.putJSON(url, body, &res)
	http.Redirect(w, r, res.RedirectTo, http.StatusFound)
}

type Middleware func(http.Handler) http.Handler

func LoginMiddleware(ab *authboss.Authboss, hydra *Hydra) Middleware {
	return func(handler http.Handler) http.Handler {
		ab.Events.After(authboss.EventAuth, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			user, err := model.GetUser(ab, &r)
			if err != nil {
				return false, err
			}

			// load stored challenge
			challenge := r.FormValue(CTXKeyChallenge)

			acceptLoginRequest(w, r, challenge, user.GetArbitrary()["user_uid"], hydra)

			return true, nil
		})

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				switch r.Method {
				// Show form
				case http.MethodGet:
					if ch := r.URL.Query().Get(hydra.getChallengeName(login)); ch != "" {
						res := initiateLogin(ch, hydra)

						// Skip login when not needed
						if res.Skip {
							acceptLoginRequest(w, r, ch, res.Subject, hydra)
							return
						}

						// Store challenge
						r = r.WithContext(context.WithValue(r.Context(), CTXKeyChallenge, ch))
						if d, ok := r.Context().Value(authboss.CTXKeyData).(authboss.HTMLData); ok {
							r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, d.MergeKV(CTXKeyChallenge, ch)))
						}
					}
				// Evaluate form
				case http.MethodPost:
					if ch := r.FormValue(CTXKeyChallenge); ch != "" {
						if d, ok := r.Context().Value(authboss.CTXKeyData).(authboss.HTMLData); ok {
							r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, d.MergeKV(CTXKeyChallenge, ch)))
						}
					}
				}
			}

			handler.ServeHTTP(w, r)
		})
	}
}
