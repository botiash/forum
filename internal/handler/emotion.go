package handler

import (
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) emotionComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/emotion/comment/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	commentId, err := strconv.Atoi(r.URL.Query().Get("id"))
	postId, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if commentId == 0 || err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	user := r.Context().Value("user").(models.User)
	if user.Username == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		h.ErrorPage(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	res := r.Form.Get("islike")
	if res == "" {
		h.ErrorPage(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if res == "like" {
		err = h.Service.EmotionServiceIR.CreateOrUpdateEmotionComment(models.Like{UserID: user.Id, CommentID: commentId, Islike: 1})
	} else if res == "dislike" {
		err = h.Service.EmotionServiceIR.CreateOrUpdateEmotionComment(models.Like{UserID: user.Id, CommentID: commentId, Islike: 0})
	} else {
		h.ErrorPage(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err != nil {
		h.ErrorPage(w, err.Error(), http.StatusBadRequest)
	}
	link := fmt.Sprintf("/post/?id=%d", postId)
	http.Redirect(w, r, link, http.StatusSeeOther)
}

func (h *Handler) emotionPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/emotion/post/" {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	postId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if postId == 0 || err != nil {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	user := r.Context().Value("user").(models.User)
	if user.Username == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		h.ErrorPage(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	res := r.Form.Get("islike")
	if res == "" {
		h.ErrorPage(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if res == "like" {
		err = h.Service.EmotionServiceIR.CreateOrUpdateEmotionPost(models.Like{UserID: user.Id, PostID: postId, Islike: 1})
	} else if res == "dislike" {
		err = h.Service.EmotionServiceIR.CreateOrUpdateEmotionPost(models.Like{UserID: user.Id, PostID: postId, Islike: 0})
	} else {
		h.ErrorPage(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err != nil {
		h.ErrorPage(w, err.Error(), http.StatusBadRequest)
	}
	link := fmt.Sprintf("/post/?id=%d", postId)
	http.Redirect(w, r, link, http.StatusSeeOther)
}
