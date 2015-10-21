package main

import (
	"os"

	"go.pedge.io/env"
	"go.pedge.io/pkg/archive"
	"go.pedge.io/proto/process"
	"go.pedge.io/protolog"
	"google.golang.org/grpc"
)

var (
	defaultEnv = map[string]string{
		"ADDRESS": "0.0.0.0:1678",
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
	protolog.SetLevel(protolog.Level_LEVEL_DEBUG)
	var dirPath string
	var err error
	if len(os.Args) >= 2 {
		dirPath = os.Args[1]
	} else {
		dirPath, err = os.Getwd()
		if err != nil {
			return err
		}
	}
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