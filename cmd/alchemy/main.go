package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wjiec/alchemy"
	"github.com/wjiec/alchemy/cmd/alchemy/internal/create"
	"github.com/wjiec/alchemy/cmd/alchemy/internal/version"
)

func main() {
	root := cobra.Command{
		Use:           "alchemy",
		Short:         "A mini-framework for quickly building HTTP & gRPC services",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.AddCommand(version.Command(), create.Command())
	if err := root.ExecuteContext(alchemy.SetupSignalHandler()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
	}
}

/*
alchemy new project <name>
>> package: github.com/xxx/<name>

alchemy new api echo --version v1
>>


*/
