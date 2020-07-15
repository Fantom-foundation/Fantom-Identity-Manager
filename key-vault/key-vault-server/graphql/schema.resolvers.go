package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"key-vault-server/graphql/generated"
	"key-vault-server/graphql/model"
)

func (r *queryResolver) Todos(_ context.Context) ([]*model.Todo, error) {
	return r.todos, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
