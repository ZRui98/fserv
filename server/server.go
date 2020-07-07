package server

import (
	"html/template"
	"net/http"

	"github.com/zrui98/fserv/config"
	"github.com/zrui98/fserv/models"

	"github.com/go-chi/chi"
	"github.com/golang/glog"
)

type Server struct {
	router chi.Router
	Config *config.Config
	users models.UserRepository
	files models.FileRepository
}

func New(c *config.Config) (*Server) {
	r := chi.NewRouter()
	pool := models.CreatePool(c.DB_URL) // use same pool, with different interfaces
	s := &Server{ router: r, Config: c, users: pool, files: pool }
	return s
}

func (s *Server) SetupRoutes() {
	s.router.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			t, err := template.ParseFiles("templates/index.html", "templates/head.tmpl", "templates/navbar.tmpl")
			if err != nil {
				glog.Errorf("Something went wrong parsing template :: %v\n", err)
			}
			t.Execute(w, nil)
		})
		r.Get("/login", s.LoginGet)
		r.Post("/login", s.LoginPost)
		r.Get("/register", s.RegisterGet)
		r.Post("/register", s.RegisterPost)
		r.Group(func(r chi.Router) {
			r.Use(s.ValidateJwt)
			r.Get("/upload", s.UploadGet)
			r.Post("/upload", s.UploadPost)
		})
		r.Get("/404", s.E404)
		r.Get("/500", s.E500)
		r.NotFound(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "/404", http.StatusSeeOther)
		})
		r.Get("/v/{fileId}", s.ViewFile)
		fs := http.FileServer(http.Dir("./static/"))
		r.Handle("/static/*", http.StripPrefix("/static/", fs))
	})
}

func (server *Server) ListenAndServe() {
	glog.Fatal(http.ListenAndServe(":2446", server.router))
}
