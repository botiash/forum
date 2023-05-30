package service

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/storage"
	"strings"
)

type ServicePostIR interface {
	CreatePost(post models.Post) error
	GetPostId(id int) (models.Post, error)
	GetAllPosts() ([]models.Post, error)
	GetCategories() ([]string, error)
	GetAllPostsByCategories(category string) ([]models.Post, error)
	GetMyPost(int) ([]models.Post, error)
	GetMyLikePost(int) ([]models.Post, error)
}

type PostService struct {
	storage storage.PostIR
}

func NewPostService(postIR storage.PostIR) ServicePostIR {
	return &PostService{
		storage: postIR,
	}
}

func (p *PostService) CreatePost(post models.Post) error {
	for x := range post.Category {
		post.Category[x] = strings.TrimSpace(post.Category[x])
		if len(post.Category[x]) == 0 {
			return fmt.Errorf("Empty category")
		}
	}
	post.Title = strings.TrimSpace(post.Title)
	if len(post.Title) == 0 {
		return fmt.Errorf("Empty Description")
	}
	post.Description = strings.TrimSpace(post.Description)
	if len(post.Description) == 0 {
		return fmt.Errorf("Empty Description")
	}

	return p.storage.CreatePost(post)
}

func (p *PostService) GetPostId(id int) (models.Post, error) {
	return p.storage.GetPostByID(id)
}

func (p *PostService) GetAllPosts() ([]models.Post, error) {
	return p.storage.GetAllPosts()
}

func (p *PostService) GetCategories() ([]string, error) {
	return p.storage.Category()
}

func (p *PostService) GetAllPostsByCategories(category string) ([]models.Post, error) {
	return p.storage.GetAllPostsByCategories(category)
}

func (p *PostService) GetMyPost(id int) ([]models.Post, error) {
	return p.storage.GetMyPost(id)
}

func (p *PostService) GetMyLikePost(id int) ([]models.Post, error) {
	return p.storage.GetMyLikedPost(id)
}
