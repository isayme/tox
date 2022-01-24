package grpc

import (
	"context"
	"fmt"

	"github.com/isayme/tox/util"
	"google.golang.org/grpc/metadata"
)

type JwtToken struct {
	key []byte
}

func newJwtToken(key []byte) *JwtToken {
	return &JwtToken{
		key: key,
	}
}

func (t *JwtToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	ss, err := util.GenerateJwtToken(t.key)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"token": ss,
	}, nil
}

func (t *JwtToken) RequireTransportSecurity() bool {
	return true
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
