package server

import (
	"net/http"
	"os"

	"github.com/golang/glog"
)

func (s *Server) fileDelete(w http.ResponseWriter, r *http.Request) int {
	urlID := r.PostFormValue("fileId")
	glog.Info("Deleting File")
	file, err := s.files.GetFileById(r.Context(), urlID)
	if err != nil {
		return http.StatusInternalServerError
	}
	if file.Owner != r.Context().Value(usernameKey).(string) {
		return http.StatusBadRequest
	}
	err = s.files.DeleteFileById(r.Context(), urlID)
	if err != nil {
		return http.StatusBadRequest
	}
	err = os.Remove(file.FilePath)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (s *Server) FileDelete(w http.ResponseWriter, r *http.Request) {
	statusCode := s.fileDelete(w, r)
	if statusCode == 500 {
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	} else if statusCode == 400 {
		http.Redirect(w, r, "/400", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/files", http.StatusSeeOther)
}

func (s *Server) FileDeleteAPI(w http.ResponseWriter, r *http.Request) {
	statusCode := s.fileDelete(w, r)
	w.WriteHeader(statusCode)
}

func (s *Server) fileUpdate(w http.ResponseWriter, r *http.Request) int {
	urlId := r.PostFormValue("fileId")
	file, err := s.files.GetFileById(r.Context(), urlId)
	if err != nil {
		return http.StatusInternalServerError
	}
	if file.Owner != r.Context().Value(usernameKey).(string) {
		return http.StatusBadRequest
	}
	isPrivate := r.PostFormValue("private")
	file.IsPrivate = isPrivate == "true"
	glog.Info(isPrivate)
	err = s.files.UpdateFile(r.Context(), file)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (s *Server) FileUpdate(w http.ResponseWriter, r *http.Request) {
	statusCode := s.fileDelete(w, r)
	if statusCode == 500 {
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	} else if statusCode == 400 {
		http.Redirect(w, r, "/400", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/files", http.StatusSeeOther)
}

func (s *Server) FileUpdateAPI(w http.ResponseWriter, r *http.Request) {
	statusCode := s.fileUpdate(w, r)
	w.WriteHeader(statusCode)
}
