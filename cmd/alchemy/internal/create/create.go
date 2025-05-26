package create

import (
	"github.com/spf13/cobra"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/create/project"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Initialize a new project and scaffold a service",
	}

	cmd.AddCommand(project.Command())
	return cmd
}
