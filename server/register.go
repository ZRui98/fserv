package server

import (
	"os"
	"net/http"
	"time"
	
	"github.com/zrui98/fserv/models"
	"github.com/zrui98/fserv/constants"
	"github.com/alexedwards/argon2id"
	"github.com/golang/glog"
)

type RegistrationErrors struct {
	Username string
	Password string
	RegistrationKey string
}

func (s *Server) RegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		glog.Errorf("Parsing form failed:: %v\n", err)
		return
	}
	glog.Info("Creating User")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	registration_key := r.PostFormValue("key")
	if (registration_key != s.Config.REGISTRATION_KEY) {
		glog.Error("Wrong Registration Key")
		renderPage(w, "templates/register.html", &RegistrationErrors{
			RegistrationKey: "Invalid Registration Key",
		})
		return
	}
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		glog.Errorf("Failed to verify password:: %v\n", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	user := &models.User{
		Username: username,
		Password: hashedPassword,
		LastLoginTime: time.Now(),
	}
	err = s.users.AddUserById(r.Context(), user)
	if err != nil {
		glog.Errorf("Querying DB failed:: %v\n", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	userDir := "files/" + username
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		os.Mkdir(userDir, os.ModePerm)
		for _, s := range constants.UserFolders {
			os.Mkdir(userDir + "/" + s, os.ModePerm)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) RegisterGet(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "templates/register.html", &RegistrationErrors{})
}
