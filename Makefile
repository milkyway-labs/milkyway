#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)

export GO111MODULE = on

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

TM_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')

###############################################################################
###                               Build flags                               ###
###############################################################################

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
	ifeq ($(OS),Windows_NT)
		GCCEXE = $(shell where gcc.exe 2> NUL)
		ifeq ($(GCCEXE),)
			$(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
		else
			build_tags += ledger
		endif
	else
		UNAME_S = $(shell uname -s)
		ifeq ($(UNAME_S),OpenBSD)
			$(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
		else
			GCC = $(shell command -v gcc 2> /dev/null)
			ifeq ($(GCC),)
				$(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
			else
				build_tags += ledger
			endif
		endif
	endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += rocksdb
endif
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += boltdb
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

###############################################################################
###                               Linker flags                              ###
###############################################################################

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=milkyway \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=milkywayd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TM_VERSION)

# Static linking
ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb
endif

# Handle RocksDB
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  $(info ################################################################)
  $(info To use rocksdb, you need to install rocksdb first)
  $(info Please follow this guide https://github.com/rockset/rocksdb-cloud/blob/master/INSTALL.md)
  $(info ################################################################)
  CGO_ENABLED=1
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
endif

# Handle BoltDB
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

###############################################################################
###                                   All                                   ###
###############################################################################

all: tools install lint test

###############################################################################
###                                 Build                                   ###
###############################################################################

build: go.sum
ifeq ($(OS),Windows_NT)
	exit 1
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/milkywayd ./cmd/milkywayd
endif

create-builder: go.sum
	$(MAKE) -C contrib/images milkyway-builder CONTEXT=$(CURDIR)

build-alpine: create-builder
	mkdir -p $(BUILDDIR)
	$(DOCKER) build -f Dockerfile --rm --tag milkywaylabs/milkyway-alpine .
	$(DOCKER) create --name milkyway-alpine --rm milkywaylabs/milkyway-alpine
	$(DOCKER) cp milkyway-alpine:/usr/bin/milkywayd $(BUILDDIR)/milkywayd
	$(DOCKER) rm milkyway-alpine

build-linux: create-builder
	mkdir -p $(BUILDDIR)
	$(DOCKER) build -f Dockerfile-ubuntu --rm --tag milkywaylabs/milkyway-linux .
	$(DOCKER) create --name milkyway-linux milkywaylabs/milkyway-linux
	$(DOCKER) cp milkyway-linux:/usr/bin/milkywayd $(BUILDDIR)/milkywayd
	$(DOCKER) rm milkyway-linux

build-reproducible: go.sum
	$(DOCKER) rm latest-build || true
	$(DOCKER) run --volume=$(CURDIR):/sources:ro \
        --env TARGET_PLATFORMS='linux/amd64 linux/arm64 darwin/amd64 windows/amd64' \
        --env APP=milkyway \
        --env VERSION=$(VERSION) \
        --env COMMIT=$(COMMIT) \
        --env LEDGER_ENABLED=$(LEDGER_ENABLED) \
        --name latest-build cosmossdk/rbuilder:latest
	$(DOCKER) cp -a latest-build:/home/builder/artifacts/ $(CURDIR)/

install: go.sum 
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/milkywayd

update-swagger-docs: statik
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
        echo "\033[92mSwagger docs are in sync\033[0m";\
    fi

.PHONY: build build-linux install update-swagger-docs

###############################################################################
###                                Protobuf                                 ###
###############################################################################

bufVer=1.47.2
bufImageName=bufbuild/buf:$(bufVer)
bufImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(bufImageName)

protoVer=0.15.2
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Generating Swagger files"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh
	$(MAKE) update-swagger-docs

proto-pulsar-gen:
	@echo "Generating Dep-Inj Protobuf files"
	@$(protoImage) sh ./scripts/protocgen-pulsar.sh

proto-format:
	@$(bufImage) format -w

proto-lint:
	@$(bufImage) lint --error-format=json ./proto

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

.PHONY: proto-all proto-gen proto-swagger-gen proto-pulsar-gen proto-format proto-lint proto-check-breaking

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	@go install github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/milkywayd -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf \
    $(BUILDDIR)/ \
    artifacts/ \
    tmp-swagger-gen/

###############################################################################
###                                   Mocks                                 ###
###############################################################################
mockgen:
	@./scripts/mockgen.sh

.PHONY: mocks

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

include sims.mk

test: test-unit
test-all: test-unit test-ledger-mock test-race test-cover

PACKAGES_UNIT=$(shell go list ./... | grep -v -e '/tests/e2e')
PACKAGES_E2E=$(shell cd tests/e2e && go list ./... | grep '/e2e')
TEST_PACKAGES=./...
TEST_TARGETS := test-unit test-unit-cover test-race test-e2e

# Test runs-specific rules. To add a new test target, just add
# a new rule, customise ARGS or TEST_PACKAGES ad libitum, and
# append the new rule to the TEST_TARGETS list.
test-unit: ARGS=-timeout=5m -tags='norace'
test-unit: TEST_PACKAGES=$(PACKAGES_UNIT)
test-unit-cover: ARGS=-timeout=5m -tags='norace' -coverprofile=coverage.txt -covermode=atomic
test-unit-cover: TEST_PACKAGES=$(PACKAGES_UNIT)
test-ledger: test_tags += cgo ledger norace
test-ledger-mock: test_tags += ledger test_ledger_mock norace
test-race: ARGS=-timeout=5m -race
test-race: TEST_PACKAGES=$(PACKAGES_UNIT)
test-e2e: ARGS=-timeout=35m -v
test-e2e: TEST_PACKAGES=$(PACKAGES_E2E)
$(TEST_TARGETS): run-tests

ARGS += -tags "$(test_tags)"
SUB_MODULES = $(shell find . -type f -name 'go.mod' -print0 | xargs -0 -n1 dirname | sort)
CURRENT_DIR = $(shell pwd)

run-tests:
ifneq (,$(shell which tparse 2>/dev/null))
	@echo "--> Running tests"
	@go test -mod=readonly -json $(ARGS) $(TEST_PACKAGES) | tparse
else
	@echo "--> Running tests"
	@go test -mod=readonly $(ARGS) $(TEST_PACKAGES)
endif

.PHONY: run-tests $(TEST_TARGETS)

###############################################################################
###                                Benchmark                                ###
###############################################################################

benchstat_cmd=golang.org/x/perf/cmd/benchstat
benchmark:
	@go test -mod=readonly -bench=. -count=$(BENCH_COUNT) -run=^a  ./... > bench-$(REF_NAME).txt
	@test -e bench-master.txt && go run $(benchstat_cmd) bench-master.txt bench-$(REF_NAME).txt || go run $(benchstat_cmd) bench-$(REF_NAME).txt
.PHONY: benchmark

###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	@echo "--> Running linter"
	@go run $(golangci_lint_cmd) run --timeout=10m

lint-fix:
	@echo "--> Running linter"
	@go run $(golangci_lint_cmd) run --fix --out-format=tab --issues-exit-code=0

.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name "*_mocks.go" -not -name '*.pb.go' -not -name '*.pulsar.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name "*_mocks.go" -not -name '*.pb.go' -not -name '*.pulsar.go' | xargs misspell -w
	find . -name '*.go' -type f -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name "*_mocks.go" -not -name '*.pb.go' -not -name '*.pulsar.go' | xargs goimports -w -local github.com/milkyway-labs/milkyway
	find . -name '*.proto' -type f -not -path "*.git*" | xargs misspell -w
.PHONY: format

###############################################################################
###                                Localnet                                 ###
###############################################################################

start-localnet-ci: build
	rm -rf ~/.milkywayd-liveness
	./build/milkywayd init liveness --chain-id liveness --home ~/.milkywayd-liveness
	./build/milkywayd config set client chain-id liveness --home ~/.milkywayd-liveness
	./build/milkywayd config set client keyring-backend test --home ~/.milkywayd-liveness
	./build/milkywayd keys add val --home ~/.milkywayd-liveness --keyring-backend test
	./build/milkywayd genesis add-genesis-account val 10000000000000000000000000stake --home ~/.milkywayd-liveness --keyring-backend test
	./build/milkywayd genesis gentx val 1000000000stake --home ~/.milkywayd-liveness --chain-id liveness --keyring-backend test
	./build/milkywayd genesis collect-gentxs --home ~/.milkywayd-liveness
	sed -i.bak'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ~/.milkywayd-liveness/config/app.toml
	./build/milkywayd start --home ~/.milkywayd-liveness --x-crisis-skip-assert-invariants

.PHONY: start-localnet-ci