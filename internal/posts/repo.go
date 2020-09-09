package posts

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vlamitin/blogg/graph/model"
	"time"
)

type PostsSubsciber func(post *model.Post)

type Repo struct {
	posts       []model.Post
	onPostAdded func(title string, description string)
}

func NewRepo(onPostAdded func(title string, description string)) *Repo {
	return &Repo{
		posts:       []model.Post{},
		onPostAdded: onPostAdded,
	}
}

func (repo *Repo) Create(postInput *model.PostInput) *model.Post {
	newPost := model.Post{
		ID:              uuid.New().String(),
		Title:           postInput.Title,
		Description:     postInput.Description,
		PublicationDate: int(time.Now().UnixNano()),
	}

	repo.posts = append(repo.posts, newPost)
	repo.onPostAdded(newPost.Title, newPost.Description)

	return &newPost
}

func (repo *Repo) Save(id string, postInput *model.PostInput) (*model.Post, error) {
	var itemToEdit model.Post

	for _, item := range repo.posts {
		if item.ID == id {
			itemToEdit = item
			break
		}

		return nil, fmt.Errorf("post not found with id %s", id)
	}

	itemToEdit.Description = postInput.Description
	itemToEdit.Title = postInput.Title

	return &itemToEdit, nil
}

func (repo *Repo) Get(id string) (*model.Post, error) {
	var requestedItem model.Post

	for _, item := range repo.posts {
		if item.ID == id {
			requestedItem = item
			break
		}
		return nil, fmt.Errorf("post not found with id %s", id)
	}

	return &requestedItem, nil
}

func (repo *Repo) GetMany(limit int, offset int) ([]*model.Post, error) {
	res := []*model.Post{}
	for i := offset; i < offset+limit-1; i++ {
		if i > len(repo.posts)-1 {
			break
		}
		res = append(res, &repo.posts[i])
	}

	return res, nil
}
