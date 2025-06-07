package templates

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

type Dockerfile struct{}

func (d *Dockerfile) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (d *Dockerfile) Path() string { return "Dockerfile" }
func (d *Dockerfile) Body() string { return dockerfileTemplate }

const dockerfileTemplate = `FROM golang:alpine AS builder

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

COPY api api
COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 go build -o app ./cmd/...


FROM alpine

COPY --from=builder /workspace/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]
`
