package service

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/storage"
	"strings"
)

type CommentServiceIR interface {
	GetCommentsByIdPost(id int) ([]models.Comment, error)
	CreateComment(id int, username, commentText string) error
}

type CommentService struct {
	storage storage.CommentIR
}

func newCommentServ(storage storage.CommentIR) CommentServiceIR {
	return &CommentService{
		storage: storage,
	}
}

func (c *CommentService) GetCommentsByIdPost(id int) ([]models.Comment, error) {
	return c.storage.GetCommentsByIdPost(id)
}

func (c *CommentService) CreateComment(id int, username, commentText string) error {
	commentText = strings.TrimSpace(commentText)
	if len(commentText) == 0 {
		return fmt.Errorf("Empty comment")
	}
	return c.storage.CreateComment(id, username, commentText)
}
