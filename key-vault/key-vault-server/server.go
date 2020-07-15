package main

import (
	"key-vault-server/config"
	"key-vault-server/graphql"
	"key-vault-server/graphql/generated"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graphql.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to %s:%s/ for GraphQL playground", cfg.RootUrl, cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
