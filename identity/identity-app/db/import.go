package db

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/volatiletech/authboss"
)

// User data definition
type ImportedUser struct {
	Id       string `json:"_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Load users from a file and inserts them into the DB
func Import(filename string, db authboss.CreatingServerStorer) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to import users from %s\n", filename)
	}
	defer func() {
		_ = f.Close()
	}()

	var users []ImportedUser
	d := json.NewDecoder(f)
	err = d.Decode(&users)
	if err != nil {
		log.Fatalf("failed to parse users: %v\n", err)
	}

	for _, u := range users {
		user := authboss.MustBeAuthable(db.New(context.Background()))

		user.PutPID(u.Username)
		user.PutPassword(u.Password)
		userRec := authboss.MustBeRecoverable(user)
		userRec.PutEmail(u.Email)

		if arbUser, ok := user.(authboss.ArbitraryUser); ok {
			arbUser.PutArbitrary(map[string]string{
				"user_uid": u.Id,
				"name":     u.Name,
			})
		}

		_ = db.Create(context.Background(), user)
	}
}
