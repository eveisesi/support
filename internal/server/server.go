package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/embersyndicate/support/internal"

	"github.com/embersyndicate/support/internal/category"
	"github.com/embersyndicate/support/internal/key"
	"github.com/embersyndicate/support/internal/ticket"
	"github.com/embersyndicate/support/internal/token"
	"github.com/embersyndicate/support/internal/user"
	"github.com/embersyndicate/support/pkg/middleware"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

// Server is a representation of this applications HTTP Server
// To it we attach out GraphQL Engine provided by GQL Gen ()
type server struct {
	logger   *logrus.Logger
	redis    *redis.Client
	newrelic *newrelic.Application

	server *http.Server

	category category.Service
	key      key.Service
	ticket   ticket.Service
	token    token.Service
	user     user.Service
}

// New returns an instance of our HTTP Server
func New(port uint, logger *logrus.Logger, redis *redis.Client, newrelic *newrelic.Application, category category.Service, key key.Service, ticket ticket.Service, token token.Service, user user.Service) *server {
	s := &server{
		logger:   logger,
		redis:    redis,
		newrelic: newrelic,

		category: category,
		key:      key,
		ticket:   ticket,
		token:    token,
		user:     user,
	}

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      s.router(),
	}

	return s

}

func (s *server) Run() error {
	s.logger.WithField("address", s.server.Addr).Info("starting http server")
	return s.server.ListenAndServe()
}

func (s *server) router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Group(func(r chi.Router) {
		r.Use(
			s.monitoring,
			middleware.RequestID,
			middleware.ContentTypeJSON,
			middleware.CORS,
			middleware.RequestLogger(s.logger),
		)

		r.Get("/.well-known/jwks.json", s.handleV1GetJWKS)

		r.Route("/v1", func(r chi.Router) {
			r.Post("/users/register", s.handleV1PostUserRegister)
			r.Post("/users/login", s.handleV1PostUserLogin)

			r.Group(func(r chi.Router) {
				r.Use(s.auth)
				r.Get("/categories", s.handleV1GetCategories)
				r.Post("/categories", s.handleV1PostCategories)
				r.Get("/categories/{categoryID}", s.handleV1GetCategory)
				r.Patch("/categories/{categoryID}", s.handleV1PatchCategory)

				r.Post("/tickets", s.handleV1PostTickets)

				r.Get("/tickets/statuses", s.handleV1GetTicketStatuses)
				r.Post("/tickets/statuses", s.handleV1PostTicketStatuses)
				r.Get("/tickets/statuses/{statusID}", s.handleV1GetTicketStatus)
				r.Patch("/tickets/statuses/{statusID}", s.handleV1PatchTicketStatus)

				r.Get("/tickets/definitions", s.handleV1GetTicketDefinitions)
				r.Post("/tickets/definitions", s.handleV1PostTicketDefinition)
				r.Get("/tickets/definitions/{definitionID}", s.handleV1GetTicketDefinition)
				r.Patch("/tickets/definitions/{definitionID}", s.handleV1PatchTicketDefinition)

				r.Get("/fields/definitions", s.handleV1GetFieldDefinitions)
				r.Post("/fields/definitions", s.handleV1PostFieldDefinitions)
				r.Get("/fields/definitions/{definitionID}", s.handleV1GetFieldDefinition)
				r.Patch("/fields/definitions/{definitionID}", s.handleV1PatchFieldDefinition)

			})
		})

	})

	return r
}

// GracefullyShutdown gracefully shuts down the HTTP API.
func (s *server) GracefullyShutdown(ctx context.Context) error {
	s.logger.Info("attempting to shutdown server gracefully")
	return s.server.Shutdown(ctx)
}

func (s *server) writeResponse(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {

	w.WriteHeader(code)

	if data != nil {
		switch d := data.(type) {
		case []byte:
			_, _ = w.Write(d)
		default:
			_ = json.NewEncoder(w).Encode(d)
		}
	}
}

func (s *server) writeError(ctx context.Context, w http.ResponseWriter, code int, err error, isNr bool) {

	if err != nil {
		var ierr internal.InternalError
		if errors.As(err, &ierr) {
			switch ierr.Level {
			case internal.LevelInternal:
				code = http.StatusInternalServerError
			case internal.LevelBad:
				code = http.StatusBadRequest
			}
		}
		msg := err.Error()
		reqID := middleware.GetRequestID(ctx)
		if reqID != "" {
			msg = fmt.Sprintf("%s (RequestID: %s)", msg, reqID)
		}
		s.writeResponse(ctx, w, code, map[string]interface{}{
			"message": msg,
		})
		return
	}

	s.writeResponse(ctx, w, code, nil)

}
