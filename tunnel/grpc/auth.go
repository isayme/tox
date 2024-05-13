package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

type jwtToken struct {
	token                    string
	requireTransportSecurity bool
}

func newJwtToken(token string, requireTransportSecurity bool) *jwtToken {
	return &jwtToken{
		token:                    token,
		requireTransportSecurity: requireTransportSecurity,
	}
}

func (t *jwtToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": t.token,
	}, nil
}

func (t *jwtToken) RequireTransportSecurity() bool {
	return t.requireTransportSecurity
}

func VerifyTokenFromContext(ctx context.Context, token string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("get context fail")
	}

	tokens, ok := md["token"]
	if !ok || len(tokens) == 0 {
		return fmt.Errorf("no token")
	}

	tokenString := tokens[0]
	if tokenString != token {
		return fmt.Errorf("token invalid")
	}

	return nil
}
