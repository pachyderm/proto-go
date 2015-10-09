package protoprocess

import "go.pedge.io/proto/stream"

type server struct {
	compressor Compressor
	opts       ServerOptions
}

func newServer(compressor Compressor, opts ServerOptions) *server {
	return &server{compressor, opts}
}

func (s *server) Handle(
	processor Processor,
	streamingBytesClient protostream.StreamingBytesClient,
	streamingBytesServer protostream.StreamingBytesServer,
) error {
	return nil
}
