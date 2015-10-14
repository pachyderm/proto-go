package protoprocess

import (
	"bytes"
	"io/ioutil"
	"os"

	"go.pedge.io/pkg/archive"
	"go.pedge.io/proto/stream"
)

type server struct {
	archiver pkgarchive.Archiver
	opts     ServerOptions
}

func newServer(archiver pkgarchive.Archiver, opts ServerOptions) *server {
	return &server{archiver, opts}
}

func (s *server) Handle(processor Processor, streamingBytesDuplexer protostream.StreamingBytesDuplexer) (retErr error) {
	buffer := bytes.NewBuffer(nil)
	if err := recv(buffer, streamingBytesDuplexer, s.getChunkSizeBytes()); err != nil {
		return err
	}
	dirPath, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(dirPath); err != nil && retErr == nil {
			retErr = err
		}
	}()
	if err := s.archiver.Decompress(buffer, dirPath); err != nil {
		return err
	}
	if err := processor.Process(dirPath); err != nil {
		return err
	}
	readCloser, err := s.archiver.Compress(dirPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := readCloser.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return send(readCloser, streamingBytesDuplexer, s.getChunkSizeBytes())
}

func (s *server) getChunkSizeBytes() int {
	if s.opts.ChunkSizeBytes > 0 {
		return s.opts.ChunkSizeBytes
	}
	return DefaultChunkSizeBytes
}
