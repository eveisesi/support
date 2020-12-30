package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/embersyndicate/support"
)

func (s *server) handleV1PostUserLogin(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var user = new(support.User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err), false)
		return
	}

	key, err := s.user.Login(ctx, user)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, key)

}

func (s *server) handleV1PostUserRegister(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var user = new(support.User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err), false)
		return
	}

	user, err = s.user.Register(ctx, user)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, user)

}
