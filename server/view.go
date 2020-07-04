package server

import (
	"net/http"
	"github.com/golang/glog"
	"github.com/go-chi/chi"
)

func (server *Server) ViewFile(w http.ResponseWriter, r *http.Request) {
	url_id := chi.URLParam(r, "fileId")
	glog.Info("Getting stuff")
	file, err := server.files.GetFileById(r.Context(), url_id)
	if err != nil {
		glog.Errorf("Failed to fetch files:: %v\n", err)
	}
	http.ServeFile(w, r, file.FilePath)
}
