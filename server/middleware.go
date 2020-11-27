package server

import (
	"context"
	"net/http"
	"time"

	"github.com/zrui98/fserv/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ContextKey string

var UsernameKey ContextKey = "usernameKey"

func (s *Server) ValidateJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwtAccessToken")
		if err != nil {
			glog.Errorf("Error with cookie: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return s.Config.JWT_KEY, nil
		})
		if err != nil {
			glog.Errorf("Error with token: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if !tkn.Valid {
			glog.Error("Token was invalid!")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		var userWithName *models.User
		userWithName, err = s.users.GetUserById(r.Context(), claims.Username)
		if err != nil {
			glog.Errorf("Getting last login time failed: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if userWithName.LastLoginTime.After(time.Unix(claims.IssuedAt, 0)) {
			glog.Error("User already logged in, token expired")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		username := userWithName.Username
		if username == "" {
			glog.Error("Error getting username")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), UsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
