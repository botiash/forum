package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	if r.URL.Path != "/" {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	user := r.Context().Value("user").(models.User)
	if r.Method != http.MethodGet {
		h.ErrorPage(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	categories, err := h.Service.GetCategories()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if r.URL.Query().Has("category") {
		category := r.URL.Query().Get("category")
		if !inSlice(category, categories) {
			h.ErrorPage(w, "Not exist", http.StatusBadRequest)
		}
		posts, err = h.Service.ServicePostIR.GetAllPostsByCategories(category)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		posts, err = h.Service.ServicePostIR.GetAllPosts()
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	info := models.InfoPosts{
		user,
		posts,
		categories,
	}
	if err := h.Temp.ExecuteTemplate(w, "homepage.html", info); err != nil {
		log.Println(err.Error())
		h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
	}
}

func inSlice(val string, slice []string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
