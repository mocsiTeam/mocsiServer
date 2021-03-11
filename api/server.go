package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/mocsiTeam/mocsiServer/api/graph"
	"github.com/mocsiTeam/mocsiServer/api/graph/generated"
	"github.com/mocsiTeam/mocsiServer/auth"
	"github.com/mocsiTeam/mocsiServer/errs"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const defaultPort = "8082"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv.SetRecoverFunc(func(ctx context.Context, e interface{}) error {
		var code int
		var err error
		switch x := e.(type) {
		case string:
			code, err = strconv.Atoi(x)
			if err != nil {
				fmt.Println(err)
				code = 0
			}
		case error:
			code, err = strconv.Atoi(x.Error())
			if err != nil {
				fmt.Println(err)
				code = 0
			}
		default:
			code = 0
		}
		err = &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: errs.CodeErr[uint(code)],
			Extensions: map[string]interface{}{
				"code":   code,
				"status": false,
			},
		}

		return err

	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
