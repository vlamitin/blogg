package graph

//go:generate go run github.com/99designs/gqlgen
import "github.com/vlamitin/blogg/internal/posts"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	postsRepo *posts.Repo
}

func NewResolver(postsRepo *posts.Repo) *Resolver {
	return &Resolver{postsRepo: postsRepo}
}
