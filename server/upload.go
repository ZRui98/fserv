package server

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gabriel-vasile/mimetype"
	"github.com/golang/glog"
	"github.com/zrui98/fserv/constants"
	"github.com/zrui98/fserv/models"
	"github.com/zrui98/fserv/utils"
)

func (s *Server) uploadPost(w http.ResponseWriter, r *http.Request) int {
	glog.Info("Uploading File")
	r.ParseMultipartForm(32 << 20)
	isPrivate := r.FormValue("private") == "true"
	uploadedFile, handler, err := r.FormFile("uploadFile")
	if err != nil {
		glog.Errorf("Error parsing form file:: %v\n", err)
		return http.StatusInternalServerError
	}
	defer uploadedFile.Close()

	mime, err := mimetype.DetectReader(uploadedFile)
	if err != nil {
		glog.Errorf("Error unrecognized mimetype:: %v\n", err)
		return http.StatusInternalServerError
	}
	folder := "txt"
	if mimetype.EqualsAny(mime.String(), constants.Image...) {
		folder = "images"
	} else if mimetype.EqualsAny(mime.String(), constants.Video...) {
		folder = "videos"
	}
	token := utils.GenerateFileName()
	filePath := path.Join(
		s.Config.ROOT_DIR,
		r.Context().Value(usernameKey).(string),
		folder,
		token+path.Ext(handler.Filename),
	)
	glog.Infof("Saving to: %s File: %s\n", filePath, handler.Filename)
	file := &models.File{
		UrlId:     token,
		FilePath:  filePath,
		FileName:  handler.Filename,
		Owner:     r.Context().Value(usernameKey).(string),
		IsPrivate: isPrivate,
	}
	err = s.files.AddFile(r.Context(), file)
	if err != nil {
		glog.Errorf("Error inserting data into DB: %v\n", err)
		return http.StatusInternalServerError
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		glog.Errorf("Error creating file pointer: %v\n", err)
		return http.StatusInternalServerError
	}
	defer f.Close()
	if _, err := uploadedFile.Seek(0, 0); err != nil {
		glog.Errorf("Error seeking uploaded file: %v\n", err)
		return http.StatusInternalServerError
	}
	io.Copy(f, uploadedFile)
	return http.StatusOK
}

// UploadPost returns website call for uploading file
func (s *Server) UploadPost(w http.ResponseWriter, r *http.Request) {
	statusCode := s.uploadPost(w, r)
	if statusCode == http.StatusOK {
		s.renderPage(w, "upload.html", nil)
		return
	}
	http.Redirect(w, r, "/500", http.StatusSeeOther)
}

// UploadPostAPI returns api call for uploading file
func (s *Server) UploadPostAPI(w http.ResponseWriter, r *http.Request) {
	statusCode := s.uploadPost(w, r)
	w.WriteHeader(statusCode)
}

// UploadGet for getting upload page
func (s *Server) UploadGet(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, "upload.html", nil)
}
