package login

import (
	"bytes"
	"encoding/json"
	"identity-app/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type flow string

const (
	login   flow = "login"
	consent flow = "consent"
	logout  flow = "logout"
)

type Hydra struct {
	baseURL *url.URL
	client  *http.Client
}

func NewHydra(cfg *config.Config) (*Hydra, error) {
	hydraUrl, err := url.Parse(cfg.HydraAdminUrl)
	if err != nil {
		log.Fatal("Unable to connect to hydra")
	}
	return &Hydra{
		baseURL: hydraUrl,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (h *Hydra) makeGetURL(f flow, challenge string) string {
	return h.makeURL("/oauth2/auth/requests/"+string(f), f, challenge)
}

func (h *Hydra) makeAcceptURL(f flow, challenge string) string {
	return h.makeURL("/oauth2/auth/requests/"+string(f)+"/accept", f, challenge)
}

func (h *Hydra) makeURL(path string, f flow, challenge string) string {
	p, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	u := h.baseURL.ResolveReference(p)

	q := u.Query()
	q.Set(h.getChallengeName(f), challenge)
	u.RawQuery = q.Encode()

	return u.String()
}

func (h *Hydra) getChallengeName(f flow) string {
	return string(f) + "_challenge"
}

func (h *Hydra) getJSON(url string, target interface{}) error {
	res, err := h.client.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	return json.NewDecoder(res.Body).Decode(target)
}

func (h *Hydra) putJSON(url string, body interface{}, target interface{}) error {
	var b io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		b = bytes.NewBuffer(jsonBody)
	}
	req, _ := http.NewRequest(http.MethodPut, url, b)

	res, err := h.client.Do(req)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	return json.NewDecoder(res.Body).Decode(target)
}
