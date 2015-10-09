package protoprocess

import "go.pedge.io/proto/stream"

type client struct {
	compressor Compressor
	opts       ClientOptions
}

func newClient(compressor Compressor, opts ClientOptions) *client {
	return &client{compressor, opts}
}

func (c *client) Process(
	dirPath string,
	streamingBytesServer protostream.StreamingBytesServer,
) error {
	return nil
}
