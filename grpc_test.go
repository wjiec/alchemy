package alchemy_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/wjiec/alchemy"
)

func TestWithGrpcServer(t *testing.T) {
	app, err := alchemy.New(t.Name(),
		alchemy.WithGrpcServer("tcp", ":0"),
	)

	assert.NoError(t, err)
	assert.NotNil(t, app)
}

func TestGrpcWithReflection(t *testing.T) {
	app, err := alchemy.New(t.Name(),
		alchemy.WithGrpcServer("tcp", ":0",
			alchemy.GrpcWithReflection(true),
		),
	)

	assert.NoError(t, err)
	assert.NotNil(t, app)
}

func TestGrpcWithServerOption(t *testing.T) {
	app, err := alchemy.New(t.Name(),
		alchemy.WithGrpcServer("tcp", ":0",
			alchemy.GrpcWithServerOption(
				grpc.MaxSendMsgSize(1<<20),
				grpc.MaxRecvMsgSize(1<<20),
			),
		),
	)

	assert.NoError(t, err)
	assert.NotNil(t, app)
}

func TestGrpcServer_Start(t *testing.T) {
	t.Run("no service", func(t *testing.T) {
		app, err := alchemy.New(t.Name(), alchemy.WithGrpcServer("tcp", ":0"))
		if assert.NoError(t, err) && assert.NotNil(t, app) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			assert.NoError(t, app.Start(ctx))
		}
	})
}
