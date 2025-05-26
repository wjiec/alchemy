package alchemy

import (
	"context"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// App represents an application with gRPC and HTTP server capabilities,
// as well as a command-line interface using Cobra for command execution.
type App struct {
	name     string
	root     *cobra.Command
	services []func(ServiceRegistrar)

	httpServer *httpServer
	grpcServer *grpcServer

	unaryInterceptors []UnaryInterceptor
}

// Start begins the execution of the application by executing the root command
// in the context provided.
func (a *App) Start(ctx context.Context) error {
	return a.root.ExecuteContext(ctx)
}

// serve starts the application's concurrent servers and handles their lifecycle,
// waiting for them to complete. An error group is used to manage lifecycle errors.
func (a *App) serve(ctx context.Context) error {
	for _, registerService := range a.services {
		registerService(a)
	}

	eg, eCtx := errgroup.WithContext(ctx)
	if a.grpcServer != nil {
		eg.Go(func() error {
			return a.grpcServer.Start(eCtx)
		})
	}
	if a.httpServer != nil {
		eg.Go(func() error {
			return a.httpServer.Start(eCtx)
		})
	}

	return eg.Wait()
}

// wrapGrpcUnaryInterceptor creates a gRPC UnaryServerInterceptor by wrapping the app's unary interceptors.
func (a *App) wrapGrpcUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// If there are no interceptors, just call the handler directly
		if len(a.unaryInterceptors) == 0 {
			return handler(ctx, req)
		}

		// Create a chain of interceptors that will eventually call the gRPC handler
		var chainedHandler grpc.UnaryHandler = func(ctx context.Context, req any) (any, error) {
			return handler(ctx, req)
		}

		// Apply interceptors in reverse order so that the first interceptor is the outermost
		for i := len(a.unaryInterceptors) - 1; i >= 0; i-- {
			interceptor, prevHandler := a.unaryInterceptors[i], chainedHandler
			chainedHandler = func(ctx context.Context, req any) (any, error) {
				return interceptor(ctx, req, info, prevHandler)
			}
		}

		return chainedHandler(ctx, req)
	}
}

// New initializes and returns a new *App instance configured with the provided name and options.
func New(name string, options ...AppOption) (*App, error) {
	app := &App{name: name}
	app.root = &cobra.Command{
		Use: name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.serve(cmd.Context())
		},
	}

	for _, applyOption := range options {
		if err := applyOption(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

// AppOption defines a function type for configuring the App instance.
type AppOption func(app *App) error

// WithSubCommand adds one or more subcommands to the root command of the App.
func WithSubCommand(cmd ...*cobra.Command) AppOption {
	return func(app *App) error {
		app.root.AddCommand(cmd...)
		return nil
	}
}

type UnaryHandler = grpc.UnaryHandler
type UnaryServerInfo = grpc.UnaryServerInfo
type UnaryInterceptor func(context.Context, any, *UnaryServerInfo, UnaryHandler) (any, error)

// WithUnaryInterceptor adds a UnaryInterceptor to the App's configuration.
func WithUnaryInterceptor(interceptor UnaryInterceptor) AppOption {
	return func(app *App) error {
		app.unaryInterceptors = append(app.unaryInterceptors, interceptor)
		return nil
	}
}
