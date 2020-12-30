package middleware

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type contextKey int64

const (
	contextKeyRequestID contextKey = iota
	contextKeyUserID
	contextKeyToken
)

func RequestID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := uuid.New()

		ctx := r.Context()

		newrelic.FromContext(ctx).AddAttribute("requestID", id.String())

		ctx = context.WithValue(ctx, contextKeyRequestID, id.String())

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func GetRequestID(ctx context.Context) string {

	req := ctx.Value(contextKeyRequestID)

	if id, ok := req.(string); ok {
		return id
	}

	return ""

}

func SetTokenOnContext(ctx context.Context, token jwt.Token) context.Context {
	return context.WithValue(ctx, contextKeyToken, token)
}

func GetTokenFromContext(ctx context.Context) jwt.Token {

	req := ctx.Value(contextKeyToken)

	if token, ok := req.(jwt.Token); ok {
		return token
	}

	return nil
}

func SetUserIDOnContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKeyUserID, id)
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {

	req := ctx.Value(contextKeyUserID)

	if id, ok := req.(string); ok {
		return id, true
	}

	return "", false

}

func GetUserObjectIDFromContext(ctx context.Context) (primitive.ObjectID, error) {

	if userID, valid := GetUserIDFromContext(ctx); valid {
		return primitive.ObjectIDFromHex(userID)
	}

	return primitive.NilObjectID, fmt.Errorf("invalid id returns from context")
}
