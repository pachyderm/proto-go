package protoplugin // import "go.pedge.io/proto/plugin"

import (
	"io"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Generator generates a file from a *descriptor.FileDescriptorProto.
type Generator interface {
	// if the returned io.Reader is nil, this means to not generate any file
	Generate(fileDescriptorProto *descriptor.FileDescriptorProto) (io.Reader, error)
}

// Plugin is a protoc plugin.
type Plugin interface {
	Handle() error
}

// PluginOptions are options for Plugins.
type PluginOptions struct {
	// If not specified, os.Stdin is used
	In io.Reader
	// If not specified, os.Stdout is used
	Out io.Writer
	// If true, the output from Generator will be formatted before write.
	GoFmt bool
}

// NewPlugin returns a new Plugin for the Generator.
func NewPlugin(extension string, generator Generator, options PluginOptions) Plugin {
	return newPlugin(extension, generator, options)
}
