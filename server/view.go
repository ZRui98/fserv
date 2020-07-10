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
	if err != nil || file == nil{
		glog.Errorf("Failed to fetch files:: %v\n", err)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}
	http.ServeFile(w, r, file.FilePath)
}
