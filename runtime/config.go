package runtime

import (
	"crypto/tls"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/piotrkowalczuk/promgrpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	pkg_runtime "runtime"
	"time"
)

var (
	interceptor = promgrpc.NewInterceptor(promgrpc.InterceptorOpts{
		TrackPeers: true,
	})
)

func createDefaultConfig() *Config {

	config := &Config{
		GrpcAddr: nil,
		GrpcInternalAddr: &Address{
			Network: "unix",
			Addr:    "tmp/server.sock",
		},
		GatewayAddr: &Address{
			Network: "tcp",
			Addr:    ":3000",
		},
		GrpcServerOption: 				NewServerOpts(),
		GatewayDialOption:               NewDialOpts(),
		GatewayServerConfig: &HTTPServerConfig{
			ReadTimeout:  8 * time.Second,
			WriteTimeout: 8 * time.Second,
			IdleTimeout:  2 * time.Minute,
		},
		MaxConcurrentStreams:     1000,
		GatewayServerMiddlewares: nil,
	}
	if pkg_runtime.GOOS == "windows" {
		config.GrpcInternalAddr = &Address{
			Network: "tcp",
			Addr:    ":5050",
		}
	}
	return config
}

// Address represents a network end point address.
type Address struct {
	Network string
	Addr    string
}

func (a *Address) createListener() (net.Listener, error) {
	if a.Network == "unix" {
		dir := filepath.Dir(a.Addr)
		f, err := os.Stat(dir)
		if err != nil {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return nil, errors.Wrap(err, "failed to create the directory")
			}
		} else if !f.IsDir() {
			return nil, errors.Errorf("file %q already exists", dir)
		}
	}
	lis, err := net.Listen(a.Network, a.Addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to listen %s %s", a.Network, a.Addr)
	}
	return lis, nil
}

type HTTPServerConfig struct {
	TLSConfig         *tls.Config
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	TLSNextProto      map[string]func(*http.Server, *tls.Conn, http.Handler)
	ConnState         func(net.Conn, http.ConnState)
}

func (c *HTTPServerConfig) applyTo(s *http.Server) {
	s.TLSConfig = c.TLSConfig
	s.ReadTimeout = c.ReadTimeout
	s.ReadHeaderTimeout = c.ReadHeaderTimeout
	s.WriteTimeout = c.WriteTimeout
	s.IdleTimeout = c.IdleTimeout
	s.MaxHeaderBytes = c.MaxHeaderBytes
	s.TLSNextProto = c.TLSNextProto
	s.ConnState = c.ConnState
}

// Config contains configurations of gRPC and Gateway server.
type Config struct {
	GrpcAddr                        *Address
	GrpcInternalAddr                *Address
	GatewayAddr                     *Address
	Servers                         []Server
	GrpcServerUnaryInterceptors     []grpc.UnaryServerInterceptor
	GrpcServerStreamInterceptors    []grpc.StreamServerInterceptor
	GatewayServerUnaryInterceptors  []grpc.UnaryClientInterceptor
	GatewayServerStreamInterceptors []grpc.StreamClientInterceptor
	GrpcServerOption                []grpc.ServerOption
	GatewayDialOption               []grpc.DialOption
	GatewayMuxOptions               []runtime.ServeMuxOption
	GatewayServerConfig             *HTTPServerConfig
	MaxConcurrentStreams            uint32
	GatewayServerMiddlewares        []HTTPServerMiddleware
}

func (c *Config) serverOptions() []grpc.ServerOption {
	return append(
		[]grpc.ServerOption{
			grpc_middleware.WithUnaryServerChain(c.GrpcServerUnaryInterceptors...),
			grpc_middleware.WithStreamServerChain(c.GrpcServerStreamInterceptors...),
			grpc.MaxConcurrentStreams(c.MaxConcurrentStreams),
		},
		c.GrpcServerOption...,
	)
}

func (c *Config) clientOptions() []grpc.DialOption {
	return append(
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithDialer(func(a string, t time.Duration) (net.Conn, error) {
				return net.Dial(c.GrpcInternalAddr.Network, a)
			}),
			grpc.WithUnaryInterceptor(
				grpc_middleware.ChainUnaryClient(c.GatewayServerUnaryInterceptors...),
			),
			grpc.WithStreamInterceptor(
				grpc_middleware.ChainStreamClient(c.GatewayServerStreamInterceptors...),
			),
		},
		c.GatewayDialOption...,
	)
}

func NewDialOpts() []grpc.DialOption {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("failed to setup logger for grpc interceptor")
	}
	opts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}

	grpc_zap.ReplaceGrpcLogger(logger)
	streamInterceptors := grpc.StreamClientInterceptor(grpc_middleware.ChainStreamClient(
		grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
		interceptor.StreamClient(),
		grpc_zap.StreamClientInterceptor(logger, opts...),
	))

	unaryInterceptors := grpc.UnaryClientInterceptor(grpc_middleware.ChainUnaryClient(
		grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
		interceptor.UnaryClient(),
		grpc_zap.UnaryClientInterceptor(logger, opts...),
	))

	dopts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithStatsHandler(interceptor),
		grpc.WithDialer(interceptor.Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("tcp", addr, timeout)
		})),
		grpc.WithUnaryInterceptor(unaryInterceptors),
		grpc.WithStreamInterceptor(streamInterceptors),
	}
	return dopts
}

func NewServerOpts() []grpc.ServerOption {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("failed to setup logger for grpc interceptor")
	}
	opts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}

	grpc_zap.ReplaceGrpcLogger(logger)
	streamInterceptors := grpc.StreamServerInterceptor(grpc_middleware.ChainStreamServer(
		grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
		interceptor.StreamServer(),
		grpc_zap.StreamServerInterceptor(logger, opts...),
	))

	unaryInterceptors := grpc.UnaryServerInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
		interceptor.UnaryServer(),
		grpc_zap.UnaryServerInterceptor(logger, opts...),
	))

	return []grpc.ServerOption{
		grpc.StatsHandler(interceptor),
		grpc.UnaryInterceptor(unaryInterceptors),
		grpc.StreamInterceptor(streamInterceptors),
	}
}
