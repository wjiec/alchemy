package templates

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

type EditorConfig struct{}

func (e *EditorConfig) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (e *EditorConfig) Path() string { return ".editorconfig" }
func (e *EditorConfig) Body() string { return editorConfigTemplate }

const editorConfigTemplate = `root = true

[*]
charset = utf-8
tab_width = 4
indent_size = 4
end_of_line = lf
indent_style = space
max_line_length = 120
insert_final_newline = true
trim_trailing_whitespace = true

[*.proto]
indent_size = 2
tab_width = 2

[{*.go,*.go2}]
indent_style = tab

[Makefile]
indent_size = tab

[{*.yaml,*.yml}]
indent_size = 2

[VERSION]
insert_final_newline = false`
