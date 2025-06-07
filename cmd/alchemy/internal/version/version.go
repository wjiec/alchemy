package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "v0.0.0 (dev)"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the alchemy version",
		Aliases: []string{"v"},
		Run: func(*cobra.Command, []string) {
			fmt.Println(Version)
		},
	}
}
