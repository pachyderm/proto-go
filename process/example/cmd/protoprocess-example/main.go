package main

import (
	"fmt"
	"os"

	"go.pedge.io/env"
	"go.pedge.io/pkg/archive"
	"go.pedge.io/proto/process"
	"google.golang.org/grpc"
)

var (
	defaultEnv = map[string]string{
		"ADDRESS": "0.0.0.0:678",
	}
)

type appEnv struct {
	Address string `env:"ADDRESS"`
}

func main() {
	env.Main(do, &appEnv{}, defaultEnv)
}

func do(appEnvObj interface{}) error {
	appEnv := appEnvObj.(*appEnv)
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: %s /path/to/dir", os.Args[1])
	}
	dirPath := os.Args[1]
	clientConn, err := grpc.Dial(appEnv.Address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	return protoprocess.NewAPIProcessor(
		protoprocess.NewClient(
			pkgarchive.NewTarArchiver(
				pkgarchive.ArchiverOptions{
					ExcludePatternsFiles: []string{
						".exampleignore",
					},
				},
			),
			protoprocess.ClientOptions{},
		),
		protoprocess.NewAPIClient(
			clientConn,
		),
	).Process(dirPath)
}
