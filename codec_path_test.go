package alchemy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy"
	"github.com/wjiec/alchemy/internal/testpb"
)

func TestPathDecoder_Decoder(t *testing.T) {
	var decoder alchemy.PathDecoder

	t.Run("has path parameter", func(t *testing.T) {
		req := NewRequest(
			WithPathVars(map[string]string{"string_value": "foo"}),
		)

		if dec := decoder.Decoder(req); assert.NotNil(t, dec) {
			var message testpb.Proto3Message
			if err := dec(&message); assert.NoError(t, err) {
				assert.Equal(t, "foo", message.StringValue)
			}
		}
	})

	t.Run("no path parameter", func(t *testing.T) {
		if dec := decoder.Decoder(NewRequest()); assert.NotNil(t, dec) {
			var message testpb.Proto3Message
			if err := dec(&message); assert.NoError(t, err) {
				assert.Empty(t, message.StringValue)
			}
		}
	})

	t.Run("type mismatched", func(t *testing.T) {
		req := NewRequest(
			WithPathVars(map[string]string{"int32_value": "foo"}),
		)

		dec := decoder.Decoder(req)
		assert.ErrorContains(t, dec(&testpb.Proto3Message{}), "type mismatch")
	})
}
