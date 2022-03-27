package auth

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bridger217/go-grpc-api/pkg/middleware"
)

// Adds state (firebase app) to authentication
type Authenticator struct {
	fb *firebase.App
}

func NewAuthenticator(fb *firebase.App) *Authenticator {
	return &Authenticator{fb: fb}
}

// Used by middleware to authenticate requests
func (a *Authenticator) Authenticate(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	authedToken, err := authenticateToken(ctx, token, a.fb)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	// // Not sure if this would be needed:
	// // grpc_ctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))

	newCtx := middleware.AddAuthTokenToContext(ctx, authedToken)

	return newCtx, nil
}

// Authenticates the token with firebase, then returns the uid
// (which should become the user's id in our storage).
func authenticateToken(ctx context.Context, token string, fb *firebase.App) (*auth.Token, error) {
	client, err := fb.Auth(ctx)
	if err != nil {
		return nil, err
	}

	authedToken, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return authedToken, nil
}
