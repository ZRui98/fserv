package server

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gabriel-vasile/mimetype"
	"github.com/golang/glog"
	"github.com/zrui98/fserv/constants"
	"github.com/zrui98/fserv/models"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	length   = uint32(len(alphabet))
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func encode(n uint32) string {
	if n == 0 {
		return "0"
	}

	b := make([]byte, 0, 512)
	for n > 0 {
		r := n % length
		n /= length
		b = append([]byte{alphabet[r]}, b...)
	}
	return string(b)
}

func (s *Server) UploadPost(w http.ResponseWriter, r *http.Request) {
	glog.Info("Uploading File")
	r.ParseMultipartForm(32 << 20)
	isPrivate := r.FormValue("private") == "true"
	uploadedFile, handler, err := r.FormFile("uploadFile")
	if err != nil {
		glog.Errorf("Error parsing form file:: %v\n", err)
		return
	}
	defer uploadedFile.Close()
	var seed []byte
	seed, err = generateRandomBytes(4)
	if err != nil {
		seed, err = generateRandomBytes(4)
		glog.Errorf("Error generating seed:: %v\n", err)
	}

	mime, err := mimetype.DetectReader(uploadedFile)
	if err != nil {
		glog.Errorf("Error unrecognized mimetype:: %v\n", err)
		return
	}
	folder := "txt"
	if mimetype.EqualsAny(mime.String(), constants.Image...) {
		folder = "images"
	} else if mimetype.EqualsAny(mime.String(), constants.Video...) {
		folder = "videos"
	}
	token := encode(binary.BigEndian.Uint32(seed))
	file_path := path.Join(
		s.Config.ROOT_DIR,
		r.Context().Value(UsernameKey).(string),
		folder,
		token+path.Ext(handler.Filename),
	)
	glog.Infof("Saving to: %s File: %s\n", file_path, handler.Filename)
	file := &models.File{
		UrlId:     token,
		FilePath:  file_path,
		FileName:  handler.Filename,
		Owner:     r.Context().Value(UsernameKey).(string),
		IsPrivate: isPrivate,
	}
	err = s.files.AddFile(r.Context(), file)
	if err != nil {
		glog.Errorf("Error inserting data into DB: %v\n", err)
		return
	}

	f, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		glog.Errorf("Error creating file pointer: %v\n", err)
		return
	}
	defer f.Close()
	if _, err := uploadedFile.Seek(0, 0); err != nil {
		glog.Errorf("Error seeking uploaded file: %v\n", err)
		return
	}
	io.Copy(f, uploadedFile)
	renderPage(w, "templates/upload.html", nil)
}

func (s *Server) UploadGet(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "templates/upload.html", nil)
}
