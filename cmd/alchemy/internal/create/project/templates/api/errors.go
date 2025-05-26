package api

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/internal/template"
)

type ErrorsGo struct{}

func (e *ErrorsGo) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (e *ErrorsGo) Path() string { return "api/errors.go" }
func (e *ErrorsGo) Body() string { return errorsGoTemplate }

const errorsGoTemplate = `package api

import "github.com/wjiec/alchemy/bizerr"

var (
	ErrEchoServiceZone = bizerr.ThousandStep.Next()
)`
