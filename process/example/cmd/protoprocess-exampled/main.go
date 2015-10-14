package main

import (
	"os"
	"path/filepath"

	"go.pedge.io/env"
	"go.pedge.io/pkg/archive"
	"go.pedge.io/proto/process"
	"go.pedge.io/proto/server"
	"google.golang.org/grpc"
)

var (
	defaultEnv = map[string]string{
		"PORT": "678",
	}
)

type appEnv struct {
	Port uint16 `env:"PORT"`
}

func main() {
	env.Main(do, &appEnv{}, defaultEnv)
}

func do(appEnvObj interface{}) error {
	appEnv := appEnvObj.(*appEnv)
	return protoserver.Serve(
		appEnv.Port,
		func(s *grpc.Server) {
			protoprocess.RegisterAPIServer(
				s,
				protoprocess.NewAPIServer(
					newProcessor(),
					protoprocess.NewServer(
						pkgarchive.NewTarArchiver(
							pkgarchive.ArchiverOptions{
								ExcludePatternsFiles: []string{
									".exampleignore",
								},
							},
						),
						protoprocess.ServerOptions{},
					),
				),
			)
		},
		protoserver.ServeOptions{},
	)
}

type processor struct{}

func newProcessor() *processor {
	return &processor{}
}

func (p *processor) Process(dirPath string) error {
	return filepath.Walk(
		dirPath,
		func(path string, info os.FileInfo, err error) (retErr error) {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			file, err := os.OpenFile(path, os.O_APPEND, info.Mode()&os.ModePerm)
			if err != nil {
				return err
			}
			defer func() {
				if err := file.Close(); err != nil && retErr == nil {
					retErr = err
				}
			}()
			_, err = file.Write([]byte("\nhello"))
			return err
		},
	)
}
