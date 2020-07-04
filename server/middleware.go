package server

import (
	"net/http"
	"context"
	"os"
	"time"

	"github.com/zrui98/fserv/models"

	"github.com/golang/glog"
	"github.com/dgrijalva/jwt-go"
)

var JWT_KEY = []byte(os.Getenv("SECRET_KEY"))
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ContextKey string
var UsernameKey ContextKey = "usernameKey"

func (s *Server) ValidateJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		c, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", 301)
			return
		}
		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			http.Redirect(w, r, "/login", 301)
			return JWT_KEY, nil
		})
		if err != nil {
			http.Redirect(w, r, "/login", 301)
			return
		}
		if !tkn.Valid {
			http.Redirect(w, r, "/login", 301)
			return
		}
		var userWithName *models.User
		userWithName, err = s.users.GetUserById(r.Context(), claims.Username)
		if err != nil {
			glog.Errorf("Getting last login time failed: %v\n", err)
			http.Redirect(w, r, "/login", 301)
			return
		}
		if time.Unix(claims.IssuedAt, 0).Before(userWithName.LastLoginTime) {
			glog.Error("User already logged in, token expired")
			http.Redirect(w, r, "/login", 301)
			return
		}
		username := userWithName.Username
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), UsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
