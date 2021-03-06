package posts

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vlamitin/blogg/graph/model"
)

type Repo struct {
	posts         []model.Post
	postAddedChan chan<- model.Post
}

func NewRepo(postAddedChan chan<- model.Post) *Repo {
	return &Repo{
		posts:         []model.Post{},
		postAddedChan: postAddedChan,
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

	select {
	case repo.postAddedChan <- newPost:
	default:
	}

	return &newPost
}

func (repo *Repo) Save(id string, postInput *model.PostInput) (*model.Post, error) {
	var itemToEdit model.Post
	var found = false

	for _, item := range repo.posts {
		if item.ID == id {
			itemToEdit = item
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("post not found with id %s", id)
	}

	itemToEdit.Description = postInput.Description
	itemToEdit.Title = postInput.Title

	return &itemToEdit, nil
}

func (repo *Repo) Get(id string) (*model.Post, error) {
	var requestedItem model.Post
	var found = false

	for _, item := range repo.posts {
		if item.ID == id {
			requestedItem = item
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("post not found with id %s", id)
	}

	return &requestedItem, nil
}

func (repo *Repo) GetMany(limit, offset int) ([]*model.Post, error) {
	res := []*model.Post{}
	for i := offset; i < offset+limit-1; i++ {
		if i > len(repo.posts)-1 {
			break
		}
		res = append(res, &repo.posts[i])
	}

	return res, nil
}
