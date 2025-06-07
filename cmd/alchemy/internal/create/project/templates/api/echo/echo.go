package echo

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

//goland:noinspection GoNameStartsWithPackageName
type EchoGo struct{}

func (e *EchoGo) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (e *EchoGo) Path() string { return "api/echo/echo.go" }
func (e *EchoGo) Body() string { return echoGoTemplate }

const echoGoTemplate = `package echo

import (
	"net/http"

	"{{ flag "repo" }}/api"
)

var (
	// ErrTextTooShort represents the too short text provided
	ErrTextTooShort = api.ErrEchoServiceZone.New(http.StatusBadRequest, "text too short")
)
`
