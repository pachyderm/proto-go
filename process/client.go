package protoprocess

import (
	"bytes"

	"go.pedge.io/pkg/archive"
	"go.pedge.io/proto/stream"
)

type client struct {
	archiver pkgarchive.Archiver
	opts     ClientOptions
}

func newClient(archiver pkgarchive.Archiver, opts ClientOptions) *client {
	return &client{archiver, opts}
}

func (c *client) Process(dirPath string, streamingBytesDuplexCloser protostream.StreamingBytesDuplexCloser) (retErr error) {
	readCloser, err := c.archiver.Compress(dirPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := readCloser.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	if err := send(readCloser, streamingBytesDuplexCloser, c.getChunkSizeBytes()); err != nil {
		return err
	}
	if err := streamingBytesDuplexCloser.CloseSend(); err != nil {
		return err
	}
	buffer := bytes.NewBuffer(nil)
	if err := recv(buffer, streamingBytesDuplexCloser, c.getChunkSizeBytes()); err != nil {
		return err
	}
	return c.archiver.Decompress(buffer, dirPath)
}

func (c *client) getChunkSizeBytes() int {
	if c.opts.ChunkSizeBytes > 0 {
		return c.opts.ChunkSizeBytes
	}
	return DefaultChunkSizeBytes
}
