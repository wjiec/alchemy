package service

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

//goland:noinspection GoNameStartsWithPackageName
type ServiceGo struct{}

func (s *ServiceGo) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (s *ServiceGo) Path() string { return "internal/service/service.go" }
func (s *ServiceGo) Body() string { return serviceGoTemplate }

const serviceGoTemplate = `package service

import (
	"github.com/google/wire"

	"{{ flag "repo" }}/internal/service/echo"
)

var ProviderSet = wire.NewSet(
	echo.ProviderSet,
)
`
