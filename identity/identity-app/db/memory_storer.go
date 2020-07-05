package db

import (
	"context"
	"fmt"
	"identity-app/config"
	"identity-app/logging"
	"identity-app/model"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/volatiletech/authboss"
	aboauth "github.com/volatiletech/authboss/oauth2"
)

const (
	memoryStorerIdentifier = "memory"
)

var (
	memoryStorer = &MemStorer{}

	_ StorerBase                       = memoryStorer
	_ authboss.CreatingServerStorer    = memoryStorer
	_ authboss.ConfirmingServerStorer  = memoryStorer
	_ authboss.RecoveringServerStorer  = memoryStorer
	_ authboss.RememberingServerStorer = memoryStorer
)

// MemStorer stores users in memory
type MemStorer struct {
	Users  map[string]model.User
	Tokens map[string][]string
}

func init() {
	RegisterStorer(func() StorerBase {
		return NewMemStorer()
	})
}

// NewMemStorer constructor
func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users:  map[string]model.User{},
		Tokens: make(map[string][]string),
	}
}

// Nothing needed to do
func (storer *MemStorer) Close() {
}

func (storer *MemStorer) CanHandle(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	return scheme == memoryStorerIdentifier
}

func (storer *MemStorer) FromConfig(_ *config.Config, _ *logging.Logger) (StorerBase, error) {
	return storer, nil
}

// Save the user
func (storer MemStorer) Save(_ context.Context, user authboss.User) error {
	u := user.(*model.User)
	storer.Users[u.Email] = *u

	fmt.Println("Saved user:", u.Name)
	return nil
}

// Load the user
func (storer MemStorer) Load(_ context.Context, key string) (user authboss.User, err error) {
	// Check to see if our key is actually an oauth2 pid
	provider, uid, err := authboss.ParseOAuth2PID(key)
	if err == nil {
		for _, u := range storer.Users {
			if u.OAuth2Provider == provider && u.OAuth2UID == uid {
				fmt.Println("Loaded OAuth2 user:", u.Email)
				return &u, nil
			}
		}

		return nil, authboss.ErrUserNotFound
	}

	u, ok := storer.Users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	fmt.Println("Loaded user:", u.Name)
	return &u, nil
}

// New user creation
func (storer MemStorer) New(_ context.Context) authboss.User {
	return &model.User{}
}

// Create the user
func (storer MemStorer) Create(_ context.Context, user authboss.User) error {
	u := user.(*model.User)

	if _, ok := storer.Users[u.Email]; ok {
		return authboss.ErrUserFound
	}

	fmt.Println("Created new user:", u.Name)
	storer.Users[u.Email] = *u
	return nil
}

// LoadByConfirmSelector looks a user up by confirmation token
func (storer MemStorer) LoadByConfirmSelector(_ context.Context, selector string) (user authboss.ConfirmableUser, err error) {
	for _, v := range storer.Users {
		if v.ConfirmSelector == selector {
			fmt.Println("Loaded user by confirm selector:", selector, v.Name)
			return &v, nil
		}
	}

	return nil, authboss.ErrUserNotFound
}

// LoadByRecoverSelector looks a user up by confirmation selector
func (storer MemStorer) LoadByRecoverSelector(_ context.Context, selector string) (user authboss.RecoverableUser, err error) {
	for _, v := range storer.Users {
		if v.RecoverSelector == selector {
			fmt.Println("Loaded user by recover selector:", selector, v.Name)
			return &v, nil
		}
	}

	return nil, authboss.ErrUserNotFound
}

// AddRememberToken to a user
func (storer MemStorer) AddRememberToken(_ context.Context, pid, token string) error {
	storer.Tokens[pid] = append(storer.Tokens[pid], token)
	fmt.Printf("Adding rm token to %s: %s\n", pid, token)
	spew.Dump(storer.Tokens)
	return nil
}

// DelRememberTokens removes all tokens for the given pid
func (storer MemStorer) DelRememberTokens(_ context.Context, pid string) error {
	delete(storer.Tokens, pid)
	fmt.Println("Deleting rm tokens from:", pid)
	spew.Dump(storer.Tokens)
	return nil
}

// UseRememberToken finds the pid-token pair and deletes it.
// If the token could not be found return ErrTokenNotFound
func (storer MemStorer) UseRememberToken(_ context.Context, pid, token string) error {
	tokens, ok := storer.Tokens[pid]
	if !ok {
		fmt.Println("Failed to find rm tokens for:", pid)
		return authboss.ErrTokenNotFound
	}

	for i, tok := range tokens {
		if tok == token {
			tokens[len(tokens)-1] = tokens[i]
			storer.Tokens[pid] = tokens[:len(tokens)-1]
			fmt.Printf("Used remember for %s: %s\n", pid, token)
			return nil
		}
	}

	return authboss.ErrTokenNotFound
}

// NewFromOAuth2 creates an oauth2 user (but not in the database, just a blank one to be saved later)
func (storer MemStorer) NewFromOAuth2(_ context.Context, provider string, details map[string]string) (authboss.OAuth2User, error) {
	switch provider {
	case "google":
		email := details[aboauth.OAuth2Email]

		var user *model.User
		if u, ok := storer.Users[email]; ok {
			user = &u
		} else {
			user = &model.User{}
		}

		// Google OAuth2 doesn't allow us to fetch real name without more complicated API calls
		// in order to do this properly in your own identity-app, look at replacing the authboss oauth2.GoogleUserDetails
		// method with something more thorough.
		user.Name = "Unknown"
		user.Email = details[aboauth.OAuth2Email]
		user.OAuth2UID = details[aboauth.OAuth2UID]
		user.Confirmed = true

		return user, nil
	}

	return nil, errors.Errorf("unknown provider %s", provider)
}

// SaveOAuth2 user
func (storer MemStorer) SaveOAuth2(_ context.Context, user authboss.OAuth2User) error {
	u := user.(*model.User)
	storer.Users[u.Email] = *u

	return nil
}

/*
func (s MemStorer) PutOAuth(uid, provider string, attr authboss.Attributes) error {
	return s.Create(uid+provider, attr)
}

func (s MemStorer) GetOAuth(uid, provider string) (result interface{}, err error) {
	user, ok := s.Users[uid+provider]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &user, nil
}

func (s MemStorer) AddToken(key, token string) error {
	s.Tokens[key] = append(s.Tokens[key], token)
	fmt.Println("AddToken")
	spew.Dump(s.Tokens)
	return nil
}

func (s MemStorer) DelTokens(key string) error {
	delete(s.Tokens, key)
	fmt.Println("DelTokens")
	spew.Dump(s.Tokens)
	return nil
}

func (s MemStorer) UseToken(givenKey, token string) error {
	toks, ok := s.Tokens[givenKey]
	if !ok {
		return authboss.ErrTokenNotFound
	}

	for i, tok := range toks {
		if tok == token {
			toks[i], toks[len(toks)-1] = toks[len(toks)-1], toks[i]
			s.Tokens[givenKey] = toks[:len(toks)-1]
			return nil
		}
	}

	return authboss.ErrTokenNotFound
}

func (s MemStorer) ConfirmUser(tok string) (result interface{}, err error) {
	fmt.Println("==============", tok)

	for _, u := range s.Users {
		if u.ConfirmToken == tok {
			return &u, nil
		}
	}

	return nil, authboss.ErrUserNotFound
}

func (s MemStorer) RecoverUser(rec string) (result interface{}, err error) {
	for _, u := range s.Users {
		if u.RecoverToken == rec {
			return &u, nil
		}
	}

	return nil, authboss.ErrUserNotFound
}
*/
