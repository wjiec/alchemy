package templates

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/internal/template"
)

type BufYaml struct {
	Validate bool
}

func (b *BufYaml) Init(fs *pflag.FlagSet, _ template.FlagConstrain) error {
	fs.BoolVar(&b.Validate, "with-validate", false, "If set, add validate support for protobuf")
	return nil
}

func (b *BufYaml) Path() string { return "buf.yaml" }
func (b *BufYaml) Body() string { return bufYamlTemplate }

const bufYamlTemplate = `version: v2
modules:
  - path: api
deps:
{{- if .Validate }}
  - buf.build/bufbuild/protovalidate
{{- end }}
  - buf.build/googleapis/googleapis
lint:
  use:
    - STANDARD
breaking:
  use:
    - FILE`
