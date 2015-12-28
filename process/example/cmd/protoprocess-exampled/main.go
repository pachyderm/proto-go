package main

import (
	"os"
	"path/filepath"

	"go.pedge.io/env"
	"go.pedge.io/pkg/archive"
	"go.pedge.io/proto/process"
	"go.pedge.io/proto/server"
	"go.pedge.io/protolog"
	"google.golang.org/grpc"
)

type appEnv struct {
	Port uint16 `env:"PORT,default=1678"`
}

func main() {
	env.Main(do, &appEnv{})
}

func do(appEnvObj interface{}) error {
	appEnv := appEnvObj.(*appEnv)
	protolog.SetLevel(protolog.LevelDebug)
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
	protolog.Debugf("processing in %s", dirPath)
	return filepath.Walk(
		dirPath,
		func(path string, info os.FileInfo, err error) (retErr error) {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			protolog.Debugf("processing %s", path)
			file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, info.Mode()&os.ModePerm)
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
