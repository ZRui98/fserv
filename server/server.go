package server

import (
	"net/http"
	"html/template"
	
	"github.com/zrui98/fserv/models"

	"github.com/go-chi/chi"
	"github.com/golang/glog"
)

type Server struct {
	router chi.Router
	users models.UserRepository
	files models.FileRepository
}

func New() (*Server) {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, nil)
	})
	pool := models.CreatePool() // use same pool, with different interfaces
	s := &Server{ router: r, users: pool, files: pool }
	r.Get("/login", s.LoginGet)
	r.Post("/login", s.LoginPost)
	r.Get("/register", s.RegisterGet)
	r.Post("/register", s.RegisterPost)
	r.Route("/upload", func(cr chi.Router) {
		cr.Use(s.ValidateJwt)
		cr.Get("/", s.UploadGet)
		cr.Post("/", s.UploadPost)
	})
	r.Get("/v/{fileId}", s.ViewFile)
	return s
}

func (server *Server) ListenAndServe() {
	glog.Fatal(http.ListenAndServe(":2446", server.router))
}
