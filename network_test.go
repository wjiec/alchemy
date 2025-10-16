package alchemy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy"
)

func TestNewAddr(t *testing.T) {
	assert.NotNil(t, alchemy.NewAddr("tcp", ":8080"))
}

func TestTCP(t *testing.T) {
	assert.NotNil(t, alchemy.TCP(":8080"))
}
