package echo

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/internal/template"
)

//goland:noinspection GoNameStartsWithPackageName
type EchoGo struct{}

func (e *EchoGo) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (e *EchoGo) Path() string { return "internal/service/echo/echo.go" }
func (e *EchoGo) Body() string { return echoGoTemplate }

const echoGoTemplate = `package echo

import (
	"github.com/google/wire"

	echov1 "{{ flag "repo" }}/internal/service/echo/v1"
)

var ProviderSet = wire.NewSet(
	echov1.New,
)
`
