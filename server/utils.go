package server

import (
	"context"
	"net/http"
	"time"

	"github.com/zrui98/fserv/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

// Claims struct containing username and jwt claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type contextKey string

var usernameKey contextKey = "usernameKey"

func (s *Server) getValidJWTUsername(ctx context.Context, c *http.Cookie) string {
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return s.Config.JWT_KEY, nil
	})
	if err != nil {
		glog.Errorf("Error with token: %v\n", err)
		return ""
	}
	if !tkn.Valid {
		glog.Error("Token was invalid!")
		return ""
	}
	var userWithName *models.User
	userWithName, err = s.users.GetUserById(ctx, claims.Username)
	if err != nil {
		glog.Errorf("Getting last login time failed: %v\n", err)
		return ""
	}
	if userWithName.LastLoginTime.After(time.Unix(claims.IssuedAt, 0)) {
		glog.Error("User already logged in, token expired")
		return ""
	}
	username := userWithName.Username
	if username == "" {
		glog.Error("Error getting username")
		return ""
	}
	return username
}
