package server

import (
	"context"
	"net/http"

	"github.com/golang/glog"
)

func (s *Server) validateJwt(r *http.Request) (int, context.Context) {
	c, err := r.Cookie("jwtAccessToken")
	if err != nil {
		glog.Errorf("Error with cookie: %v\n", err)
		return http.StatusForbidden, nil
	}
	username := s.getValidJWTUsername(r.Context(), c)
	if username == "" {
		return http.StatusForbidden, nil
	}
	ctx := context.WithValue(r.Context(), usernameKey, username)
	return http.StatusOK, ctx
}

// ValidateJwt middleware for validating JWT in website
func (s *Server) ValidateJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode, ctx := s.validateJwt(r)
		if statusCode != http.StatusOK {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ValidateJwtAPI middlware for validating JWT in API
func (s *Server) ValidateJwtAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode, ctx := s.validateJwt(r)
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
