package protoprocess

import "io"

type tarCompressor struct {
	excludePatternsFiles []string
}

func newTarCompressor(opts CompressorOptions) *tarCompressor {
	excludePatternsFiles := opts.ExcludePatternsFiles
	if excludePatternsFiles == nil {
		excludePatternsFiles = DefaultExcludePatternsFiles
	}
	return &tarCompressor{excludePatternsFiles}
}

func (c *tarCompressor) Compress(dirPath string) (io.ReadCloser, error) {
	return nil, nil
}

func (c *tarCompressor) Decompress(reader io.Reader, dirPath string) error {
	return nil
}
