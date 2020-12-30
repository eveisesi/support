package server

import (
	"net/http"
)

func (s *server) handleV1GetJWKS(w http.ResponseWriter, r *http.Request) {
	s.writeResponse(r.Context(), w, http.StatusOK, s.key.GetPublicJWKS())
}
