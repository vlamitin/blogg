package main

import (
	"github.com/gorilla/websocket"
	"github.com/vlamitin/blogg/internal/posts"
	"github.com/vlamitin/blogg/internal/subscribers"
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

	subs := subscribers.NewSubscribers()
	postsRepo := posts.NewRepo(subs.Dispatch)

	subscribeHandler := getSubscribeHandler(subs.AddSubscriber)

	postsResolver := graph.NewResolver(postsRepo)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: postsResolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.HandleFunc("/subscribe", subscribeHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getSubscribeHandler(onConnection func(sub *websocket.Conn)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		onConnection(c)
		c.WriteMessage(websocket.TextMessage, []byte("subscribed"))
	}
}
