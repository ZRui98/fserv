package server

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/zrui98/fserv/models"
)

type FilesData struct {
	UserFiles []models.File
}

func filesGet(w http.ResponseWriter, r *http.Request, s *Server) (int, []models.File) {
	username := r.Context().Value(usernameKey)
	if username == nil {
		glog.Error("Could not get username fron context, something has gone wrong")
		return http.StatusBadRequest, nil
	}
	files, err := s.files.GetFilesForUser(r.Context(), username.(string))
	if err != nil {
		glog.Error("Connection to DB failed:: %v\n", err)
		return http.StatusInternalServerError, nil
	}
	return http.StatusOK, files
}

func (s *Server) FilesGet(w http.ResponseWriter, r *http.Request) {
	statusCode, files := filesGet(w, r, s)
	if statusCode == 500 {
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	} else if statusCode == 400 {
		http.Redirect(w, r, "/400", http.StatusSeeOther)
		return
	}
	s.renderPage(w, "files.html", &FilesData{
		UserFiles: files,
	})
}

func (s *Server) FilesGetAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
