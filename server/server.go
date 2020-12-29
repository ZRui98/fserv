package server

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/zrui98/fserv/config"
	"github.com/zrui98/fserv/models"
)

type Server struct {
	router    chi.Router
	Config    *config.Config
	users     models.UserRepository
	files     models.FileRepository
	templates map[string]*template.Template
}

func New(c *config.Config) *Server {
	r := chi.NewRouter()
	pool := models.CreatePool(c.DB_URL) // use same pool, with different interfaces
	s := &Server{router: r, Config: c, users: pool, files: pool}
	return s
}

func (s *Server) SetupRoutes() {
	s.router.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			s.renderPage(w, "index.html", nil)
		})
		r.Get("/login", s.LoginGet)
		r.Post("/login", s.LoginPost)
		r.Get("/register", s.RegisterGet)
		r.Post("/register", s.RegisterPost)
		r.Group(func(r chi.Router) {
			r.Use(s.ValidateJwt)
			r.Get("/upload", s.UploadGet)
			r.Post("/upload", s.UploadPost)
			r.Get("/files", s.FilesGet)
			r.Post("/delete", s.FileDelete)
			r.Post("/update", s.FileUpdate)
		})
		r.Get("/400", s.E400)
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
	glog.Info("Server started on PORT 2446")
	glog.Fatal(http.ListenAndServe(":2446", server.router))
}

func (server *Server) ParseTemplates() {
	server.templates = make(map[string]*template.Template)
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			name := info.Name()
			templ := template.Must(template.New(name).Funcs(templateFunctions).ParseFiles("templates/"+name, "templates/head.tmpl", "templates/navbar.tmpl"))
			server.templates[name] = templ
			if err != nil {
				glog.Errorf("Error parsing file: %v\n", err)
			}
		}
		return err
	})
	if err != nil {
		panic(err)
	}
}

func (server *Server) renderPage(w http.ResponseWriter, templateName string, v interface{}) {
	if t, ok := server.templates[templateName]; ok {
		t.Execute(w, &v)
	} else {
		glog.Errorf("Template named %s not found!\n", templateName)
	}
}
