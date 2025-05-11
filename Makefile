VERSION_BRANCH ?= release
VERSION_MAJOR ?= 1
VERSION_MINOR ?= 0
VERSION_PATCH ?= 0

RELEASE_FILE = demo-$(VERSION_BRANCH).bin.v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_PATCH).tar.gz

GIT_URL ?= github.com/suzi1037/pcmd
GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_DATE ?= $(shell date '+%Y-%m-%d-%H:%M:%S')

PKG = github.com/suzi1037/pcmd

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

CODE_GEN = app/code-gen
CODE_GEN_MAIN = $(CODE_GEN)/main.go

DEMO_ROOT = app/demo
DEMO_CMD_PKG = $(PKG)/$(DEMO_ROOT)/cmd
DEMO_MAIN = $(DEMO_ROOT)/main.go

DEMO_LDFLAGS = -s -w \
	-X $(DEMO_CMD_PKG).gitBranch=$(VERSION_BRANCH) \
	-X $(DEMO_CMD_PKG).verMajor=$(VERSION_MAJOR) \
	-X $(DEMO_CMD_PKG).verMinor=$(VERSION_MINOR) \
	-X $(DEMO_CMD_PKG).verPatch=$(VERSION_PATCH) \
	-X $(DEMO_CMD_PKG).gitURL=$(GIT_URL) \
	-X $(DEMO_CMD_PKG).gitCommit=$(GIT_COMMIT) \
	-X $(DEMO_CMD_PKG).buildDate=$(BUILD_DATE)

SERVER_BIN = demo
DEMO_BUILD = CGO_ENABLED=0 go build -ldflags "$(DEMO_LDFLAGS)" -o bin/$(SERVER_BIN) $(DEMO_MAIN)

-include test/private/Makefile.mk

.PHONY: all
all: fmt vet build

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build: gen
	$(DEMO_BUILD)

.PHONY: self-gen
self-gen:
	go run $(CODE_GEN_MAIN)

.PHONY: generate gen
generate: self-gen controller-gen
	$(CONTROLLER_GEN) object paths=./$(DEMO_ROOT)/apis/v1/...
	$(CONTROLLER_GEN) object paths=./$(DEMO_ROOT)/apis/v1beta1/...
	#$(CONTROLLER_GEN) object paths=./$(DEMO_ROOT)/.../...

gen: generate

.PHONY: local-release
local-release:
	goreleaser release --snapshot --clean

.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT) run

.PHONY: dev
dev:
	go run $(DEMO_MAIN) $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

%:
	@echo ""

.PHONY: clean
clean:
	@rm -fr bin *.log

##@ Test

.PHONY: test-gen
test-gen: self-gen controller-gen
	$(CONTROLLER_GEN) object paths=./test/pcmd/.../...

##@ Dependencies

CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
CONTROLLER_TOOLS_VERSION ?= v0.16.4

GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.59.1

UPX ?= $(LOCALBIN)/upx

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN)
$(CONTROLLER_GEN): $(LOCALBIN)
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

.PHONY: upx
upx: $(UPX)
$(UPX): $(LOCALBIN)
	@command -v upx > /dev/null 2>&1 || yum install -y upx

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef