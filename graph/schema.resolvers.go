package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vlamitin/blogg/graph/generated"
	"github.com/vlamitin/blogg/graph/model"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input *model.PostInput) (*model.Post, error) {
	return r.postsRepo.Create(input), nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input *model.PostInput) (*model.Post, error) {
	return r.postsRepo.Save(id, input)
}

func (r *queryResolver) GetPost(ctx context.Context, id string) (*model.Post, error) {
	return r.postsRepo.Get(id)
}

func (r *queryResolver) GetPosts(ctx context.Context, limit int, offset int) ([]*model.Post, error) {
	return r.postsRepo.GetMany(limit, offset)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
