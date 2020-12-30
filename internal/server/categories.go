package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/embersyndicate/support"
	"github.com/go-chi/chi"
)

func (s *server) handleV1GetCategories(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	categories, err := s.category.Categories(ctx)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, categories)

}

func (s *server) handleV1PostCategories(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var category = new(support.Category)
	err := json.NewDecoder(r.Body).Decode(category)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, err, false)
		return
	}

	_, err = s.category.CreateCategory(ctx, category)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusCreated, nil)

}

func (s *server) handleV1GetCategory(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := chi.URLParam(r, "categoryID")
	if id == "" {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("categoryID is required, empty value received"), false)
		return
	}

	category, err := s.category.Category(ctx, id)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, category)

}

func (s *server) handleV1PatchCategory(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := chi.URLParam(r, "categoryID")
	if id == "" {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("categoryID is required, empty value received"), false)
		return
	}

	category, err := s.category.Category(ctx, id)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, err, false)
		return
	}

	err = json.NewDecoder(r.Body).Decode(category)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to decode response body: %w", err), false)
		return
	}

	_, err = s.category.UpdateCategory(ctx, id, category)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError,
			fmt.Errorf("failed to update category"),
			false,
		)
		return
	}

	s.writeResponse(ctx, w, http.StatusNoContent, nil)

}
