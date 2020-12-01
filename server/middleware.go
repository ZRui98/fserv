package server

import (
	"context"
	"net/http"

	"github.com/golang/glog"
)

func (s *Server) ValidateJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwtAccessToken")
		if err != nil {
			glog.Errorf("Error with cookie: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		username := s.GetValidJWTUsername(c, r.Context())
		if username == "" {
			http.Redirect(w, r, "login", http.StatusSeeOther)
		}
		ctx := context.WithValue(r.Context(), UsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
