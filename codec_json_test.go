package alchemy_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy"
	"github.com/wjiec/alchemy/internal/testpb"
)

func TestJsonDecoder_Decoder(t *testing.T) {
	factory := alchemy.JsonDecoder{}

	t.Run("simple", func(t *testing.T) {
		req := NewRequest(
			WithBody(bytes.NewReader([]byte(`{"string_value": "foo", "int64_value": 42}`))),
			WithHeader(http.Header{"Content-Type": []string{factory.ContentType(nil)}}),
		)

		if dec := factory.Decoder(req); assert.NotNil(t, dec) {
			var message testpb.Proto3Message
			if err := dec(&message); assert.NoError(t, err) {
				assert.Equal(t, "foo", message.StringValue)
				assert.Equal(t, int64(42), message.Int64Value)
			}
		}
	})

	t.Run("nested", func(t *testing.T) {
		req := NewRequest(
			WithBody(bytes.NewReader([]byte(`{"string_value": "foo", "int64_value": 42}`))),
			WithRouteDesc(&alchemy.RouteDesc{
				RequestField: alchemy.KeyPath{
					Name:     "nested_value",
					Accessor: func(a any) any { return &a.(*testpb.Proto3Message).NestedValue },
				},
			}),
			WithHeader(http.Header{"Content-Type": []string{factory.ContentType(nil)}}),
		)

		if dec := factory.Decoder(req); assert.NotNil(t, dec) {
			var message testpb.Proto3Message
			if err := dec(&message); assert.NoError(t, err) {
				assert.Equal(t, "foo", message.NestedValue.StringValue)
				assert.Equal(t, int64(42), message.NestedValue.Int64Value)
			}
		}
	})
}
