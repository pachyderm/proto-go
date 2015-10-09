package protoprocess // import "go.pedge.io/proto/process"

import (
	"io"

	"go.pedge.io/proto/stream"
)

var (
	// DefaultExcludePatternsFiles is the slice of file names
	// that will contain patterns of files to exclude on upload.
	DefaultExcludePatternsFiles = []string{
		".protoprocessignore",
	}
)

// Processor processes the contents of a directory.
//
// The result of a call to a Processor is that the directory
// may have modified contents.
type Processor interface {
	Process(dirPath string) error
}

// Compressor compresses and decompresses a directory.
type Compressor interface {
	Compress(dirPath string) (io.ReadCloser, error)
	Decompress(reader io.Reader, dirPath string) error
}

// CompressorOptions are options to the construction of a Compressor.
type CompressorOptions struct {
	// If not set, DefaultExcludePatternsFiles is used.
	ExcludePatternsFiles []string
}

// NewTarCompressor returns a new Compressor for tar.
func NewTarCompressor(opts CompressorOptions) Compressor {
	return newTarCompressor(opts)
}

// Client calls a remote Processor.
type Client interface {
	Process(
		dirPath string,
		streamingBytesServer protostream.StreamingBytesServer,
	) error
}

// ClientOptions are options to the construction of a Client.
type ClientOptions struct {
}

// NewClient returns a new Client.
func NewClient(compressor Compressor, opts ClientOptions) Client {
	return newClient(compressor, opts)
}

// Server wraps a Processor.
type Server interface {
	Handle(
		processor Processor,
		streamingBytesClient protostream.StreamingBytesClient,
		streamingBytesServer protostream.StreamingBytesServer,
	) error
}

// ServerOptions are options to the construction of a Server.
type ServerOptions struct {
}

// NewServer returns a new Server.
func NewServer(compressor Compressor, opts ServerOptions) Server {
	return newServer(compressor, opts)
}
