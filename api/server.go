package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"

	"github.com/mocsiTeam/mocsiServer/api/graph"
	"github.com/mocsiTeam/mocsiServer/api/graph/generated"
	"github.com/mocsiTeam/mocsiServer/errs"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func NewGraphQLServer() *handler.Server {
	var resolver graph.Resolver
	resolver.InitResolver()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}))
	srv.SetRecoverFunc(func(ctx context.Context, e interface{}) error {
		return recoverFunc(ctx, e)
	})

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}

func recoverFunc(ctx context.Context, e interface{}) error {
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

}
