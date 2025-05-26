package bizerr_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wjiec/alchemy/bizerr"
)

func TestNewZone(t *testing.T) {
	assert.NotNil(t, bizerr.NewZone(10000))
	assert.NotNil(t, bizerr.NewZone(10000, 20000))
}

func TestZone_New(t *testing.T) {
	t.Run("unlimited", func(t *testing.T) {
		z := bizerr.NewZone(10000)
		require.NotNil(t, z)

		assert.NotNil(t, z.New(http.StatusBadRequest, "bad user"))
		assert.NotNil(t, z.New(http.StatusBadRequest, "bad pass"))
	})

	t.Run("limited", func(t *testing.T) {
		z := bizerr.NewZone(10000, 10001)
		require.NotNil(t, z)

		assert.NotNil(t, z.New(http.StatusUnauthorized, "bad user"))
		assert.Panics(t, func() {
			assert.NotNil(t, z.New(http.StatusUnauthorized, "bad pass"))
		})
	})
}

func TestStep_Reset(t *testing.T) {
	assert.NotPanics(t, func() {
		bizerr.ThousandStep.Reset(0)
		bizerr.TenThousandStep.Reset(0)
	})
}

func TestStep_Next(t *testing.T) {
	assert.NotNil(t, bizerr.ThousandStep.Next())
	assert.NotNil(t, bizerr.TenThousandStep.Next())
}
