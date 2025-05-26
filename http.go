package alchemy

import (
	"context"
	"net"
	"net/http"
	"net/textproto"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/wjiec/alchemy/bizerr"
	"github.com/wjiec/alchemy/errs"
)

const (
	DefaultGracefulShutdownTimeout = 3 * time.Second
)

// WithHttpServer sets the HTTP server for the App, enabling the application to handle
// HTTP requests using the specified server configuration.
func WithHttpServer(network, address string, options ...HttpOption) AppOption {
	return func(app *App) error {
		app.httpServer = &httpServer{
			network: network,
			address: address,

			codec:    NewHttpDynamicCodec(),
			fallback: mux.NewRouter(),

			unaryInterceptor:      app.wrapGrpcUnaryInterceptor(),
			gracefulTimeout:       DefaultGracefulShutdownTimeout,
			outgoingHeaderMatcher: DefaultOutgoingHeaderMatcher,
		}
		for _, applyHttpOption := range options {
			if err := applyHttpOption(app.httpServer); err != nil {
				return err
			}
		}

		return nil
	}
}

// grpcServer represents a HTTP server.
type httpServer struct {
	network string
	address string

	codec    CodecFactory
	fallback *mux.Router

	services              []func(*mux.Router)
	gracefulTimeout       time.Duration
	unaryInterceptor      grpc.UnaryServerInterceptor
	errorHandlers         []HttpErrorHandler
	respDecorators        []HttpResponseDecorator
	metadataAnnotators    []HttpMetadataAnnotator
	outgoingHeaderMatcher HttpOutgoingHeaderMatcher
}

// Start initiates the HTTP server and begins serving requests.
func (hs *httpServer) Start(ctx context.Context) error {
	l, err := net.Listen(hs.network, hs.address)
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	for _, registerService := range hs.services {
		registerService(router)
	}
	router.NotFoundHandler = hs.fallback

	return errs.Ignore(hs.serve(ctx, l, router), http.ErrServerClosed)
}

// serve starts the HTTP server using the provided listener and handler.
//
// It also sets up graceful shutdown triggered by the context's cancellation.
func (hs *httpServer) serve(ctx context.Context, l net.Listener, h http.Handler) error {
	server := http.Server{Handler: h}

	go func() {
		<-ctx.Done()
		shutdownCtx, forceShutdown := context.WithTimeout(context.Background(), hs.gracefulTimeout)
		defer forceShutdown()

		_ = server.Shutdown(shutdownCtx)
	}()

	return server.Serve(l)
}

type outgoingMetadataKey struct{}

// wrapHttpHandler creates an HTTP handler that wraps a gRPC method handler.
func (hs *httpServer) wrapHttpHandler(route *RouteDesc, srv any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() { _ = req.Body.Close() }()
		ctx := NewContextWithRouteDesc(req.Context(), route)
		ctx = NewContextWithHttpRequest(ctx, req)
		ctx = NewContextWithHttpResponseWriter(ctx, w)

		incomingMetadata := metadata.MD{}
		for _, annotator := range hs.metadataAnnotators {
			annotator(ctx, req, incomingMetadata)
		}
		ctx = metadata.NewIncomingContext(ctx, incomingMetadata)

		outgoingMetadata := metadata.MD{}
		ctx = context.WithValue(ctx, outgoingMetadataKey{}, outgoingMetadata)
		ctx = metadata.NewOutgoingContext(ctx, outgoingMetadata)

		req = req.WithContext(ctx)
		resp, err := route.Handler(srv, ctx, hs.codec.Decoder(req), hs.unaryInterceptor)
		hs.writeResponse(ctx, w, req, resp, err)
	})
}

// writeResponse writes the response data to the HTTP response writer.
func (hs *httpServer) writeResponse(ctx context.Context, w http.ResponseWriter, req *http.Request, resp any, err error) {
	if err == nil {
		if enc := hs.codec.Encoder(w, req); enc != nil {
			var buf []byte
			if buf, err = enc(resp); err == nil {
				hs.forwardResponseServerMetadata(ctx, w)
				_, _ = w.Write(buf)
			}
		}
	}

	if err != nil {
		hs.defaultErrorHandler(ctx, w, req, err)
		return
	}
}

// defaultErrorHandler processes and writes error responses for HTTP requests.
func (hs *httpServer) defaultErrorHandler(ctx context.Context, w http.ResponseWriter, req *http.Request, err error) {
	// Apply all registered error handlers in sequence
	for _, errHandler := range hs.errorHandlers {
		err = errHandler(req.Context(), req, err)
	}

	// Set appropriate HTTP status code from biz error if present
	if bizErr, ok := bizerr.FromError(err); ok {
		w.WriteHeader(int(bizErr.Status()))
	}

	hs.forwardResponseServerMetadata(ctx, w)
	if enc := hs.codec.Encoder(w, req); enc != nil {
		statusErr := status.Convert(err)
		statusErr.Details()

		if buf, wErr := enc(status.Convert(err).Proto()); wErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"code": 13, "message": "internal server error"}`))
		} else {
			_, _ = w.Write(buf)
		}
	}
}

// forwardResponseServerMetadata forwards gRPC metadata from the outgoing context to HTTP response headers.
// It applies header matchers to determine which metadata should be forwarded and how header names should be mapped.
func (hs *httpServer) forwardResponseServerMetadata(ctx context.Context, w http.ResponseWriter) {
	if outgoingMetadata, ok := metadata.FromOutgoingContext(ctx); ok {
		for key, values := range outgoingMetadata {
			if header, forwarded := hs.outgoingHeaderMatcher(key); forwarded {
				for _, value := range values {
					w.Header().Add(header, value)
				}
			}
		}
	}
}

// HttpOption used to configures a httpServer instance.
type HttpOption func(server *httpServer) error

// HttpWithAdditionalHandler  adds an HTTP handler for a specific method and path pattern.
func HttpWithAdditionalHandler(method, pattern string, handler http.Handler) HttpOption {
	return func(server *httpServer) error {
		server.fallback.Handle(pattern, handler).Methods(method)
		return nil
	}
}

// HttpWithGracefulShutdownTimeout configures the timeout duration for graceful server shutdown.
//
// This option sets how long the server will wait for existing connections
// to complete before forcefully terminating them.
func HttpWithGracefulShutdownTimeout(timeout time.Duration) HttpOption {
	return func(server *httpServer) error {
		server.gracefulTimeout = timeout
		return nil
	}
}

// httpRequestContextKey is how we find the [*http.Request] in a context.Context.
type httpRequestContextKey struct{}

// NewContextWithHttpRequest returns a new Context, derived from ctx, which carries
// the provided [*http.Request].
func NewContextWithHttpRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, httpRequestContextKey{}, req)
}

// HttpRequestFromContext returns a [*http.Request] from ctx.
func HttpRequestFromContext(ctx context.Context) (*http.Request, bool) {
	if raw := ctx.Value(httpRequestContextKey{}); raw != nil {
		return raw.(*http.Request), true
	}
	return nil, false
}

// httpResponseWriterKey is how we find the [http.ResponseWriter] in a context.Context.
type httpResponseWriterKey struct{}

// NewContextWithHttpResponseWriter returns a new Context, derived from ctx, which carries
// the provided [http.ResponseWriter].
func NewContextWithHttpResponseWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, httpResponseWriterKey{}, w)
}

// HttpResponseWriterFromContext returns a [http.ResponseWriter] from ctx.
func HttpResponseWriterFromContext(ctx context.Context) (http.ResponseWriter, bool) {
	if raw := ctx.Value(httpResponseWriterKey{}); raw != nil {
		return raw.(http.ResponseWriter), true
	}
	return nil, false
}

// HttpWithCodecFactory configures the codec factory used by the application.
func HttpWithCodecFactory(factory CodecFactory) HttpOption {
	return func(hs *httpServer) error {
		hs.codec = factory
		return nil
	}
}

// HttpRoutingErrorHandler defines a function type for handling HTTP routing errors.
type HttpRoutingErrorHandler func(ctx context.Context, req *http.Request) (any, error)

// serve creates an HTTP handler from the HttpRoutingErrorHandler function.
func (f HttpRoutingErrorHandler) serve(writeResponse func(context.Context, http.ResponseWriter, *http.Request, any, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() { _ = req.Body.Close() }()

		ctx := NewContextWithHttpRequest(req.Context(), req)
		ctx = NewContextWithHttpResponseWriter(ctx, w)
		resp, err := f(ctx, req.WithContext(ctx))
		writeResponse(req.Context(), w, req, resp, err)
	})
}

// HttpWithNotFoundHandler returns an HttpOption that configures the not-found handler.
//
// This handler is called when no route matches the requested URL path.
func HttpWithNotFoundHandler(handler HttpRoutingErrorHandler) HttpOption {
	return func(hs *httpServer) error {
		hs.fallback.NotFoundHandler = handler.serve(hs.writeResponse)
		return nil
	}
}

// HttpWithMethodNotAllowedHandler returns an HttpOption that configures the method-not-allowed handler.
//
// This handler is called when a route matches the URL path but doesn't support the requested HTTP method.
func HttpWithMethodNotAllowedHandler(handler HttpRoutingErrorHandler) HttpOption {
	return func(hs *httpServer) error {
		hs.fallback.MethodNotAllowedHandler = handler.serve(hs.writeResponse)
		return nil
	}
}

// HttpErrorHandler defines a function type for processing and transforming errors in HTTP handlers.
//
// This allows for custom error handling, logging, and error transformation before the response is sent.
type HttpErrorHandler func(context.Context, *http.Request, error) error

// HttpWithErrorHandler returns an HttpOption that adds an error handler to the HTTP server.
//
// Multiple error handlers can be added, and they will be executed in the order they were added.
func HttpWithErrorHandler(handler HttpErrorHandler) HttpOption {
	return func(hs *httpServer) error {
		hs.errorHandlers = append(hs.errorHandlers, handler)
		return nil
	}
}

// HttpResponseDecorator defines a function type for transforming HTTP responses before they are encoded.
//
// It takes the original response value and returns a potentially modified response value.
// This allows for wrapping, augmenting, or completely changing the response structure.
type HttpResponseDecorator func(resp any) any

// HttpWithResponseDecorator returns an HttpOption that adds a response decorator to the HTTP server.
//
// Multiple decorators can be added, and they will be executed in the order they were added.
// Each decorator can transform the response before it's encoded and sent to the client.
func HttpWithResponseDecorator(decorator HttpResponseDecorator) HttpOption {
	return func(hs *httpServer) error {
		hs.respDecorators = append(hs.respDecorators, decorator)
		return nil
	}
}

// HttpMetadataAnnotator defines a function type for annotating gRPC metadata from HTTP requests.
//
// This allows for extracting information from HTTP requests and adding it to gRPC metadata.
type HttpMetadataAnnotator func(context.Context, *http.Request, metadata.MD)

// HttpWihMetadataAnnotator returns an HttpOption that adds a metadata annotator to the HTTP server.
// Multiple annotators can be added, and they will be executed in the order they were added.
// Each annotator can modify the metadata that will be passed to the gRPC handler.
func HttpWihMetadataAnnotator(annotator HttpMetadataAnnotator) HttpOption {
	return func(hs *httpServer) error {
		hs.metadataAnnotators = append(hs.metadataAnnotators, annotator)
		return nil
	}
}

// HttpOutgoingHeaderMatcher defines a function type for determining which gRPC metadata keys
// should be forwarded to HTTP responses and how they should be mapped.
type HttpOutgoingHeaderMatcher func(string) (string, bool)

// DefaultOutgoingHeaderMatcher provides the default logic for determining which gRPC metadata
// should be forwarded to HTTP responses and how they should be mapped to HTTP headers.
func DefaultOutgoingHeaderMatcher(key string) (string, bool) {
	switch key = textproto.CanonicalMIMEHeaderKey(key); {
	case slices.Contains([]string{"Cache-Control", "Cookie", "Pragma"}, key):
		return key, true
	case strings.HasPrefix(key, runtime.MetadataHeaderPrefix):
		return key[len(runtime.MetadataHeaderPrefix):], true
	}
	return "", false
}

// HttpWithOutgoingHeaderMatcher returns an HttpOption that sets a header matcher to the HTTP server.
func HttpWithOutgoingHeaderMatcher(matcher HttpOutgoingHeaderMatcher) HttpOption {
	return func(hs *httpServer) error {
		hs.outgoingHeaderMatcher = matcher
		return nil
	}
}

// SendHeader adds a header to the outgoing metadata in the context.
//
// This header will be included in the HTTP response when forwarded by the server.
func SendHeader(ctx context.Context, header string, values ...string) {
	if raw := ctx.Value(outgoingMetadataKey{}); raw != nil {
		outgoingMetadata := raw.(metadata.MD)
		outgoingMetadata.Set(runtime.MetadataHeaderPrefix+header, values...)
	}
}
