package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang/glog"
	"github.com/dgrijalva/jwt-go"
)

func (s *Server) LoginPost(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			glog.Errorf("Parsing form failed:: %v\n", err)
			return
		}
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		user, err := s.users.GetUserById(r.Context(), username)
		if err != nil {
			glog.Errorf("Querying DB failed:: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		match, err := argon2id.ComparePasswordAndHash(password, user.Password)
		if err != nil {
			glog.Errorf("Error validating password:: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !match {
			w.WriteHeader(http.StatusUnauthorized)
			glog.Error("Password did not match")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		loginTimestamp := time.Now()
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
			w.WriteHeader(http.StatusInternalServerError)
			http.Redirect(w, r, "/500", http.StatusSeeOther)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name: "token",
			Value: tokenString,
			Expires: expirationTime,
			Secure: true,
			SameSite: http.SameSiteNoneMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) LoginGet(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/login.html", "templates/head.tmpl", "templates/navbar.tmpl")
	t.Execute(w, nil)
}
