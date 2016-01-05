package protoserver // import "go.pedge.io/proto/server"

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	"go.pedge.io/proto/time"
	"go.pedge.io/proto/version"
	"go.pedge.io/protolog"

	"golang.org/x/net/context"

	"github.com/gengo/grpc-gateway/runtime"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

var (
	// ErrMustSpecifyPort is the error if no port is specified to Serve.
	ErrMustSpecifyPort = errors.New("must specify port")
	// ErrMustSpecifyRegisterFunc is the error if no registerFunc is specifed to Serve.
	ErrMustSpecifyRegisterFunc = errors.New("must specify registerFunc")
)

// ServeOptions represent optional fields for serving.
type ServeOptions struct {
	HTTPPort         uint16
	HTTPRegisterFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
	Version          *protoversion.Version
}

// Serve serves stuff.
func Serve(
	port uint16,
	registerFunc func(*grpc.Server),
	opts ServeOptions,
) (retErr error) {
	if port == 0 {
		return ErrMustSpecifyPort
	}
	if registerFunc == nil {
		return ErrMustSpecifyRegisterFunc
	}
	defer func(start time.Time) { logServerFinished(start, retErr) }(time.Now())
	s := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32))
	registerFunc(s)
	if opts.Version != nil {
		protoversion.RegisterAPIServer(s, protoversion.NewAPIServer(opts.Version, protoversion.APIServerOptions{}))
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcErrC := make(chan error)
	httpErrC := make(chan error)
	errCCount := 1
	go func() { grpcErrC <- s.Serve(listener) }()
	if opts.HTTPPort != 0 && (opts.Version != nil || opts.HTTPRegisterFunc != nil) {
		time.Sleep(1 * time.Second)
		ctx, cancel := context.WithCancel(context.Background())
		mux := runtime.NewServeMux()
		conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", port), grpc.WithInsecure())
		if err != nil {
			glog.Flush()
			cancel()
			return err
		}
		go func() {
			<-ctx.Done()
			_ = conn.Close()
		}()
		if opts.Version != nil {
			if err := protoversion.RegisterAPIHandler(ctx, mux, conn); err != nil {
				_ = conn.Close()
				glog.Flush()
				cancel()
				return err
			}
		}
		if opts.HTTPRegisterFunc != nil {
			if err := opts.HTTPRegisterFunc(ctx, mux, conn); err != nil {
				_ = conn.Close()
				glog.Flush()
				cancel()
				return err
			}
		}
		httpAddress := fmt.Sprintf(":%d", opts.HTTPPort)
		httpServer := &http.Server{
			Addr:    httpAddress,
			Handler: mux,
		}
		gracefulServer := &graceful.Server{
			Timeout: 1 * time.Second,
			ShutdownInitiated: func() {
				glog.Flush()
				cancel()
			},
			Server: httpServer,
		}
		errCCount++
		go func() { httpErrC <- gracefulServer.ListenAndServe() }()
	}
	protolog.Info(
		&ServerStarted{
			Port:     uint32(port),
			HttpPort: uint32(opts.HTTPPort),
		},
	)
	var errs []error
	grpcStopped := false
	for i := 0; i < errCCount; i++ {
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
				s.Stop()
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
}
