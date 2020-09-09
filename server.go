package main

import (
	"github.com/vlamitin/blogg/internal/posts"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vlamitin/blogg/graph"
	"github.com/vlamitin/blogg/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	postsRepo := posts.NewRepo()
	postsResolver := graph.NewResolver(postsRepo)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: postsResolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
