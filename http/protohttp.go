package protohttp // import "go.pedge.io/proto/http"

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

var (
	// ErrInvalidAuthorization says that authorization was present on a context.Context, but was invalid.
	ErrInvalidAuthorization = errors.New("pkghttp: invalid authorization on context")

	jsonMarshaler = &jsonpb.Marshaler{}
)

// Do does the given method and url over http, using the given context and request.
//
// The http request body is the serialized request.
// The response will be the deserialized http response body.
// If basic auth is present, it will be added.
//
// Returns the http response body.
func Do(
	ctx context.Context,
	httpClient *http.Client,
	method string,
	url string,
	request proto.Message,
	response proto.Message,
) (retVal []byte, retErr error) {
	data, err := jsonMarshaler.MarshalToString(request)
	if err != nil {
		return nil, err
	}
	httpRequest, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	basicAuth, err := BasicAuthFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if basicAuth != nil {
		httpRequest.SetBasicAuth(basicAuth.Username, basicAuth.Password)
	}
	httpResponse, err := ctxhttp.Do(ctx, httpClient, httpRequest)
	if err != nil {
		return nil, err
	}
	defer func() {
		if httpResponse.Body != nil {
			if err := httpResponse.Body.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}
	}()
	var body []byte
	if httpResponse.Body != nil {
		body, err = ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return body, err
		}
	}
	if httpResponse.StatusCode != http.StatusOK {
		return body, errors.New(httpResponse.Status)
	}
	return body, jsonpb.UnmarshalString(string(body), response)
}

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
