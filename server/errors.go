package server

import (
	"html/template"
	"net/http"
)

type ErrorData struct {
	ErrorHeader  string
	ErrorMessage string
}

func (s *Server) E500(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/error.html", "templates/head.tmpl", "templates/navbar.tmpl")
	viewData := &ErrorData{
		ErrorHeader:  "500",
		ErrorMessage: "That's the way she goes. Sometimes she goes, sometimes she doesn't.",
	}
	t.Execute(w, viewData)
}

func (s *Server) E404(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/error.html", "templates/head.tmpl", "templates/navbar.tmpl")
	viewData := &ErrorData{
		ErrorHeader:  "404",
		ErrorMessage: "Where's Randy's BBQ?",
	}
	t.Execute(w, viewData)
}

func (s *Server) E400(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/error.html", "templates/head.tmpl", "templates/navbar.tmpl")
	viewData := &ErrorData{
		ErrorHeader:  "400",
		ErrorMessage: "Frig off!",
	}
	t.Execute(w, viewData)
}
