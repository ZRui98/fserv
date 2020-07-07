package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang/glog"
	"github.com/dgrijalva/jwt-go"
)

type LoginErrors struct {
	Username string
	Password string
}

func (s *Server) LoginPost(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			glog.Errorf("Parsing form failed:: %v\n", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		user, err := s.users.GetUserById(r.Context(), username)
		if err != nil {
			glog.Errorf("Querying DB failed:: %v\n", err)
			renderLogin(w, &LoginErrors{
				Username: "Invalid Username",
			})
			return
		}
		match, err := argon2id.ComparePasswordAndHash(password, user.Password)
		if err != nil {
			glog.Errorf("Error validating password:: %v\n", err)
			renderLogin(w, &LoginErrors{
				Password: "Invalid Password",
			})
			return
		}

		if !match {
			glog.Error("Password did not match")
			renderLogin(w, &LoginErrors{
				Password: "Invalid Password",
			})
			return
		}
		loginTimestamp := time.Now().Add(-1*time.Second)
		err = s.users.UpdateUserLastLoginTime(r.Context(), username, loginTimestamp)
		expirationTime := time.Now().Add(3 * time.Hour)
		claims := &Claims {
			Username: username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
				IssuedAt: loginTimestamp.Unix() + 1,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(s.Config.JWT_KEY)
		if err != nil {
			glog.Errorf("Error creating JWT:: %v\n", err)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name: "jwtsecret",
			Value: tokenString,
			Expires: expirationTime,
			SameSite: http.SameSiteNoneMode,
			Secure: true,
			Path: "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) LoginGet(w http.ResponseWriter, r *http.Request) {
	renderLogin(w, &LoginErrors{})
}

func renderLogin(w http.ResponseWriter, v *LoginErrors) {
	t, _ := template.ParseFiles("templates/login.html", "templates/head.tmpl", "templates/navbar.tmpl")
	t.Execute(w, v)
}
