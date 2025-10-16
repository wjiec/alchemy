package alchemy_test

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/wjiec/alchemy"
)

func TestNew(t *testing.T) {
	app, err := alchemy.New("foobar")
	assert.NoError(t, err)
	assert.NotNil(t, app)
}

func TestApp_Start(t *testing.T) {
	app, err := alchemy.New("foobar")
	if assert.NoError(t, err) {
		err = app.Start(context.Background())
		assert.NoError(t, err)
	}
}

func TestWithSubCommand(t *testing.T) {
	echoCommand := &cobra.Command{
		Use: "echo",
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println(strings.Join(args, "\n"))
		},
	}

	assert.NotNil(t, alchemy.WithSubCommand(echoCommand))
}

func TestWithUnaryInterceptor(t *testing.T) {
	loggerInterceptor := func(ctx context.Context, req any, info *alchemy.UnaryServerInfo, handler alchemy.UnaryHandler) (any, error) {
		slog.Info("new request coming", "info", info)
		return handler(ctx, req)
	}

	assert.NotNil(t, alchemy.WithUnaryInterceptor(loggerInterceptor))
}

func TestWithResetUnaryInterceptors(t *testing.T) {
	assert.NotNil(t, alchemy.WithResetUnaryInterceptors())
}
