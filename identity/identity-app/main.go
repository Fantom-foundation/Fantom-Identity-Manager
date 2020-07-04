package main

import (
	"context"
	abrenderer "github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/defaults"
	"identity-app/config"
	"identity-app/db"
	"identity-app/login"
	"identity-app/model"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/justinas/nosurf"

	"github.com/volatiletech/authboss"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	_ "github.com/volatiletech/authboss/auth"
	_ "github.com/volatiletech/authboss/logout"

	//_ "github.com/volatiletech/authboss/remember" // uncomment to enable remembering functionality
	//_ "github.com/volatiletech/authboss/recover" // uncomment to enable user recovery functionality
	_ "github.com/volatiletech/authboss/register"
)

const (
	sessionCookieName = "fantomX"
)

var (
	cfg *config.Config
)

var (
	ab        = authboss.New()
	database  = db.NewMemStorer()
	schemaDec = schema.NewDecoder()

	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
)

func authbossSetup() {
	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	ab.Config.Paths.Mount = "/auth"
	ab.Config.Core.ViewRenderer = abrenderer.NewHTML(ab.Config.Paths.Mount, "ab_views")
	ab.Config.Modules.LogoutMethod = http.MethodGet
	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name", "user_uid"}

	defaults.SetCore(&ab.Config, false, false)

	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]+`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength:  4,
		MinSymbols: 0,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength:       2,
		MinSymbols:      0,
		AllowWhitespace: true,
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		UseUsername: false,
		ReadJSON:    false,
		Rulesets: map[string][]defaults.Rules{
			"register": {emailRule, passwordRule, nameRule},
		},
		Confirms: map[string][]string{
			"register": {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": {"email", "name", "password", "user_uid"},
		},
	}
}

func main() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}

	cookieStore = abclientstate.NewCookieStorer([]byte(os.Getenv("COOKIE_STORE_KEY")), nil)
	cookieStore.Secure = false
	sessionStore = abclientstate.NewSessionStorer(sessionCookieName, []byte(os.Getenv("SESSION_STORE_KEY")), nil)

	cStore := sessionStore.Store.(*sessions.CookieStore)
	cStore.Options.Secure = false
	cStore.MaxAge(int((7 * 24 * time.Hour) / time.Second))

	authbossSetup()

	if filename := os.Getenv("IMPORT_USERS"); filename != "" {
		log.Printf("Importing users from file: %s\n", filename)
		db.Import(filename, database)
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	rootURL := os.Getenv("ROOT_URL")
	if rootURL == "" {
		rootURL = "http://localhost:" + port
	}
	_, err = url.Parse(rootURL)
	if err != nil {
		panic("invalid root URL passed")
	}
	ab.Config.Paths.RootURL = rootURL
	ab.Config.Paths.RegisterOK = rootURL

	if err := ab.Init(); err != nil {
		panic(err)
	}

	schemaDec.IgnoreUnknownKeys(true)

	mux := chi.NewRouter()

	mux.Use(logger,
		nosurf.NewPure,
		ab.LoadClientStateMiddleware,
		dataInjector,
		authboss.ModuleListMiddleware(ab),
	)

	mux.Route(ab.Config.Paths.Mount, func(mux chi.Router) {
		mws := chi.Chain(
			login.LoginMiddleware(ab),
			login.LogoutMiddleware(ab),
			login.RegisterMiddleware(ab),
		)
		mux.Mount("/", http.StripPrefix(ab.Config.Paths.Mount, mws.Handler(ab.Config.Core.Router)))
		mux.Mount("/consent", login.Consent(ab))

		fs := http.FileServer(http.Dir("static"))
		mux.Mount("/static/", http.StripPrefix(ab.Config.Paths.Mount+"/static/", fs))
	})

	log.Printf("Listening on port %s", port)
	log.Println(http.ListenAndServe(":"+port, mux))
}

func dataInjector(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := layoutData(w, &r)
		r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, data))
		handler.ServeHTTP(w, r)
	})
}

// layoutData is passing pointers to pointers be able to edit the current pointer
// to the request. This is still safe as it still creates a new request and doesn't
// modify the old one, it just modifies what we're pointing to in our methods so
// we're able to skip returning an *http.Request everywhere
func layoutData(w http.ResponseWriter, r **http.Request) authboss.HTMLData {
	var loggedIn bool
	var currentUserName string

	if user, err := model.GetUser(ab, r); user != nil && err == nil {
		loggedIn = true
		currentUserName = user.Name
	}

	return authboss.HTMLData{
		"loggedin":          loggedIn,
		"current_user_name": currentUserName,
		"csrf_token":        nosurf.Token(*r),
		"flash_success":     authboss.FlashSuccess(w, *r),
		"flash_error":       authboss.FlashError(w, *r),
	}
}
