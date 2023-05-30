package handler

import (
	"forum/internal/models"
	"log"
	"net/http"
)

func (h *Handler) myPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := r.Context().Value("user").(models.User)
	posts, err := h.Service.GetMyPost(user.Id)
	if err != nil {
		h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info := models.InfoPosts{
		User:     user,
		Posts:    posts,
		Category: nil,
	}

	if err := h.Temp.ExecuteTemplate(w, "myPost.html", info); err != nil {
		log.Println(err.Error())
		h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) myLikedPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := r.Context().Value("user").(models.User)
	posts, err := h.Service.GetMyLikePost(user.Id)
	if err != nil {
		h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info := models.InfoPosts{
		user,
		posts,
		nil,
	}
	if err := h.Temp.ExecuteTemplate(w, "myLikedPost.html", info); err != nil {
		log.Println(err.Error())
		h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}
