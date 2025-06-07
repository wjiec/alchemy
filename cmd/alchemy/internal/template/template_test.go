package template_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

func TestNewBuilder(t *testing.T) {
	assert.NotNil(t, template.NewBuilder(pflag.CommandLine))
}

type TestTemplate struct {
	BoolValue bool
}

func (t *TestTemplate) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (t *TestTemplate) Path() string { return "foobar" }
func (t *TestTemplate) Body() string { return "{{ if .BoolValue }}foo{{ else }}bar{{ end }}" }

func TestBuilder_AddTemplate(t *testing.T) {
	b := template.NewBuilder(pflag.CommandLine)
	require.NotNil(t, b)

	assert.NotPanics(t, func() {
		b.AddTemplate(&TestTemplate{})
	})
}

func TestBuilder_Init(t *testing.T) {
	b := template.NewBuilder(pflag.CommandLine)
	require.NotNil(t, b)

	b.AddTemplate(&TestTemplate{})
	assert.NoError(t, b.Init(&cobra.Command{}))
}

func TestBuilder_Build(t *testing.T) {
	b := template.NewBuilder(pflag.CommandLine)
	require.NotNil(t, b)

	b.AddTemplate(&TestTemplate{})

	workspace := t.TempDir()
	if err := b.Build(workspace, map[string]string{}); assert.NoError(t, err) {
		fp, err := os.Open(filepath.Join(workspace, "foobar"))
		if assert.NoError(t, err) {
			defer func() { _ = fp.Close() }()

			content, err := io.ReadAll(fp)
			if assert.NoError(t, err) {
				assert.Equal(t, "bar", string(content))
			}
		}
	}
}
