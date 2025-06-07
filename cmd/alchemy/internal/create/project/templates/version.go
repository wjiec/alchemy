package templates

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

type Version struct{}

func (v *Version) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (v *Version) Path() string { return "VERSION" }
func (v *Version) Body() string { return `v0.0.1` }
