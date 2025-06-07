package template

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

// Builder is responsible for rendering templates to files in a workspace.
//
// It manages a collection of templates and handles the rendering process.
type Builder struct {
	env       func(string) string
	flagSet   *pflag.FlagSet
	templates []Template
}

// NewBuilder creates a new template builder instance.
func NewBuilder(fs *pflag.FlagSet) *Builder {
	return &Builder{flagSet: fs}
}

type FlagConstrain interface {
	MarkFlagRequired(name string) error
}

// Init initializes all templates in the builder.
func (b *Builder) Init(fc FlagConstrain) error {
	for _, elem := range b.templates {
		if err := elem.Init(b.flagSet, fc); err != nil {
			return err
		}
	}
	return nil
}

// Build renders all templates to the specified workspace directory.
func (b *Builder) Build(workspace string, envs map[string]string) (err error) {
	if err = b.prepareWorkspace(workspace); err != nil {
		return errors.Wrap(err, "prepare workspace")
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(workspace)
		}
	}()

	b.env = func(name string) string { return envs[name] }
	for _, elem := range b.templates {
		templatePath := envSubst(elem.Path(), envs)
		if err = b.renderTo(filepath.Join(workspace, templatePath), elem); err != nil {
			return errors.Wrap(err, "render template")
		}
	}

	return
}

// prepareWorkspace ensures the workspace directory exists and is empty.
//
// It creates the directory if it doesn't exist, and verifies it's an empty directory.
func (b *Builder) prepareWorkspace(workspace string) error {
	if stat, err := os.Stat(workspace); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return os.MkdirAll(workspace, 0755)
	} else if stat != nil && !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", workspace)
	}

	dir, err := os.Open(workspace)
	if err != nil {
		return err
	}
	defer func() { _ = dir.Close() }()

	if entries, err := dir.ReadDir(1); err != nil && err != io.EOF {
		return err
	} else if len(entries) > 0 {
		return fmt.Errorf("%s is not empty", workspace)
	}

	return nil
}

// renderTo renders a single template to the specified file path.
func (b *Builder) renderTo(file string, target Template) error {
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	renderer, err := template.New(file).Funcs(b.funcMap()).Parse(target.Body())
	if err != nil {
		return errors.Wrap(err, "failed to parse template")
	}

	fp, err := os.Create(file)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer func() { _ = fp.Close() }()

	return renderer.Execute(fp, target)
}

// funcMap returns a map of custom template functions that can be used in templates.
func (b *Builder) funcMap() template.FuncMap {
	return template.FuncMap{
		"flag": flag(b.flagSet),
		"env": func(name string) string {
			if b.env == nil {
				return ""
			}
			return b.env(name)
		},
	}
}

// Template defines the interface for renderable templates.
type Template interface {
	// Init initializes the template with command-line flags.
	Init(*pflag.FlagSet, FlagConstrain) error

	// Path returns the relative file path where this template should be rendered
	Path() string

	// Body returns the template content
	Body() string
}

// AddTemplate adds one or more templates to the builder.
// These templates will be rendered when Build is called.
func (b *Builder) AddTemplate(templates ...Template) {
	b.templates = append(b.templates, templates...)
}

// envSubst performs environment variable substitution in a text string.
// It replaces placeholders in the format ${VARIABLE} with their corresponding values
// from the provided environment map.
func envSubst(text string, envs map[string]string) string {
	for k, v := range envs {
		text = strings.ReplaceAll(text, fmt.Sprintf("${%s}", strings.ToUpper(k)), v)
	}
	return text
}
