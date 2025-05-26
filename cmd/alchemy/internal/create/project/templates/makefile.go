package templates

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/internal/template"
)

type Makefile struct{}

func (m *Makefile) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (m *Makefile) Path() string { return "Makefile" }
func (m *Makefile) Body() string { return makefileTemplate }

const makefileTemplate = `APP_VERSION := $(shell cat VERSION)

.PHONY: default help
.DEFAULT_GOAL: help

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Development

.PHONY: prepare
prepare: ## Generate protocol buffer dependencies and export them
	- buf build
	if buf dep update &>/dev/null; then buf export --output third_party/proto; fi

.PHONY: generate
generate: ## Generate code from protocol buffer definitions
	buf generate --exclude-path third_party


##@ Build

.PHONY: build
build:
	go build -o bin/app ./cmd/...
`
