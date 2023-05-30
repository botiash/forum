package service

import (
	"forum/internal/models"
	"forum/internal/storage"
)

type EmotionServiceIR interface {
	CreateOrUpdateEmotionComment(models.Like) error
	CreateOrUpdateEmotionPost(models.Like) error
}

type EmotionService struct {
	storage storage.ReactionIR
}

func NewEmotionService(storage storage.ReactionIR) EmotionServiceIR {
	return &EmotionService{
		storage,
	}
}

func (e *EmotionService) CreateOrUpdateEmotionPost(postEmo models.Like) error {
	exists1, err := e.storage.EmotionPostExistsFull(postEmo)
	if err != nil {
		return err
	}
	if exists1 {
		postEmo.Islike = -1
		return e.storage.UpdateEmotionPost(postEmo)
	}

	exists, err := e.storage.EmotionPostExists(postEmo.PostID, postEmo.UserID)
	if err != nil {
		return err
	}
	if exists {
		err = e.storage.UpdateEmotionPost(postEmo)
		if err != nil {
			return err
		}
	} else {
		err = e.storage.CreateEmotionPost(postEmo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EmotionService) CreateOrUpdateEmotionComment(commentEmo models.Like) error {
	exists1, err := e.storage.EmotionCommentExistsFull(commentEmo)
	if err != nil {
		return err
	}
	if exists1 {
		commentEmo.Islike = -1
		return e.storage.UpdateEmotionComment(commentEmo)
	}

	exists, err := e.storage.EmotionCommentExists(commentEmo.CommentID, commentEmo.UserID)
	if err != nil {
		return err
	}
	if exists {
		err = e.storage.UpdateEmotionComment(commentEmo)
		if err != nil {
			return err
		}
	} else {
		err = e.storage.CreateEmotionComment(commentEmo)
		if err != nil {
			return err
		}
	}

	return nil
}
