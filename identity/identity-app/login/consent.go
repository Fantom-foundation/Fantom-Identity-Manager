package login

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/volatiletech/authboss"
	"identity-app/model"
)

type consentResponse struct {
	Skip                         bool     `json:"skip"`
	RequestedScope               []string `json:"requested_scope"`
	RequestedAccessTokenAudience []string `json:"requested_access_token_audience"`
}

func initiateConsent(challenge string, hydra *Hydra) consentResponse {
	var res consentResponse
	url := hydra.makeGetURL(consent, challenge)
	_ = hydra.getJSON(url, &res)

	return res
}

type acceptConsentResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func acceptConsentRequest(w http.ResponseWriter, r *http.Request, challenge string, getRes consentResponse, user *model.User, hydra *Hydra) {
	var res acceptConsentResponse
	url := hydra.makeAcceptURL(consent, challenge)
	var idToken = IDToken{
		Uid:  user.UserUid,
		Name: user.Name,
	}
	body := map[string]interface{}{
		"grant_scope":                 getRes.RequestedScope,
		"grant_access_token_audience": getRes.RequestedAccessTokenAudience,
		"session": map[string]interface{}{
			"id_token": idToken,
		},
	}
	_ = hydra.putJSON(url, body, &res)
	http.Redirect(w, r, res.RedirectTo, http.StatusFound)
}

type IDToken struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

func Consent(ab *authboss.Authboss, hydra *Hydra) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if ch := r.URL.Query().Get(hydra.getChallengeName(consent)); ch != "" {
			// Auto consent to every request
			getRes := initiateConsent(ch, hydra)
			if user, err := model.GetUser(ab, &r); err == nil {
				acceptConsentRequest(w, r, ch, getRes, user, hydra)
			}
		}
	})

	return mux
}
