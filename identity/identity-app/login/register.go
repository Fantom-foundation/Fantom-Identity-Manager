package login

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/authboss"
	"identity-app/model"
	"net/http"
)

func RegisterMiddleware(ab *authboss.Authboss) Middleware {
	return func(handler http.Handler) http.Handler {
		ab.Events.After(authboss.EventRegister, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
			user, err := model.GetUser(ab, &r)
			if err != nil {
				return false, err
			}

			// load stored challenge
			challenge := r.Context().Value(CTXKeyChallenge).(string)

			acceptLoginRequest(w, r, challenge, user.GetArbitrary()["user_uid"])

			return true, nil
		})

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/register" {
				switch r.Method {
				// Show form
				case http.MethodGet:
					if ch := r.URL.Query().Get(getChallengeName(login)); ch != "" {
						res := initiateLogin(ch)

						// Skip login when not needed
						if res.Skip {
							acceptLoginRequest(w, r, ch, res.Subject)
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
						r = r.WithContext(context.WithValue(r.Context(), CTXKeyChallenge, ch))
					}

					// Add generated values to saved user
					r.Form.Set("user_uid", uuid.NewV4().String())
				}
			}

			handler.ServeHTTP(w, r)
		})
	}
}
