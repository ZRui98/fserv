package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang/glog"
	"github.com/dgrijalva/jwt-go"
)

func (server *Server) LoginPost(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			glog.Errorf("Parsing form failed:: %v\n", err)
			return
		}
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		user, err := server.users.GetUserById(r.Context(), username)
		if err != nil {
			glog.Errorf("Querying DB failed:: %v\n", err)
			return
		}
		match, err := argon2id.ComparePasswordAndHash(password, user.Password)
		if err != nil {
			glog.Errorf("Error validating password:: %v\n", err)
			return
		}

		if !match {
			w.WriteHeader(http.StatusUnauthorized)
			glog.Error("Password did not match")
			return
		}
		loginTimestamp := time.Now()
		err = server.users.UpdateUserLastLoginTime(r.Context(), username, loginTimestamp)
		expirationTime := time.Now().Add(3 * time.Hour)
		claims := &Claims {
			Username: username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
				IssuedAt: loginTimestamp.Unix() + 1,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(JWT_KEY)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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

func (server *Server) LoginGet(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, nil)
}
