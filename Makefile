BIN = $(CURDIR)/bin
GO = go
M = $(shell printf "\033[34;1m▶\033[0m")
EXECUTABLE = go-hello

EXTLDFLAGS = -extldflags "-static"
TAGS = netgo osusergo static_build
GOBUILD = $(GO) build -a -trimpath -ldflags '-s -w $(EXTLDFLAGS)' -tags '$(TAGS)'

BUILD = CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) $(GOBUILD) -o $(BIN)/$(EXECUTABLE)-$(1)-$(2)

.DEFAULT_GOAL := help

.PHONY: build
build: lint | $(BIN) ; $(info $(M) building executable…) @ ## Build program binary
	@$(call BUILD,linux,amd64)

.PHONY: docker-build
docker-build:
	@docker build --force-rm --compres-t $(shell basename $(CURDIR)) .

.PHONY: docker-run
docker-run:
	@docker run --rm -it -p 8080:80 $(shell basename $(CURDIR))

# Tools

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(PACKAGE)…)
	@tmp=$$(mktemp -d); \
	env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(PACKAGE) || ret=$$?; \
	rm -rf $$tmp ; exit $$ret

GOLINT = $(BIN)/golint
$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

STATICCHECK = $(BIN)/staticcheck
$(BIN)/staticcheck: PACKAGE=honnef.co/go/tools/cmd/staticcheck

.PHONY: lint
golint: | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	@$(GOLINT) -set_exit_status ./...

.PHONY: staticcheck
staticcheck: | $(STATICCHECK) ; $(info $(M) running staticcheck…) @ ## Run staticcheck
	@$(STATICCHECK) ./...

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@$(GO) fmt ./...

.PHONY: vet
vet: ; $(info $(M) running govet…) @ ## Run govet on all source files
	@$(GO) vet ./...

.PHONY: lint
lint: fmt golint vet staticcheck

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@$(GO) clean -i ./...
	@rm -rf $(BIN)

.PHONY: help
help: ## Provides help information on available commands
	@awk -F ':.*?## ' '/^[a-zA-Z0-9[:space:]_.-]+:.*?##/ {printf "\033[36m%-18s\033[0m %s\n", $$1, $$NF}' $(MAKEFILE_LIST)
