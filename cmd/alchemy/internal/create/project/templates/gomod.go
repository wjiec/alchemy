package templates

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

type GoMod struct {
	Repo string
}

func (g *GoMod) Init(fs *pflag.FlagSet, fc template.FlagConstrain) error {
	fs.String("repo", "", "name to use for go module")
	return fc.MarkFlagRequired("repo")
}

func (g *GoMod) Path() string { return "go.mod" }
func (g *GoMod) Body() string { return `module {{ flag "repo" }}` }
