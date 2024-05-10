package grpc

import (
	"context"
	"fmt"

	"github.com/isayme/tox/util"
	"google.golang.org/grpc/metadata"
)

type jwtToken struct {
	key                      []byte
	requireTransportSecurity bool
}

func newJwtToken(key []byte, requireTransportSecurity bool) *jwtToken {
	return &jwtToken{
		key:                      key,
		requireTransportSecurity: requireTransportSecurity,
	}
}

func (t *jwtToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	ss, err := util.GenerateJwtToken(t.key)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"token": ss,
	}, nil
}

func (t *jwtToken) RequireTransportSecurity() bool {
	return t.requireTransportSecurity
}

func VerifyTokenFromContext(ctx context.Context, key []byte) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("get context fail")
	}

	tokens, ok := md["token"]
	if !ok || len(tokens) == 0 {
		return fmt.Errorf("no token")
	}

	tokenString := tokens[0]
	err := util.ValidateJwtToken(tokenString, key)
	return err
}
