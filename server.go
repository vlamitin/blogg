package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vlamitin/blogg/graph/model"
	"github.com/vlamitin/blogg/internal/notifier"
	"github.com/vlamitin/blogg/internal/posts"

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

	newPostsChan := make(chan model.Post)
	newPostsMessageCh := make(chan string)

	go func() {
		for {
			newPost := <-newPostsChan
			newPostsMessageCh <- fmt.Sprintf("New post! Title: '%s'", newPost.Title)
		}
	}()

	postsRepo := posts.NewRepo(newPostsChan)
	postsNotifier := notifier.NewNotifier(newPostsMessageCh)
	postsSubsHandler := notifier.NewWsHandler(postsNotifier)
	go postsNotifier.Run()

	postsResolver := graph.NewResolver(postsRepo)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: postsResolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.Handle("/subscribe", postsSubsHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
