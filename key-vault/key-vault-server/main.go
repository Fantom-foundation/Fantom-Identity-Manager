package main

import (
	"key-vault-server/config"
	"key-vault-server/graphql"
	"key-vault-server/logging"
	"log"
)

var (
	cfg *config.Config
)

func main() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}

	generalLogger := logging.NewLogger(cfg)

	graphql.StartServer(cfg, generalLogger)
}
