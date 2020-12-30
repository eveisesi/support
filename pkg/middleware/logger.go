package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

const (
	logEntryKeyRequestID = "requestID"
	logEntryKeyMethod    = "method"
	logEntryURI          = "uri"
	logEntryRemoteAddr   = "remote_addr"
)

func RequestLogger(logger *logrus.Logger) func(http.Handler) http.Handler {
	return middleware.RequestLogger(&requestLogger{Entry: logrus.NewEntry(logger)})
}

type requestLogger struct {
	Entry *logrus.Entry
}

func (l *requestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {

	fields := map[string]interface{}{
		logEntryKeyMethod:  r.Method,
		logEntryURI:        r.RequestURI,
		logEntryRemoteAddr: r.RemoteAddr,
	}

	if reqID := GetRequestID(r.Context()); reqID != "" {
		fields[logEntryKeyRequestID] = reqID
	}

	return &requestLogEntry{Logger: l.Entry.WithContext(r.Context()).WithFields(fields)}

}

// StructuredLoggerEntry holds our FieldLogger entry
type requestLogEntry struct {
	Logger *logrus.Entry
}

// Write will write to logger entry once the http.Request is complete
func (l *requestLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	e := l.Logger.WithFields(logrus.Fields{
		"status":     status,
		"elapsed_ms": elapsed,
	})

	if _, ok := l.Logger.Data["error"]; ok {
		e.Errorln()
		return
	}
	if status >= http.StatusBadRequest {
		e.Errorln()
		return
	}

	e.Infoln()

}

func (l *requestLogEntry) Panic(v interface{}, stack []byte) {}

// LogEntrySetField will set a new field on a log entry
func LogEntrySetField(ctx context.Context, field string, value interface{}) {
	newrelic.FromContext(ctx).AddAttribute(field, value)
	if entry, ok := ctx.Value(middleware.LogEntryCtxKey).(*requestLogEntry); ok {
		entry.Logger = entry.Logger.WithField(field, value)
	}
}

func LogEntrySetError(ctx context.Context, err error) {
	newrelic.FromContext(ctx).NoticeError(err)
	if entry, ok := ctx.Value(middleware.LogEntryCtxKey).(*requestLogEntry); ok {
		entry.Logger = entry.Logger.WithError(err)
	}
}

// LogEntrySetFields will set a map of key/value pairs on a log entry
func LogEntrySetFields(ctx context.Context, fields map[string]interface{}) {
	for field, value := range fields {
		newrelic.FromContext(ctx).AddAttribute(field, value)
	}
	if entry, ok := ctx.Value(middleware.LogEntryCtxKey).(*requestLogEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
