ALCHEMY_VERSION := $(shell cat VERSION)
ALCHEMY_PACKAGE := $(shell go list -m)

.PHONY: default help build install warm-up
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

.PHONY: generate
generate: ## Generate code from protocol buffer definitions.
	buf build && buf generate


##@ Build

.PHONY: alchemy
alchemy:
	go build \
		-ldflags="-X '$(ALCHEMY_PACKAGE)/cmd/alchemy/internal/version.Version=$(ALCHEMY_VERSION)'" \
		-o bin/alchemy ./cmd/alchemy

.PHONY: protoc-gen-alchemy
protoc-gen-alchemy:
	cd cmd/protoc-gen-alchemy && go build \
		-ldflags="-X '$(ALCHEMY_PACKAGE)/cmd/protoc-gen-alchemy/internal/gengo.Version=$(ALCHEMY_VERSION)'" \
		-o ../../bin/protoc-gen-alchemy .
