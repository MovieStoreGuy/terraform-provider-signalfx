TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=signalfx

SRC_ROOT        := $(shell git rev-parse --show-toplevel)
SRC_GO_FILES    := $(shell find $(SRC_ROOT) -name '*.go')
TOOLS_MOD_DIR   := $(SRC_ROOT)/internal/tools
TOOLS_BIN_DIR   := $(SRC_ROOT)/.tools
TOOLS_MOD_REGEX := "\s+_\s+\".*\""
TOOLS_PKG_NAMES := $(shell grep -E $(TOOLS_MOD_REGEX) < $(TOOLS_MOD_DIR)/tools.go | tr -d " _\"" | grep -vE '/v[0-9]+$$')
TOOLS_BIN_NAMES := $(addprefix $(TOOLS_BIN_DIR)/, $(notdir $(shell echo $(TOOLS_PKG_NAMES))))

ADDLICENCESE   := $(TOOLS_BIN_DIR)/addlicense
GOVULNCHECK    := $(TOOLS_BIN_DIR)/govulncheck
GOLANGCI_LINT  := $(TOOLS_BIN_DIR)/golangci-lint
WEBSITE_PLUGIN := $(TOOLS_BIN_DIR)/tfplugindocs

default: build

.PHONY: install-tools
install-tools: $(TOOLS_BIN_NAMES)

$(TOOLS_BIN_DIR):
	mkdir -p $@

$(TOOLS_BIN_NAMES): $(TOOLS_BIN_DIR) $(TOOLS_MOD_DIR)/go.mod
	go -C $(TOOLS_MOD_DIR) build -o $@ -trimpath $(filter %/$(notdir $@),$(TOOLS_PKG_NAMES))

.PHONY: addlicense
addlicense: $(ADDLICENCESE)
	@ADDLICENCESEOUT=`$(ADDLICENCESE) -y "" -c 'Splunk, Inc.' -l mpl -s=only $(SRC_GO_FILES) 2>&1`; \
		if [ "$$ADDLICENCESEOUT" ]; then \
			echo "$(ADDLICENCESE) FAILED => add License errors:\n"; \
			echo "$$ADDLICENCESEOUT\n"; \
			exit 1; \
		else \
			echo "Add License finished successfully"; \
		fi

.PHONY: checklicense
checklicense: $(ADDLICENCESE)
	@ADDLICENCESEOUT=`$(ADDLICENCESE) -check $(SRC_GO_FILES) 2>&1`; \
		if [ "$$ADDLICENCESEOUT" ]; then \
			echo "$(ADDLICENCESE) FAILED => add License errors:\n"; \
			echo "$$ADDLICENCESEOUT\n"; \
			echo "Use 'make addlicense' to fix this."; \
			exit 1; \
		else \
			echo "Check License finished successfully"; \
		fi

.PHONY: govulncheck
govulncheck: $(GOVULNCHECK)
	$(GOVULNCHECK)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run -v

.PHONY: lint-fix
lint-fix:
	$(GOLANGCI_LINT) run -v --fix

build:
	go build

test:
	go test --cover --race -v --timeout 30s ./...

test-with-cover:
	mkdir -p $(PWD)/coverage/unit || true
	go test --race --timeout 300s --cover ./... \
		-covermode=atomic \
		-args -test.gocoverdir="$(PWD)/coverage/unit"
	go tool covdata textfmt -i=./coverage/unit -o ./coverage.txt


testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

fmt:
	gofmt -w $(GOFMT_FILES)


test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

check-docs: gen-docs
	@if [ "`git status --porcelain docs/`" ];then \
		git diff;\
		echo "Changes to documentation are not committed. Please run 'make gen-docs' and commit the changes" && \
		echo `git status --porcelain docs/` &&\
		exit 1;\
	fi


gen-docs: $(WEBSITE_PLUGIN)
	$(WEBSITE_PLUGIN)

test-docs: $(WEBSITE_PLUGIN)
	$(WEBSITE_PLUGIN) validate 

.PHONY: build test testacc vet fmt fmtcheck errcheck gen-docs check-docs
