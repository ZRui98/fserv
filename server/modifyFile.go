package server

import (
	"net/http"
	"os"

	"github.com/golang/glog"
)

func (s *Server) FileDelete(w http.ResponseWriter, r *http.Request) {
	urlId := r.PostFormValue("fileId")
	glog.Info("Deleting File")
	file, err := s.files.GetFileById(r.Context(), urlId)
	if err != nil {
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	if file.Owner != r.Context().Value(UsernameKey).(string) {
		http.Redirect(w, r, "/400", http.StatusSeeOther)
		return
	}
	err = s.files.DeleteFileById(r.Context(), urlId)
	if err != nil {
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	err = os.Remove(file.FilePath)
	if err != nil {
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/files", http.StatusSeeOther)
}

func (s *Server) FileUpdate(w http.ResponseWriter, r *http.Request) {
	urlId := r.PostFormValue("fileId")
	file, err := s.files.GetFileById(r.Context(), urlId)
	if err != nil {
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	if file.Owner != r.Context().Value(UsernameKey).(string) {
		http.Redirect(w, r, "/400", http.StatusSeeOther)
		return
	}
	isPrivate := r.PostFormValue("private")
	file.IsPrivate = isPrivate == "true"
	glog.Info(isPrivate)
	err = s.files.UpdateFile(r.Context(), file)
	if err != nil {
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/files", http.StatusSeeOther)
}
