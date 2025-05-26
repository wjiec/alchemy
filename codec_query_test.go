package alchemy_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy"
	"github.com/wjiec/alchemy/internal/testpb"
)

func TestQueryDecoder_Decoder(t *testing.T) {
	var decoder alchemy.QueryDecoder

	t.Run("correct", func(t *testing.T) {
		req := NewRequest(WithQuery(url.Values{
			"string_value":    []string{"foo"},
			"int32_value":     []string{"42"},
			"double_value":    []string{"3.14"},
			"repeated_string": []string{"one", "two"},
		}))

		if dec := decoder.Decoder(req); assert.NotNil(t, dec) {
			var message testpb.Proto3Message
			if err := dec(&message); assert.NoError(t, err) {
				assert.Equal(t, "foo", message.StringValue)
				assert.Equal(t, int32(42), message.Int32Value)
				assert.Equal(t, 3.14, message.DoubleValue)
				assert.Contains(t, message.RepeatedString, "one")
				assert.Contains(t, message.RepeatedString, "two")
			}
		}
	})

	t.Run("with path parameter", func(t *testing.T) {
		req := NewRequest(WithQuery(url.Values{
			"string_value": []string{"foo"},
			"int32_value":  []string{"42"},
		}), WithRouteDesc(&alchemy.RouteDesc{PathParameters: []string{"string_value"}}))

		if dec := decoder.Decoder(req); assert.NotNil(t, dec) {
			var message testpb.Proto3Message
			if err := dec(&message); assert.NoError(t, err) {
				assert.Empty(t, message.StringValue)
				assert.Equal(t, int32(42), message.Int32Value)
			}
		}
	})
}
