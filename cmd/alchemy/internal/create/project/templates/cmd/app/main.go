package app

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

type MainGo struct {
	WithGrpc    bool
	WithoutHttp bool
}

func (m *MainGo) Init(fs *pflag.FlagSet, _ template.FlagConstrain) error {
	fs.BoolVar(&m.WithGrpc, "with-grpc", false, "If set, enabling grpc service for project")
	fs.BoolVar(&m.WithoutHttp, "without-http", false, "If set, disabling http service for project")

	return nil
}

func (m *MainGo) Path() string { return "cmd/${PROJECT_NAME}/main.go" }
func (m *MainGo) Body() string { return mainGoTemplate }

const mainGoTemplate = `package main

import (
	"log/slog"

	"github.com/wjiec/alchemy"

	echov1api "{{ flag "repo" }}/api/echo/v1"
	echov1 "{{ flag "repo" }}/internal/service/echo/v1"
)

func Setup(v1Echo *echov1.EchoService) (*alchemy.App, error) {
	return alchemy.New("{{ env "PROJECT_NAME" }}",
{{- if not .WithoutHttp }}
		alchemy.WithHttpServer("tcp", ":8080"),
{{- end }}
{{- if .WithGrpc }}
		alchemy.WithGrpcServer("tcp", ":8081"),
{{- end }}
		alchemy.WithServiceRegister[echov1api.EchoServiceServer](echov1api.RegisterEchoServiceAlchemyServer, v1Echo),
	)
}

func main() {
	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	if err = app.Start(alchemy.SetupSignalHandler()); err != nil {
		slog.Error("fatal error", "error", err)
	}
}`
