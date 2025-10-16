package alchemy_test

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/wjiec/alchemy"
)

type FlagAddr struct{ addr string }

func (a *FlagAddr) Network(_ context.Context) string { return "tcp" }
func (a *FlagAddr) String(_ context.Context) string  { return a.addr }

func TestWithGrpcServer(t *testing.T) {
	t.Run("fixed listener", func(t *testing.T) {
		app, err := alchemy.New(t.Name(),
			alchemy.WithGrpcServer(alchemy.TCP(":0")),
		)

		assert.NoError(t, err)
		assert.NotNil(t, app)
	})

	t.Run("dynamic listener", func(t *testing.T) {
		var flagAddr FlagAddr
		app, err := alchemy.New(t.Name(),
			alchemy.WithGrpcServer(&flagAddr),
			alchemy.WithBeforeStart(func(ctx context.Context, root *cobra.Command) error {
				root.PersistentFlags().StringVar(&flagAddr.addr, "grpc-addr", ":8080", "")
				return nil
			}),
		)

		assert.NoError(t, err)
		assert.NotNil(t, app)
	})
}

func TestGrpcWithReflection(t *testing.T) {
	app, err := alchemy.New(t.Name(),
		alchemy.WithGrpcServer(alchemy.TCP(":0"),
			alchemy.GrpcWithReflection(true),
		),
	)

	assert.NoError(t, err)
	assert.NotNil(t, app)
}

func TestGrpcWithServerOption(t *testing.T) {
	app, err := alchemy.New(t.Name(),
		alchemy.WithGrpcServer(alchemy.TCP(":0"),
			alchemy.GrpcWithServerOption(
				grpc.MaxSendMsgSize(1<<20),
				grpc.MaxRecvMsgSize(1<<20),
			),
		),
	)

	assert.NoError(t, err)
	assert.NotNil(t, app)
}

func TestGrpcWithServices(t *testing.T) {
	type Fake struct{}

	app, err := alchemy.New(t.Name(),
		alchemy.WithGrpcServer(alchemy.TCP(":0"),
			alchemy.GrpcWithServices(func(s grpc.ServiceRegistrar) {
				s.RegisterService(&grpc.ServiceDesc{}, &Fake{})
			}),
		),
	)

	assert.NoError(t, err)
	assert.NotNil(t, app)
}

func TestGrpcServer_Start(t *testing.T) {
	t.Run("no service", func(t *testing.T) {
		app, err := alchemy.New(t.Name(), alchemy.WithGrpcServer(alchemy.TCP(":0")))
		if assert.NoError(t, err) && assert.NotNil(t, app) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			assert.NoError(t, app.Start(ctx))
		}
	})
}
