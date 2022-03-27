package middleware

import (
	"context"

	"firebase.google.com/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Custom key type to avoid ctx collisions
type contextKey string

const (
	contextKeyAuthToken = contextKey("auth-token")
)

// Place the firebase auth token in context
func AddAuthTokenToContext(ctx context.Context, tk *auth.Token) context.Context {
	return context.WithValue(ctx, contextKeyAuthToken, tk)
}

// Parse out the user id
func UserIdFromContext(ctx context.Context) (string, error) {
	t, ok := ctx.Value(contextKeyAuthToken).(*auth.Token)
	if !ok {
		return "", status.Errorf(codes.Internal, "missing auth token in context")
	}
	return t.UID, nil
}
