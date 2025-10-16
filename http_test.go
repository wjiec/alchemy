package alchemy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/wjiec/alchemy"
	"github.com/wjiec/alchemy/bizerr"
)

func TestWithHttpServer(t *testing.T) {
	app, err := alchemy.New(t.Name(),
		alchemy.WithHttpServer(alchemy.TCP(":0")),
	)

	assert.NoError(t, err)
	assert.NotNil(t, app)
}

func TestHttpWithAdditionalHandler(t *testing.T) {
	pong := func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("pong")) }
	assert.NotNil(t, alchemy.HttpWithAdditionalHandler(http.MethodGet, "/ping", http.HandlerFunc(pong)))
}

func TestHttpWithGracefulShutdownTimeout(t *testing.T) {
	assert.NotNil(t, alchemy.HttpWithGracefulShutdownTimeout(time.Second))
}

func TestNewContextWithHttpRequest(t *testing.T) {
	assert.NotNil(t, alchemy.NewContextWithHttpRequest(context.Background(), &http.Request{}))
}

func TestHttpRequestFromContext(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		req := new(http.Request)
		ctx := alchemy.NewContextWithHttpRequest(context.Background(), req)
		if assert.NotNil(t, ctx) {
			req1, found := alchemy.HttpRequestFromContext(ctx)
			if assert.True(t, found) && assert.NotNil(t, req1) {
				assert.Equal(t, req, req1)
			}
		}
	})

	t.Run("not exists", func(t *testing.T) {
		req, found := alchemy.HttpRequestFromContext(context.Background())
		assert.False(t, found)
		assert.Nil(t, req)
	})
}

func TestNewContextWithHttpResponseWriter(t *testing.T) {
	assert.NotNil(t, alchemy.NewContextWithHttpResponseWriter(context.Background(), httptest.NewRecorder()))
}

func TestHttpResponseWriterFromContext(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx := alchemy.NewContextWithHttpResponseWriter(context.Background(), w)
		if assert.NotNil(t, ctx) {
			w1, found := alchemy.HttpResponseWriterFromContext(ctx)
			if assert.True(t, found) && assert.NotNil(t, w1) {
				assert.Equal(t, w, w1)
			}
		}
	})

	t.Run("not exists", func(t *testing.T) {
		req, found := alchemy.HttpResponseWriterFromContext(context.Background())
		assert.False(t, found)
		assert.Nil(t, req)
	})
}

type NoopCodec struct{}

func (NoopCodec) Decoder(*http.Request) func(any) error { return nil }
func (NoopCodec) Encoder(http.ResponseWriter, *http.Request) func(any) ([]byte, error) {
	return func(any) ([]byte, error) {
		return nil, nil
	}
}

func TestHttpWithCodecFactory(t *testing.T) {
	assert.NotNil(t, alchemy.HttpWithCodecFactory(&NoopCodec{}))
}

func TestHttpWithNotFoundHandler(t *testing.T) {
	assert.NotNil(t, alchemy.HttpWithNotFoundHandler(func(ctx context.Context, req *http.Request) (any, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}))
}

func TestHttpWithMethodNotAllowedHandler(t *testing.T) {
	assert.NotNil(t, alchemy.HttpWithMethodNotAllowedHandler(func(ctx context.Context, req *http.Request) (any, error) {
		return nil, status.Error(codes.Unimplemented, "unimplemented")
	}))
}

func TestHttpWithErrorHandler(t *testing.T) {
	assert.NotNil(t, alchemy.HttpWithErrorHandler(func(ctx context.Context, req *http.Request, err error) error {
		if _, ok := bizerr.FromError(err); !ok {
			err = bizerr.New(10000, http.StatusInternalServerError, err.Error())
		}
		return err
	}))
}

func TestHttpWithResponseDecorator(t *testing.T) {
	assert.NotNil(t, alchemy.HttpWithResponseDecorator(func(resp any) any {
		return map[string]any{"data": resp}
	}))
}

func TestHttpWihMetadataAnnotator(t *testing.T) {
	assert.NotNil(t, alchemy.HttpWihMetadataAnnotator(func(ctx context.Context, req *http.Request, md metadata.MD) {
		for _, cookie := range req.Cookies() {
			md.Set(cookie.Name, cookie.Value)
		}
	}))
}

func TestHttpWithOutgoingHeaderMatcher(t *testing.T) {
	const MetadataHeaderPrefix = "X-Http-Header"
	assert.NotNil(t, alchemy.HttpWithOutgoingHeaderMatcher(func(key string) (string, bool) {
		if strings.HasPrefix(key, MetadataHeaderPrefix) {
			return key[len(MetadataHeaderPrefix):], true
		}
		return "", false
	}))
}

func TestSendHeader(t *testing.T) {
	t.Run("no metadata", func(t *testing.T) {
		assert.NotPanics(t, func() {
			alchemy.SendHeader(context.Background(), "Location", "http://localhost/home")
		})
	})
}
