package protohttp // import "go.pedge.io/proto/http"

// do not actually use this package, this was really just for testing

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"
)

var (
	// ErrInvalidAuthorization says that authorization was present on a context.Context, but was invalid.
	ErrInvalidAuthorization = errors.New("pkghttp: invalid authorization on context")
)

// GetRequestMetadata gets the request metadata for gRPC.
func (c *BasicAuth) GetRequestMetadata(ctx context.Context, uris ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": c.GetAuthorization(),
	}, nil
}

// RequireTransportSecurity says whether BasicAuth requires transport security.
func (c *BasicAuth) RequireTransportSecurity() bool {
	// yaaaaaaaaaa
	return false
}

// GetAuthorization gets the request authorization.
func (c *BasicAuth) GetAuthorization() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.Username, c.Password))))
}

// NewContext returns a new context.Context with the basic auth attached.
func (c *BasicAuth) NewContext(ctx context.Context) context.Context {
	return metadata.NewContext(ctx, metadata.Pairs("Authorization", c.GetAuthorization()))
}

// BasicAuthFromContext gets the basic auth from the specified context.Context.
//
// If no basic auth is present, BasicAuthFromContext returns nil.
func BasicAuthFromContext(ctx context.Context) (*BasicAuth, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, nil
	}
	authorization, ok := md["Authorization"]
	if !ok {
		authorization, ok = md["authorization"]
		if !ok {
			return nil, nil
		}
	}
	if len(authorization) != 1 {
		return nil, ErrInvalidAuthorization
	}
	if !strings.HasPrefix(authorization[0], "Basic ") {
		return nil, ErrInvalidAuthorization
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authorization[0], "Basic "))
	if err != nil {
		return nil, err
	}
	split := strings.SplitN(string(decoded), ":", 2)
	if len(split) != 2 {
		return nil, ErrInvalidAuthorization
	}
	return &BasicAuth{
		Username: split[0],
		Password: split[1],
	}, nil
}
