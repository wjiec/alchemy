package errs_test

import (
	"io"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy/errs"
)

func TestIgnore(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		assert.Nil(t, errs.Ignore(io.EOF, io.EOF))
	})

	t.Run("wrapped", func(t *testing.T) {
		assert.Nil(t, errs.Ignore(errors.Wrap(io.EOF, "wrapped"), io.EOF))
		assert.Nil(t, errs.Ignore(errors.Wrap(io.EOF, "wrapped"), io.EOF))
	})

	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, errs.Ignore(nil, io.EOF))
	})
}
