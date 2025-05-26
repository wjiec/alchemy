package alchemy_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy"
	"github.com/wjiec/alchemy/internal/testpb"
)

type Empty struct{}

func (Empty) Read([]byte) (int, error) { return 0, io.EOF }

func NewRequest(options ...RequestOption) *http.Request {
	builder := &RequestBuilder{RouteDesc: &alchemy.RouteDesc{}, Body: Empty{}}
	for _, applyOption := range options {
		applyOption(builder)
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost/api?"+builder.Query.Encode(), builder.Body)
	if err != nil {
		panic(err)
	}
	if len(builder.PathVars) != 0 {
		req = mux.SetURLVars(req, builder.PathVars)
	}

	req.Header = builder.Header
	return req.WithContext(alchemy.NewContextWithRouteDesc(req.Context(), builder.RouteDesc))
}

type RequestBuilder struct {
	Query     url.Values
	Body      io.Reader
	Header    http.Header
	PathVars  map[string]string
	RouteDesc *alchemy.RouteDesc
}

type RequestOption func(builder *RequestBuilder)

func WithRouteDesc(desc *alchemy.RouteDesc) RequestOption {
	return func(b *RequestBuilder) {
		b.RouteDesc = desc
	}
}

func WithQuery(query url.Values) RequestOption {
	return func(b *RequestBuilder) {
		b.Query = query
	}
}

func WithBody(body io.Reader) RequestOption {
	return func(b *RequestBuilder) {
		b.Body = body
	}
}

func WithHeader(header http.Header) RequestOption {
	return func(b *RequestBuilder) {
		b.Header = header
	}
}

func WithPathVars(pathVars map[string]string) RequestOption {
	return func(b *RequestBuilder) {
		b.PathVars = pathVars
	}
}

func TestNewHttpDynamicCodec(t *testing.T) {
	assert.NotNil(t, alchemy.NewHttpDynamicCodec())
}

func TestHttpDynamicCodec_Decoder(t *testing.T) {
	codec := alchemy.NewHttpDynamicCodec()

	t.Run("multi source", func(t *testing.T) {
		req := NewRequest(
			WithQuery(url.Values{"double_value": []string{"3.14"}, "bool_value": []string{"yes"}}),
			WithPathVars(map[string]string{"string_value": "nice", "bool_value": "false"}),
			WithBody(bytes.NewReader([]byte(`{"string_value": "foo", "int64_value": 42}`))),
			WithRouteDesc(&alchemy.RouteDesc{
				PathParameters: []string{"string_value", "bool_value"},
				RequestField: alchemy.KeyPath{
					Name:     "nested_value",
					Accessor: func(a any) any { return &a.(*testpb.Proto3Message).NestedValue },
				},
			}),
			WithHeader(http.Header{"Content-Type": []string{"application/json"}}),
		)

		if dec := codec.Decoder(req); assert.NotNil(t, dec) {
			var message testpb.Proto3Message
			if err := dec(&message); assert.NoError(t, err) {
				assert.Equal(t, false, message.BoolValue)
				assert.Equal(t, 3.14, message.DoubleValue)
				assert.Equal(t, "nice", message.StringValue)
				assert.Equal(t, "foo", message.NestedValue.StringValue)
				assert.Equal(t, int64(42), message.NestedValue.Int64Value)
			}
		}
	})
}
