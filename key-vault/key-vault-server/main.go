package main

import (
	"key-vault-server/config"
	"key-vault-server/graphql"
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

	graphql.StartServer(cfg)
}
