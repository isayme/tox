package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

type jwtToken struct {
	token              string
	insecureSkipVerify bool
}

func newJwtToken(token string, insecureSkipVerify bool) *jwtToken {
	return &jwtToken{
		token:              token,
		insecureSkipVerify: insecureSkipVerify,
	}
}

func (t *jwtToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": t.token,
	}, nil
}

func (t *jwtToken) RequireTransportSecurity() bool {
	return !t.insecureSkipVerify
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
