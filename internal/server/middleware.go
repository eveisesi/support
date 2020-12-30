package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (s *server) monitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		txn := s.newrelic.StartTransaction(r.URL.Path)
		txn.SetWebRequestHTTP(r)
		rw := txn.SetWebResponse(w)
		defer txn.End()

		r = newrelic.RequestWithTransactionContext(r, txn)

		next.ServeHTTP(rw, r)

		rctx := chi.RouteContext(r.Context())
		name := rctx.RoutePattern()

		// ignore invalid routes
		if name == "/*" {
			txn.Ignore()
			return
		}

		txn.SetName(r.Method + " " + name)
	})
}
