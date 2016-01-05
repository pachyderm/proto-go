package protoserver // import "go.pedge.io/proto/server"

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"time"

	"go.pedge.io/env"
	"go.pedge.io/pkg/http"
	"go.pedge.io/proto/time"
	"go.pedge.io/proto/version"
	"go.pedge.io/protolog"

	"golang.org/x/net/context"

	"github.com/gengo/grpc-gateway/runtime"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

var (
	// ErrMustSpecifyRegisterFunc is used when a register func is nil.
	ErrMustSpecifyRegisterFunc = errors.New("must specify registerFunc")
)

// ServeEnv are environment variables for serving.
type ServeEnv struct {
	// Default is 7070.
	GRPCPort uint16 `env:"GRPC_PORT,default=7070"`
}

// GetServeEnv gets a ServeEnv using environment variables.
func GetServeEnv() (ServeEnv, error) {
	var serveEnv ServeEnv
	if err := env.Populate(&serveEnv); err != nil {
		return ServeEnv{}, err
	}
	return serveEnv, nil
}

// ServeOptions represent optional fields for serving.
type ServeOptions struct {
	Version *protoversion.Version
}

// Serve serves stuff.
func Serve(
	appName string,
	registerFunc func(*grpc.Server),
	serveEnv ServeEnv,
	options ServeOptions,
) (retErr error) {
	defer func(start time.Time) { logServerFinished(start, retErr) }(time.Now())
	if registerFunc == nil {
		return ErrMustSpecifyRegisterFunc
	}
	if serveEnv.GRPCPort == 0 {
		serveEnv.GRPCPort = 7070
	}
	appEnv, err := pkghttp.GetAppEnv()
	if err != nil {
		return err
	}
	if _, err := pkghttp.SetupAppEnv(appName, appEnv); err != nil {
		return err
	}
	grpcServer := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32))
	registerFunc(grpcServer)
	if options.Version != nil {
		protoversion.RegisterAPIServer(grpcServer, protoversion.NewAPIServer(options.Version, protoversion.APIServerOptions{}))
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serveEnv.GRPCPort))
	if err != nil {
		return err
	}
	return grpcServer.Serve(listener)
}

// ServeWithHTTPOptions represent optional fields for serving.
type ServeWithHTTPOptions struct {
	ServeOptions
	HTTPHandlerModifier func(http.Handler) (http.Handler, error)
}

// ServeWithHTTP serves stuff.
func ServeWithHTTP(
	appName string,
	registerFunc func(*grpc.Server),
	httpRegisterFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error,
	serveEnv ServeEnv,
	handlerEnv pkghttp.HandlerEnv,
	options ServeWithHTTPOptions,
) (retErr error) {
	defer func(start time.Time) { logServerFinished(start, retErr) }(time.Now())
	if registerFunc == nil || httpRegisterFunc == nil {
		return ErrMustSpecifyRegisterFunc
	}
	if serveEnv.GRPCPort == 0 {
		serveEnv.GRPCPort = 7070
	}
	if handlerEnv.Port == 0 {
		handlerEnv.Port = 8080
	}
	appEnv, err := pkghttp.GetAppEnv()
	if err != nil {
		return err
	}
	if _, err := pkghttp.SetupAppEnv(appName, appEnv); err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32))
	registerFunc(grpcServer)
	if options.Version != nil {
		protoversion.RegisterAPIServer(grpcServer, protoversion.NewAPIServer(options.Version, protoversion.APIServerOptions{}))
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serveEnv.GRPCPort))
	if err != nil {
		return err
	}
	grpcErrC := make(chan error)
	go func() { grpcErrC <- grpcServer.Serve(listener) }()

	time.Sleep(1 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", serveEnv.GRPCPort), grpc.WithInsecure())
	if err != nil {
		cancel()
		return err
	}
	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()
	mux := runtime.NewServeMux()
	if options.Version != nil {
		if err := protoversion.RegisterAPIHandler(ctx, mux, conn); err != nil {
			cancel()
			return err
		}
	}
	if err := httpRegisterFunc(ctx, mux, conn); err != nil {
		cancel()
		return err
	}
	var handler http.Handler = mux
	if options.HTTPHandlerModifier != nil {
		handler, err = options.HTTPHandlerModifier(mux)
		if err != nil {
			cancel()
			return err
		}
	}
	httpErrC := make(chan error)
	go func() { httpErrC <- pkghttp.ListenAndServeHandler(handler, handlerEnv) }()
	protolog.Info(
		&ServerStarted{
			Port:     uint32(serveEnv.GRPCPort),
			HttpPort: uint32(handlerEnv.Port),
		},
	)
	var errs []error
	grpcStopped := false
	for i := 0; i < 2; i++ {
		select {
		case grpcErr := <-grpcErrC:
			if grpcErr != nil {
				errs = append(errs, fmt.Errorf("grpc error: %s", grpcErr.Error()))
			}
			grpcStopped = true
		case httpErr := <-httpErrC:
			if httpErr != nil {
				errs = append(errs, fmt.Errorf("http error: %s", httpErr.Error()))
			}
			if !grpcStopped {
				grpcServer.Stop()
				_ = listener.Close()
				grpcStopped = true
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}
	return nil
}

func logServerFinished(start time.Time, err error) {
	if err != nil {
		protolog.Error(
			&ServerFinished{
				Error:    err.Error(),
				Duration: prototime.DurationToProto(time.Since(start)),
			},
		)
	} else {
		protolog.Info(
			&ServerFinished{
				Duration: prototime.DurationToProto(time.Since(start)),
			},
		)
	}
	glog.Flush()
}
