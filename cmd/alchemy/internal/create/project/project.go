package project

import (
	"errors"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/create/project/templates"
	"github.com/wjiec/alchemy/cmd/alchemy/internal/create/project/templates/_internal/service"
	srvecho "github.com/wjiec/alchemy/cmd/alchemy/internal/create/project/templates/_internal/service/echo"
	srcechov1 "github.com/wjiec/alchemy/cmd/alchemy/internal/create/project/templates/_internal/service/echo/v1"
	"github.com/wjiec/alchemy/cmd/alchemy/internal/create/project/templates/api"
	echoapi "github.com/wjiec/alchemy/cmd/alchemy/internal/create/project/templates/api/echo"
	echov1api "github.com/wjiec/alchemy/cmd/alchemy/internal/create/project/templates/api/echo/v1"
	"github.com/wjiec/alchemy/cmd/alchemy/internal/create/project/templates/cmd/app"
	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <name>",
		Short: "Create a new project",
	}

	b := template.NewBuilder(cmd.Flags())
	b.AddTemplate(&templates.GoMod{}, &templates.Version{})
	b.AddTemplate(&templates.BufYaml{}, &templates.BufGenYaml{})
	b.AddTemplate(&templates.EditorConfig{}, &templates.GitIgnore{}, &templates.Dockerfile{}, &templates.Makefile{})
	b.AddTemplate(&app.MainGo{}, &app.WireGo{})
	b.AddTemplate(&api.ErrorsGo{}, &echoapi.EchoGo{}, &echov1api.EchoProto{})
	b.AddTemplate(&service.ServiceGo{}, &srvecho.EchoGo{}, &srcechov1.EchoImplGo{})
	if err := b.Init(cmd); err != nil {
		cmd.PersistentPreRunE = func(*cobra.Command, []string) error {
			return err
		}
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("project name is required")
		}

		abs, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		return b.Build(abs, map[string]string{"PROJECT_NAME": args[0]})
	}

	return cmd
}
