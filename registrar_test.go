package alchemy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy"
)

type EchoRequest struct{ Text string }
type EchoResponse struct{ Text string }

type EchoServiceServer interface {
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
}

type FakeEchoServiceImpl struct{}

func (FakeEchoServiceImpl) Echo(_ context.Context, req *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{Text: req.Text}, nil
}

func RegisterEchoServiceServer(s alchemy.ServiceRegistrar, srv EchoServiceServer) {
	s.RegisterService(&alchemy.ServiceDesc{}, srv)
}

func TestWithServiceRegister(t *testing.T) {
	assert.NotNil(t, alchemy.WithServiceRegister[EchoServiceServer](RegisterEchoServiceServer, &FakeEchoServiceImpl{}))
}
