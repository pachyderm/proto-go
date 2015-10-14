package protoprocess // import "go.pedge.io/proto/process"

import (
	"go.pedge.io/pkg/archive"
	"go.pedge.io/proto/stream"
)

var (
	// DefaultChunkSizeBytes is the default chunk size used for sending and receiving.
	DefaultChunkSizeBytes = 65536
)

// Processor processes the contents of a directory.
//
// The result of a call to a Processor is that the directory
// may have modified contents.
type Processor interface {
	Process(dirPath string) error
}

// NewAPIProcessor returns a new Processor that calls an APIClient.
func NewAPIProcessor(client Client, apiClient APIClient) Processor {
	return newAPIProcessor(client, apiClient)
}

// Client calls a remote Processor.
type Client interface {
	Process(dirPath string, streamingBytesDuplexCloser protostream.StreamingBytesDuplexCloser) error
}

// ClientOptions are options to the construction of a Client.
type ClientOptions struct {
	ChunkSizeBytes int
}

// NewClient returns a new Client.
func NewClient(archiver pkgarchive.Archiver, opts ClientOptions) Client {
	return newClient(archiver, opts)
}

// Server wraps a Processor.
type Server interface {
	Handle(processor Processor, streamingBytesDuplexer protostream.StreamingBytesDuplexer) error
}

// ServerOptions are options to the construction of a Server.
type ServerOptions struct {
	ChunkSizeBytes int
}

// NewServer returns a new Server.
func NewServer(archiver pkgarchive.Archiver, opts ServerOptions) Server {
	return newServer(archiver, opts)
}

// NewAPIServer returns a new APIServer for the given Server.
func NewAPIServer(processor Processor, server Server) APIServer {
	return newAPIServer(processor, server)
}
