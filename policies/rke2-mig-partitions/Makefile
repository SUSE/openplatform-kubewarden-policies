ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_DIR := $(abspath $(ROOT_DIR)/bin)

SOURCE_FILES := $(shell find . -type f -name '*.go')

GOLANGCI_LINT_VER := v2.1.6
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT := $(BIN_DIR)/$(GOLANGCI_LINT_BIN)

policy.wasm: $(SOURCE_FILES) go.mod go.sum
	GOOS=wasip1 GOARCH=wasm go build -gcflags=all="-B" -ldflags="-w -s" -o policy.wasm ./cmd
	wasm-opt --enable-bulk-memory -Oz -o policy.wasm policy.wasm

annotated-policy.wasm: policy.wasm metadata.yml
	kwctl annotate -m metadata.yml -u README.md -o annotated-policy.wasm policy.wasm

.PHONY: test
test:
	go test ./internal/...

.PHONY: e2e-tests
e2e-tests: annotated-policy.wasm
	bats e2e.bats

golangci-lint: $(GOLANGCI_LINT) ## Install a local copy of golang ci-lint.
$(GOLANGCI_LINT): ## Install golangci-lint.
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VER)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	go vet ./...
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fix

.PHONY: clean
clean:
	go clean
	rm -f policy.wasm annotated-policy.wasm

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: push_gitlab
push_gitlab:
	kwctl push --docker-config-json-path=`pwd` ./annotated-policy.wasm $(CI_REGISTRY)/itpe/core/open-platform/kubewarden-policies/policy-pod-mig-partitions:$(VERSION)

.PHONY: push_harbor
push_harbor:
	kwctl push --docker-config-json-path=`pwd` ./annotated-policy.wasm harbor.op-prg2-0-ingress.op.suse.org/policy-pod-mig-partitions/policy-pod-mig-partitions:$(VERSION)
