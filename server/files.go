package server

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/zrui98/fserv/models"
)

type FilesData struct {
	UserFiles []models.File
}

func (s *Server) FilesGet(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(UsernameKey)
	if username == nil {
		glog.Error("Could not get username fron context, something has gone wrong")
	}
	files, err := s.files.GetFilesForUser(r.Context(), username.(string))
	if err != nil {
		glog.Error("Connection to DB failed:: %v\n", err)
	}
	s.renderPage(w, "files.html", &FilesData{
		UserFiles: files,
	})
}
