
NAME?=aim


PACKAGES=$(shell go list ./...)
BUILDDIR?=$(CURDIR)/dist
OUTPUT?=$(BUILDDIR)/$(NAME)

BUILD_TAGS?=$(NAME)


COMMIT_HASH := $(shell git rev-parse --short HEAD)
LD_FLAGS = -X github.com/BlockPILabs/aa-scan/version.GitCommitHash="$(COMMIT_HASH)"
BUILD_FLAGS = -mod=readonly -ldflags "$(LD_FLAGS)"

CGO_ENABLED ?= 0

# handle nostrip
ifeq (,$(findstring nostrip,$(COMETBFT_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
  LD_FLAGS += -s -w
endif

# handle race
ifeq (race,$(findstring race,$(COMETBFT_BUILD_OPTIONS)))
  CGO_ENABLED=1
  BUILD_FLAGS += -race
endif

# handle cleveldb
ifeq (cleveldb,$(findstring cleveldb,$(COMETBFT_BUILD_OPTIONS)))
  CGO_ENABLED=1
  BUILD_TAGS += cleveldb
endif

# handle badgerdb
ifeq (badgerdb,$(findstring badgerdb,$(COMETBFT_BUILD_OPTIONS)))
  BUILD_TAGS += badgerdb
endif

# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COMETBFT_BUILD_OPTIONS)))
  CGO_ENABLED=1
  BUILD_TAGS += rocksdb
endif

# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COMETBFT_BUILD_OPTIONS)))
  BUILD_TAGS += boltdb
endif

# allow users to pass additional flags via the conventional LDFLAGS variable
LD_FLAGS += $(LDFLAGS)

# Process Docker environment varible TARGETPLATFORM
# in order to build binary with correspondent ARCH
# by default will always build for linux/amd64
TARGETPLATFORM ?=
GOOS ?= linux
GOARCH ?= amd64
GOARM ?=

ifeq (linux/arm,$(findstring linux/arm,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=arm
	GOARM=7
endif

ifeq (linux/arm/v6,$(findstring linux/arm/v6,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=arm
	GOARM=6
endif

ifeq (linux/arm64,$(findstring linux/arm64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=arm64
	GOARM=7
endif

ifeq (linux/386,$(findstring linux/386,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=386
endif

ifeq (linux/amd64,$(findstring linux/amd64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=amd64
endif

ifeq (linux/mips,$(findstring linux/mips,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mips
endif

ifeq (linux/mipsle,$(findstring linux/mipsle,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mipsle
endif

ifeq (linux/mips64,$(findstring linux/mips64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mips64
endif

ifeq (linux/mips64le,$(findstring linux/mips64le,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mips64le
endif

ifeq (linux/riscv64,$(findstring linux/riscv64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=riscv64
endif





all: build

generate: clean
	go generate ./internal/entity


build: clean generate
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o $(OUTPUT) ./cmd/aim/
.PHONY: build

install: clean generate
	CGO_ENABLED=$(CGO_ENABLED) go install $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o $(OUTPUT) ./cmd/aim/

.PHONY: install

migrate: clean generate
	CGO_ENABLED=$(CGO_ENABLED) go run  ./cmd/aim/ --home ./.aim migrate

.PHONY: migrate

clean:
	rm -rf dist/*
	rm -rf internal/entity/ent
.PHONY: clean

