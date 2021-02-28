package server

import (
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

// LoginErrors struct for login errors for username or password
type LoginErrors struct {
	Username string
	Password string
}

func (s *Server) loginPost(w http.ResponseWriter, r *http.Request) (int, *LoginErrors) {
	err := r.ParseForm()
	if err != nil {
		glog.Errorf("Parsing form failed:: %v\n", err)
		return http.StatusInternalServerError, nil
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	user, err := s.users.GetUserById(r.Context(), username)
	if err != nil {
		glog.Errorf("Querying DB failed:: %v\n", err)
		return http.StatusBadRequest, &LoginErrors{
			Password: "Failed to Login! Refresh the page?",
		}
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.Password)
	if err != nil {
		glog.Errorf("Error validating password:: %v\n", err)
		return http.StatusForbidden, &LoginErrors{
			Password: "Invalid Password",
		}
	}
	if !match {
		glog.Error("Password did not match")
		return http.StatusForbidden, &LoginErrors{
			Password: "Invalid Password",
		}
	}

	loginTimestamp := time.Now().Add(-1 * time.Second)
	err = s.users.UpdateUserLastLoginTime(r.Context(), username, loginTimestamp)
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  loginTimestamp.Unix() + 1,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.Config.JWT_KEY)
	if err != nil {
		glog.Errorf("Error creating JWT:: %v\n", err)
		return http.StatusForbidden, nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwtAccessToken",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   s.Config.SECURE_COOKIE,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	return http.StatusOK, nil
}

func (s *Server) LoginPost(w http.ResponseWriter, r *http.Request) {
	statusCode, loginErrors := s.loginPost(w, r)
	if loginErrors != nil {
		s.renderPage(w, "login.html", loginErrors)
		return
	}
	if statusCode == http.StatusInternalServerError {
		http.Redirect(w, r, "/500", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) LoginPostAPI(w http.ResponseWriter, r *http.Request) {
	statusCode, _ := s.loginPost(w, r)
	w.WriteHeader(statusCode)
}

func (s *Server) LoginGet(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, "login.html", LoginErrors{})
}
