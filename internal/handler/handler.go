package handler

import (
	"html/template"
	"net/http"

	"forum/internal/service"
)

type Handler struct {
	Mux     *http.ServeMux
	Temp    *template.Template
	Service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Mux:     http.NewServeMux(),
		Temp:    template.Must(template.ParseGlob("./ui/html/*.html")),
		Service: services,
	}
}

func (h *Handler) InitRoutes() {
	h.Mux.HandleFunc("/", h.middleWareGetUser(h.homePage))

	h.Mux.HandleFunc("/signup", h.signUp)
	h.Mux.HandleFunc("/signin", h.signIn)

	h.Mux.HandleFunc("/post/", h.middleWareGetUser(h.postPage))
	h.Mux.HandleFunc("/post/create", h.middleWareGetUser(h.createPost))
	h.Mux.HandleFunc("/post/myPost", h.middleWareGetUser(h.myPost))
	h.Mux.HandleFunc("/post/myLikedPost", h.middleWareGetUser(h.myLikedPost))

	h.Mux.HandleFunc("/emotion/post/", h.middleWareGetUser(h.emotionPost))
	h.Mux.HandleFunc("/emotion/comment/", h.middleWareGetUser(h.emotionComment))

	h.Mux.HandleFunc("/logout", h.logOut)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	h.Mux.Handle("/static/", http.StripPrefix("/static", fileServer))
}
