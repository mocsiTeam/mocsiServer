package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/go-chi/chi"

	"github.com/mocsiTeam/mocsiServer/api"
	"github.com/mocsiTeam/mocsiServer/auth"
	"github.com/mocsiTeam/mocsiServer/db"
)

const defaultPort = "8082"

func init() {
	err := db.CreateDB()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := api.NewGraphQLServer()

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
