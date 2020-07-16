package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"key-vault-server/auth"
	"key-vault-server/config"
	"key-vault-server/graphql/generated"
	"key-vault-server/logging"
	"log"
	"net/http"
)

func StartServer(cfg *config.Config, logger *logging.Logger) {
	router := chi.NewRouter()
	// Setup user Authorization service
	router.Use(auth.AuthenticationMiddleware(cfg))

	// Main GraphQl resolver
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))

	// Define routes
	router.Handle("/", playground.Handler("GraphQL playground", cfg.GraphqlEntrypoint))
	router.Handle(cfg.GraphqlEntrypoint, srv)

	// Start server
	log.Printf("connect to %s:%s/ for GraphQL playground", cfg.RootUrl, cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
