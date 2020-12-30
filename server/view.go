package server

import (
	"io/ioutil"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/zrui98/fserv/constants"
)

type TextData struct {
	TextContents string
	FileName     string
}

func (server *Server) ViewFile(w http.ResponseWriter, r *http.Request) {
	url_id := chi.URLParam(r, "fileId")
	file, err := server.files.GetFileById(r.Context(), url_id)
	if err != nil || file == nil {
		glog.Errorf("Failed to fetch files:: %v\n", err)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}

	glog.Info(file.IsPrivate)
	if file.IsPrivate {
		c, err := r.Cookie("jwtAccessToken")
		if err != nil {
			glog.Errorf("Error with cookie:: %v\n", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
		username := server.GetValidJWTUsername(c, r.Context())
		if username != file.Owner {
			glog.Error("Tried to access private file!")
			http.Redirect(w, r, "/400", http.StatusSeeOther)
			return
		}
	}
	mime, err := mimetype.DetectFile(file.FilePath)
	if err != nil {
		glog.Errorf("Error unrecognized mimetype:: %v\n", err)
		return
	}
	if mimetype.EqualsAny(mime.String(), constants.Image...) || mimetype.EqualsAny(mime.String(), constants.Video...) {
		http.ServeFile(w, r, file.FilePath)
	} else {
		rawcontent, err := ioutil.ReadFile(file.FilePath)
		if err != nil {
			glog.Errorf("Error getting text file contents for file %s", file.FilePath)
		}
		server.renderPage(w, "album.html", &TextData{
			TextContents: string(rawcontent),
			FileName:     file.FileName,
		})
	}
}
