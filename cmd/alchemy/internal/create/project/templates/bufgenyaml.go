package templates

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/internal/template"
)

type BufGenYaml struct{}

func (b *BufGenYaml) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (b *BufGenYaml) Path() string { return "buf.gen.yaml" }
func (b *BufGenYaml) Body() string { return bufGenYamlTemplate }

const bufGenYamlTemplate = `version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: {{ flag "repo" }}/api
  disable:
{{- if eq (flag "with-validate") "true" }}
    - file_option: go_package
      module: buf.build/bufbuild/protovalidate
{{- end }}
    - file_option: go_package
      module: buf.build/googleapis/googleapis
plugins:
  - remote: buf.build/protocolbuffers/go
    out: api
    opt: paths=source_relative
  - remote: buf.build/grpc/go
    out: api
    opt:
      - paths=source_relative
  - local: protoc-gen-alchemy
    out: api
    opt:
      - paths=source_relative
inputs:
  - directory: api`
