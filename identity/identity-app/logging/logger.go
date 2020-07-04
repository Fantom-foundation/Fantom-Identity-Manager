package logging

import (
	"fmt"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	"identity-app/config"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/op/go-logging"
	"github.com/volatiletech/authboss"
)

type Logger struct {
	Logger       logging.Logger
	cfg          *config.Config
	sessionStore *abclientstate.SessionStorer
}

func NewLogger(cfg *config.Config, sessionStore *abclientstate.SessionStorer) *Logger {
	logger := logging.MustGetLogger("identity")
	return &Logger{
		Logger:       *logger,
		cfg:          cfg,
		sessionStore: sessionStore,
	}
}

func (logger *Logger) RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n%s %s %s\n", r.Method, r.URL.Path, r.Proto)

		if logger.cfg.Debug {
			session, err := logger.sessionStore.Get(r, logger.cfg.SessionCookieName)
			if err == nil {
				fmt.Print("Session: ")
				first := true
				for k, v := range session.Values {
					if first {
						first = false
					} else {
						fmt.Print(", ")
					}
					fmt.Printf("%s = %v", k, v)
				}
				fmt.Println()
			}
		}

		if logger.cfg.DebugCTX {
			if val := r.Context().Value(authboss.CTXKeyData); val != nil {
				fmt.Printf("CTX Data: %s", spew.Sdump(val))
			}
			if val := r.Context().Value(authboss.CTXKeyValues); val != nil {
				fmt.Printf("CTX Values: %s", spew.Sdump(val))
			}
		}

		h.ServeHTTP(w, r)
	})
}
