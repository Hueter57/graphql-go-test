package graph

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Hueter57/graphql-go-test/internal"
	"github.com/Hueter57/graphql-go-test/middlewares/auth"
)

var Directive internal.DirectiveRoot = internal.DirectiveRoot{
	IsAuthenticated: IsAuthenticated,
}

func IsAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	if _, ok := auth.GetUserName(ctx); !ok {
		return nil, errors.New("not authenticated")
	}
	return next(ctx)
}
