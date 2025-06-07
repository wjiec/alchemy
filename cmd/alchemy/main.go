package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/create"
	"github.com/wjiec/alchemy/cmd/alchemy/internal/version"
)

const (
	appName  = "alchemy"
	appShort = "A mini-framework for quickly building a HTTP & gRPC services"
)

func main() {
	root := cobra.Command{Use: appName, Short: appShort, SilenceErrors: true}
	root.AddCommand(version.Command(), create.Command())
	if err := root.ExecuteContext(context.Background()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
	}
}
