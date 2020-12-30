package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/embersyndicate/support"
	"github.com/go-chi/chi"
)

func (s *server) handleV1PostTickets(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var ticket = new(support.Ticket)
	err := json.NewDecoder(r.Body).Decode(ticket)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err), false)
		return
	}

	_, err = s.ticket.CreateTicket(ctx, ticket)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusCreated, nil)
}

func (s *server) handleV1GetTicketStatuses(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	statuses, err := s.ticket.TicketStatuses(ctx)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, statuses)

}

func (s *server) handleV1PostTicketStatuses(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var status = new(support.TicketStatus)
	err := json.NewDecoder(r.Body).Decode(status)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err), false)
		return
	}

	_, err = s.ticket.CreateTicketStatus(ctx, status)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusCreated, nil)

}

func (s *server) handleV1GetTicketStatus(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := chi.URLParam(r, "statusID")
	if id == "" {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("statusID is required, empty value received"), false)
		return
	}

	status, err := s.ticket.TicketStatus(ctx, id)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, status)

}

func (s *server) handleV1PatchTicketStatus(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := chi.URLParam(r, "statusID")
	if id == "" {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("statusID is required, empty value received"), false)
		return
	}

	status, err := s.ticket.TicketStatus(ctx, id)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	err = json.NewDecoder(r.Body).Decode(status)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err), false)
		return
	}

	_, err = s.ticket.UpdateTicketStatus(ctx, id, status)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusNoContent, nil)

}

func (s *server) handleV1GetTicketDefinitions(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	definitions, err := s.ticket.TicketDefinitions(ctx)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, definitions)

}

func (s *server) handleV1PostTicketDefinition(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var definition = new(support.TicketDefinition)
	err := json.NewDecoder(r.Body).Decode(definition)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err), false)
		return
	}

	definition, err = s.ticket.CreateTicketDefinition(ctx, definition)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusCreated, definition)

}

func (s *server) handleV1GetTicketDefinition(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := chi.URLParam(r, "definitionID")
	if id == "" {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("definitionID is required, empty value received"), false)
		return
	}

	definition, err := s.ticket.TicketDefinition(ctx, id)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, definition)

}

func (s *server) handleV1PatchTicketDefinition(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := chi.URLParam(r, "definitionID")
	if id == "" {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("definitionID is required, empty value received"), false)
		return
	}

	definition, err := s.ticket.TicketDefinition(ctx, id)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	err = json.NewDecoder(r.Body).Decode(definition)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err), false)
		return
	}

	definition, err = s.ticket.UpdateTicketDefinition(ctx, id, definition)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusNoContent, definition)

}

func (s *server) handleV1GetFieldDefinitions(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	fields, err := s.ticket.FieldDefinitions(ctx)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, fields)

}

func (s *server) handleV1PostFieldDefinitions(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var definition = new(support.FieldDefinition)
	err := json.NewDecoder(r.Body).Decode(definition)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err), false)
		return
	}

	definition, err = s.ticket.CreateFieldDefinition(ctx, definition)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusCreated, definition)

}

func (s *server) handleV1GetFieldDefinition(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := chi.URLParam(r, "definitionID")
	if id == "" {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("definitionID is required, empty value received"), false)
		return
	}

	definitionID, err := s.ticket.FieldDefinition(ctx, id)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, definitionID)

}

func (s *server) handleV1PatchFieldDefinition(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := chi.URLParam(r, "definitionID")
	if id == "" {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("definitionID is required, empty value received"), false)
		return
	}

	definition, err := s.ticket.FieldDefinition(ctx, id)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	err = json.NewDecoder(r.Body).Decode(definition)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to read request body: %w", err), false)
		return
	}

	_, err = s.ticket.UpdateFieldDefinition(ctx, id, definition)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err, false)
		return
	}

	s.writeResponse(ctx, w, http.StatusNoContent, nil)

}
