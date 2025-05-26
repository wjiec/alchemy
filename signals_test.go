package alchemy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy"
)

func TestSignals(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		assert.NotNil(t, alchemy.SetupSignalHandler())
	})

	t.Run("panic", func(t *testing.T) {
		assert.Panics(t, func() {
			alchemy.SetupSignalHandler()
		})
	})
}
