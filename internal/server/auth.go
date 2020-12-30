package server

import (
	"net/http"
	"strings"

	"github.com/embersyndicate/support/pkg/middleware"
)

func (s *server) auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parsed, err := s.token.ParseAndVerifyToken(ctx, authHeader[7:])
		if err != nil {
			s.writeError(ctx, w, http.StatusUnauthorized, err, false)
			return
		}

		id, err := s.token.GetUserIDFromToken(parsed)
		if err != nil {
			s.writeError(ctx, w, http.StatusUnauthorized, err, false)
			return
		}

		ctx = middleware.SetUserIDOnContext(ctx, id)
		ctx = middleware.SetTokenOnContext(ctx, parsed)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

}
