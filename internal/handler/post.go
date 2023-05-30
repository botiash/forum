package handler

import (
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	user := r.Context().Value("user").(models.User)
	if !user.IsAuth {
		h.ErrorPage(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		title := r.FormValue("title")
		description := r.FormValue("description")
		categories := r.Form["category"]
		for _, category := range categories {
			if len(category) == 0 || len(category) >= 20 {
				h.ErrorPage(w, "INVALID CATEGORY, category should be shorter than 20 symbols and not empty", http.StatusBadRequest)
				return
			}
		}
		if len(description) > 300 || len(description) == 0 {
			h.ErrorPage(w, "description should be shorter than 300 symbols and not empty", http.StatusBadRequest)
			return
		}
		if len(title) == 0 || len(title) >= 20 {
			h.ErrorPage(w, "INVALID TITLE, title should be shorter than 20 symbols and not empty", http.StatusBadRequest)
			return
		}
		if err := h.Service.ServicePostIR.CreatePost(models.Post{
			Title:       title,
			Description: description,
			Category:    categories,
			Author:      user.Username,
		}); err != nil {
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	case http.MethodGet:
		if err := h.Temp.ExecuteTemplate(w, "postCreate.html", nil); err != nil {
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		h.ErrorPage(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) postPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/" {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if id == 0 || err != nil {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	user := r.Context().Value("user").(models.User)

	post, err := h.Service.ServicePostIR.GetPostId(id)
	if err != nil {
		h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	comments, err := h.Service.GetCommentsByIdPost(id)
	if err != nil {
		h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		model := models.Info{
			User:    user,
			Post:    post,
			Comment: comments,
		}
		if err := h.Temp.ExecuteTemplate(w, "post.html", model); err != nil {
			log.Println(err.Error())
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.ErrorPage(w, "Bad request", http.StatusBadRequest)
			return
		}
		commentText := r.FormValue("text")

		if commentText == "" {
			h.ErrorPage(w, "comment field not found", http.StatusBadRequest)

			return
		}
		if len(commentText) > 300 {
			h.ErrorPage(w, "comment should be shorter than 300 symbols", http.StatusBadRequest)
			return
		}
		if err := h.Service.CommentServiceIR.CreateComment(id, user.Username, commentText); err != nil {
			h.ErrorPage(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, r.URL.Path+fmt.Sprintf("/?id=%d", id), http.StatusSeeOther)
	default:
		h.ErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
