package handler

import (
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		h.ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		rpassword := r.FormValue("password1")
		if err := h.Service.Auth.CreateUser(models.User{Email: email, Username: username, Password: password, RepeatPassword: rpassword}); err != nil {
			h.ErrorPage(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	case http.MethodGet:

		if err := h.Temp.ExecuteTemplate(w, "signup.html", nil); err != nil {
			h.ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		h.ErrorPage(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
