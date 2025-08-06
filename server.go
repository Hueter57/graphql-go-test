package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aarondl/sqlboiler/v4/boil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Hueter57/graphql-go-test/graph"
	"github.com/Hueter57/graphql-go-test/graph/resolver"
	"github.com/Hueter57/graphql-go-test/graph/services"
	"github.com/Hueter57/graphql-go-test/internal"
	"github.com/Hueter57/graphql-go-test/middlewares/auth"
)

const (
	defaultPort = "8080"
	dbFile      = "./mygraphql.db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=on", dbFile))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	service := services.New(db)

	srv := handler.New(internal.NewExecutableSchema(internal.Config{
		Resolvers: &resolver.Resolver{
			Srv:     service,
			Loaders: graph.NewLoaders(service),
		},
		Directives: graph.Directive,
		Complexity: graph.ComplexityConfig(),
	}))
	srv.Use(extension.FixedComplexityLimit(20))

	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		log.Println("before OperationHandler")
		res := next(ctx)
		defer log.Println("after OperationHandler")
		return res
	})
	srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		log.Println("before ResponseHandler")
		res := next(ctx)
		defer log.Println("after ResponseHandler")
		return res
	})
	srv.AroundRootFields(func(ctx context.Context, next graphql.RootResolver) graphql.Marshaler {
		log.Println("before RootResolver")
		res := next(ctx)
		defer func() {
			var b bytes.Buffer
			res.MarshalGQL(&b)
			log.Println("after RootResolver", b.String())
		}()
		return res
	})
	srv.AroundFields(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
		log.Println("before Resolver")
		res, err = next(ctx)
		defer log.Println("after Resolver", res)
		return
	})

	boil.DebugMode = true

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", auth.AuthMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
