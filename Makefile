PREFIX ?= $(shell pwd)

NAME = fsql
PKG = github.com/kshvmdn/$(NAME)
MAIN = $(PKG)/cmd/$(NAME)

DIST_DIR := ${PREFIX}/dist
DIST_DIRS := find . -type d | sed "s|^\./||" | grep -v \\. | tr '\n' '\0' | xargs -0 -I '{}'

SRCS := $(shell find . -type f -name '*.go')
PKGS := $(shell go list ./... | grep -v /vendor)

VERSION := $(shell cat VERSION)
GITCOMMIT := $(shell git rev-parse --short HEAD)
ifneq ($(shell git status --porcelain --untracked-files=no),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif

LDFLAGS := ${LDFLAGS} \
	-X $(PKG)/meta.GITCOMMIT=${GITCOMMIT} \
	-X $(PKG)/meta.VERSION=${VERSION}

.PHONY: build
build: $(SRCS) VERSION
	@echo "+ $@"
	@go build -ldflags "${LDFLAGS}" -o $(NAME) -v $(MAIN)

.PHONY: install
install:
	@echo "+ $@"
	@go install $(PKGS)

.PHONY: get-tools
get-tools:
	@echo "+ $@"
	@go get -u -v github.com/golang/lint/golint

.PHONY: clean
clean:
	@echo "+ $@"
	$(RM) $(NAME)
	$(RM) -r $(DIST_DIR)

.PHONY: fmt
fmt:
	@echo "+ $@"
	@test -z "$$(gofmt -s -l . 2>&1 | grep -v ^vendor/ | tee /dev/stderr)" || \
		(echo >&2 "+ please format Go code with 'gofmt -s', or use 'make fmt-save'" && false)

.PHONY: fmt-save
fmt-save:
	@echo "+ $@"
	@gofmt -s -l . 2>&1 | grep -v ^vendor/ | xargs gofmt -s -l -w

.PHONY: vet
vet:
	@echo "+ $@"
	@go vet $(PKGS)

.PHONY: lint
lint:
	@echo "+ $@"
	$(if $(shell which golint || echo ''),, \
		$(error Please install golint: `make get-tools`))
	@test -z "$$(golint ./... 2>&1 | grep -v ^vendor/ | grep -v mock/ | tee /dev/stderr)"

.PHONY: test
test:
	@echo "+ $@"
	@go test -race $(PKGS)

.PHONY: coverage
coverage:
	@echo "+ $@"
	@for pkg in $(PKGS); do \
		go test -test.short -race -coverprofile="../../../$$pkg/coverage.txt" $${pkg} || exit 1; \
	done

.PHONY: bootstrap-dist
bootstrap-dist:
	@echo "+ $@"
	@go get -u -v github.com/franciscocpg/gox

.PHONY: build-all
build-all: $(SRCS) VERSION
	@echo "+ $@"
	@gox -verbose \
		-ldflags "${LDFLAGS}" \
		-os="darwin freebsd netbsd openbsd linux windows" \
		-arch="386 amd64 arm arm64" \
		-osarch="!darwin/arm !darwin/arm64" \
		-output="$(DIST_DIR)/{{.OS}}-{{.Arch}}/{{.Dir}}" $(MAIN)

.PHONY: dist
dist: build-all
	@echo "+ $@"
	@cd $(DIST_DIR) && \
		$(DIST_DIRS) cp ../LICENSE {} && \
		$(DIST_DIRS) cp ../README.md {} && \
		$(DIST_DIRS) tar -zcf fsql-${VERSION}-{}.tar.gz {} && \
		$(DIST_DIRS) zip -r -q fsql-${VERSION}-{}.zip {} && \
		$(DIST_DIRS) rm -rf {} && \
		cd ..
