package alchemy

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// WithGrpcServer sets the gRPC server for the App, enabling the application to handle
// gRPC requests using the specified server configuration.
func WithGrpcServer(addr Addr, options ...GrpcOption) AppOption {
	return func(app *App) error {
		app.grpcServer = &grpcServer{addr: addr}
		for _, applyGrpcOption := range options {
			if err := applyGrpcOption(app.grpcServer); err != nil {
				return err
			}
		}

		app.grpcServer.unaryInterceptor = app.wrapGrpcUnaryInterceptor()
		return nil
	}
}

// grpcServer represents a gRPC server.
type grpcServer struct {
	addr       Addr
	reflection bool

	options          []grpc.ServerOption
	services         []func(grpc.ServiceRegistrar)
	unaryInterceptor grpc.UnaryServerInterceptor
}

// Start initiates the gRPC server and begins serving requests.
func (gs *grpcServer) Start(ctx context.Context) error {
	l, err := net.Listen(gs.addr.Network(ctx), gs.addr.String(ctx))
	if err != nil {
		return err
	}

	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.ChainUnaryInterceptor(gs.unaryInterceptor))
	grpcOptions = append(grpcOptions, gs.options...)

	server := grpc.NewServer(grpcOptions...)
	for _, registerService := range gs.services {
		registerService(server)
	}

	// Enable server reflection if specified.
	if gs.reflection {
		reflection.Register(server)
	}

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	// Start serving on the configured listener.
	return server.Serve(l)
}

// GrpcOption used to configure a gRPC server instance.
type GrpcOption func(server *grpcServer) error

// GrpcWithReflection enables or disables server reflection based on the provided boolean flag
func GrpcWithReflection(reflection bool) GrpcOption {
	return func(server *grpcServer) error {
		server.reflection = reflection
		return nil
	}
}

// GrpcWithServerOption appends fallback [grpc.ServerOption] to the gRPC server's configuration.
func GrpcWithServerOption(options ...grpc.ServerOption) GrpcOption {
	return func(server *grpcServer) error {
		server.options = append(server.options, options...)
		return nil
	}
}

// GrpcWithServices registers a provided service with the gRPC server.
func GrpcWithServices(register func(grpc.ServiceRegistrar)) GrpcOption {
	return func(server *grpcServer) error {
		server.services = append(server.services, register)
		return nil
	}
}
