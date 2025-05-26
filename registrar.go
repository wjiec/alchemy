package alchemy

import (
	"context"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

// WithServiceRegister returns an AppOption that registers a service
// using the provided register function.
func WithServiceRegister[T any](register func(ServiceRegistrar, T), srv T) AppOption {
	return func(app *App) error {
		app.services = append(app.services, func(s ServiceRegistrar) {
			register(s, srv)
		})
		return nil
	}
}

// ServiceRegistrar wraps a single method that supports service registration.
type ServiceRegistrar interface {
	// RegisterService registers a service and its implementation to the
	// concrete type implementing this interface.
	RegisterService(desc *ServiceDesc, srv any)
}

// RegisterService registers a service and its implementation to
// the underlying gRPC and HTTP server.
func (a *App) RegisterService(desc *ServiceDesc, srv any) {
	if a.grpcServer != nil {
		a.grpcServer.services = append(a.grpcServer.services, func(s grpc.ServiceRegistrar) {
			s.RegisterService(desc.GrpcServiceDesc, srv)
		})
	}

	if a.httpServer != nil {
		a.httpServer.services = append(a.httpServer.services, func(r *mux.Router) {
			for i := range desc.Routes {
				route := &desc.Routes[i]
				r.Handle(route.PathPattern, a.httpServer.wrapHttpHandler(route, srv)).Methods(route.HttpMethod)
			}
		})
	}
}

// ServiceDesc represents a service description containing both gRPC service details
// and HTTP route information for API gateway integration.
type ServiceDesc struct {
	GrpcServiceDesc *grpc.ServiceDesc // the underlying gRPC service descriptor information
	Routes          []RouteDesc       // the HTTP route descriptions that map to this service's methods
}

// RouteDesc defines an HTTP route mapping for a gRPC method.
type RouteDesc struct {
	HttpMethod     string             // the HTTP verb (GET, POST, PUT, DELETE, etc.)
	PathPattern    string             // the URL path pattern for this route
	Handler        grpc.MethodHandler // the gRPC method handler function that processes the request
	RequestField   KeyPath            // Specifies which field in the request message should be parsed from the HTTP request body
	ResponseField  KeyPath            // Specifies which field in the response message to use as the HTTP response body
	PathParameters []string           // List of path parameters extracted from the URL path and mapped to request fields
}

// routeDescKey is how we find the [*RouteDesc] in a context.Context.
type routeDescKey struct{}

// NewContextWithRouteDesc returns a new Context, derived from ctx, which carries
// the provided [*RouteDesc].
func NewContextWithRouteDesc(ctx context.Context, desc *RouteDesc) context.Context {
	return context.WithValue(ctx, routeDescKey{}, desc)
}

// RouteDescFromContext returns a [*RouteDesc] from ctx.
func RouteDescFromContext(ctx context.Context) (*RouteDesc, bool) {
	if raw := ctx.Value(routeDescKey{}); raw != nil {
		return raw.(*RouteDesc), true
	}
	return nil, false
}

// KeyPath represents a path to access nested values within a data structure.
type KeyPath struct {
	Name     string        // the full path of the field
	Accessor func(any) any // a function that returns the nested value at the specified path
}
