package app

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/internal/template"
)

type WireGo struct{}

func (w *WireGo) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (w *WireGo) Path() string { return "cmd/${PROJECT_NAME}/wire.go" }
func (w *WireGo) Body() string { return wireGoTemplate }

const wireGoTemplate = `//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/wjiec/alchemy"

	"{{ flag "repo" }}/internal/service"
)

func NewApp() (*alchemy.App, error) {
	panic(wire.Build(
		Setup,
		service.ProviderSet,
	))
}
`
