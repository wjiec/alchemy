package bizerr_test

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wjiec/alchemy/bizerr"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, bizerr.New(10000, http.StatusBadRequest, "bad user"))
}

func TestFromError(t *testing.T) {
	t.Run("simple error", func(t *testing.T) {
		bizErr, ok := bizerr.FromError(errors.New("something wrong"))
		assert.False(t, ok)
		assert.Nil(t, bizErr)
	})

	t.Run("bizerr", func(t *testing.T) {
		var rawErr error = bizerr.New(10000, http.StatusBadRequest, "bad user")

		bizErr, ok := bizerr.FromError(rawErr)
		if assert.True(t, ok) && assert.NotNil(t, bizErr) {
			assert.Equal(t, rawErr.(*bizerr.Error).Code(), bizErr.Code())
			assert.Equal(t, rawErr.(*bizerr.Error).Status(), bizErr.Status())
		}
	})

	t.Run("status error", func(t *testing.T) {
		bizErr, ok := bizerr.FromError(status.Error(codes.InvalidArgument, "invalid argument"))
		assert.False(t, ok)
		assert.Nil(t, bizErr)
	})
}

func TestError_Code(t *testing.T) {
	bizErr := bizerr.ThousandStep.Next().New(http.StatusBadRequest, "bad request")
	assert.NotZero(t, bizErr.Code())
}

func TestError_Status(t *testing.T) {
	bizErr := bizerr.ThousandStep.Next().New(http.StatusBadRequest, "bad request")
	assert.NotZero(t, bizErr.Status())
}

func TestError_Error(t *testing.T) {
	bizErr := bizerr.ThousandStep.Next().New(http.StatusBadRequest, "bad request")
	assert.NotEmpty(t, bizErr.Error())
}

func TestError_GRPCStatus(t *testing.T) {
	bizErr := bizerr.ThousandStep.Next().New(http.StatusBadRequest, "bad request")
	assert.NotPanics(t, func() {
		assert.NotNil(t, bizErr.GRPCStatus())
	})
}

func TestError_Equals(t *testing.T) {
	bizErr1 := bizerr.ThousandStep.Next().New(http.StatusBadRequest, "bad request")
	bizErr2 := bizerr.ThousandStep.Next().New(http.StatusBadRequest, "bad request")

	assert.False(t, bizErr1.Equals(bizErr2))
	assert.False(t, bizErr1.Equals(bizErr2.With(errors.New("another error"))))
}

func TestError_With(t *testing.T) {
	bizErr := bizerr.ThousandStep.Next().New(http.StatusBadRequest, "bad request")
	bizErr1 := bizErr.With(errors.New("another error"))
	if assert.NotNil(t, bizErr1) {
		t.Log(bizErr.GRPCStatus().Err())
		t.Log(bizErr1.GRPCStatus().Err())
	}
}
