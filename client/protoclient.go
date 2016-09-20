package protoclient // import "go.pedge.io/proto/client"

import (
	"fmt"

	"go.pedge.io/pkg/cobra"
	"go.pedge.io/proto/version"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// NewVersionCommand creates a new command to print the version of the client and server.
func NewVersionCommand(clientVersion *protoversion.Version, clientConnFunc func() (*grpc.ClientConn, error)) *cobra.Command {
	return &cobra.Command{
		Use:  "version",
		Long: "Print the version.",
		Run: pkgcobra.RunFixedArgs(0, func(args []string) error {
			clientConn, err := clientConnFunc()
			if err != nil {
				return err
			}
			serverVersion, err := protoversion.GetServerVersion(clientConn)
			if err != nil {
				return err
			}
			fmt.Printf("Client: %s\nServer: %s\n", clientVersion.VersionString(), serverVersion.VersionString())
			return nil
		}),
	}
}
