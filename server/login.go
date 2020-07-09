package server

import (
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
			renderPage(w, "templates/login.html", LoginErrors{
				Password: "Invalid Password",
			})
			return
		}

		match, err := argon2id.ComparePasswordAndHash(password, user.Password)
		if err != nil {
			glog.Errorf("Error validating password:: %v\n", err)
			renderPage(w, "templates/login.html", LoginErrors{
				Password: "Invalid Password",
			})
			return
		}
		if !match {
			glog.Error("Password did not match")
			renderPage(w, "templates/login.html", LoginErrors{
				Password: "Invalid Password",
			})
			return
		}

		loginTimestamp := time.Now().Add(-1*time.Second)
		err = s.users.UpdateUserLastLoginTime(r.Context(), username, loginTimestamp)
		expirationTime := time.Now().Add(30 * time.Minute)
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
			Name: "jwtAccessToken",
			Value: tokenString,
			Expires: expirationTime,
			HttpOnly: true,
			Secure: true,
			SameSite: http.SameSiteNoneMode,
			Path: "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) LoginGet(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "templates/login.html", LoginErrors{})
}
