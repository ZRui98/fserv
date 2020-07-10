package server

import (
	"os"
	"net/http"
	"time"
	"path"
	
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
	if len(username) < 3 {
		glog.Errorf("Username Too Short!")
		renderPage(w, "templates/register.html", &RegistrationErrors{
			Username: "Username too short",
		})
		return
	}
	userExists, _ := s.users.GetUserById(r.Context(), username)
	if len(userExists.Username) > 0 {
		glog.Errorf("Failed to claim unique username:: %v\n", userExists.Username)
		renderPage(w, "templates/register.html", &RegistrationErrors{
			Username: "Invalid Username (taken or bad form)",
		})
		return
	}
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
	if err != nil {
		glog.Errorf("Querying DB failed:: %v\n", err)
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	userDir := path.Join(s.Config.ROOT_DIR, username)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		err = os.Mkdir(userDir, os.ModePerm)
		if err != nil {
			glog.Errorf("Failed to make user dir:: %v\n", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
		for _, p := range constants.UserFolders {
			err = os.Mkdir(userDir + "/" + p, os.ModePerm)
			if err != nil {
				glog.Errorf("Failed to make directory:: %v\n", err)
				http.Redirect(w, r, "/500", http.StatusSeeOther)
				return
			}
		}
	}
	err = s.users.AddUserById(r.Context(), user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) RegisterGet(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "templates/register.html", &RegistrationErrors{})
}
