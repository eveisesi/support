package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
}

// New returns an instance of our HTTP Server
func New(port uint, logger *logrus.Logger, redis *redis.Client, newrelic *newrelic.Application) *server {
	s := &server{
		logger:   logger,
		redis:    redis,
		newrelic: newrelic,
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

	return r
}

// GracefullyShutdown gracefully shuts down the HTTP API.
func (s *server) GracefullyShutdown(ctx context.Context) error {
	s.logger.Info("attempting to shutdown server gracefully")
	return s.server.Shutdown(ctx)
}

func (s *server) writeResponse(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {

	if code != http.StatusOK {
		w.WriteHeader(code)
	}

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (s *server) writeError(ctx context.Context, w http.ResponseWriter, code int, err error) {

	// If err is not nil, actually pass in a map so that the output to the wire is {"error": "text...."} else just let it fall through
	if err != nil {
		newrelic.FromContext(ctx).NoticeError(err)
		s.writeResponse(ctx, w, code, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	s.writeResponse(ctx, w, code, nil)

}
