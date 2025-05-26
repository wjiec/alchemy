package v1

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/internal/template"
)

type EchoImplGo struct{}

func (e *EchoImplGo) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (e *EchoImplGo) Path() string { return "internal/service/echo/v1/echo.go" }
func (e *EchoImplGo) Body() string { return echoImplTemplate }

const echoImplTemplate = `package v1

import (
	"context"

	"{{ flag "repo" }}/api/echo"
	echov1api "{{ flag "repo" }}/api/echo/v1"
)

type EchoService struct {
	echov1api.UnimplementedEchoServiceServer
}

func (e *EchoService) Echo(_ context.Context, req *echov1api.EchoRequest) (*echov1api.EchoResponse, error) {
	if len(req.Text) == 0 {
		return nil, echo.ErrTextTooShort
	}
	return &echov1api.EchoResponse{Text: req.Text}, nil
}

func New() (*EchoService, error) {
	return &EchoService{}, nil
}`
