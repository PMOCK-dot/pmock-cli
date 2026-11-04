BINARY  := pmock-cli
MODULE  := github.com/PMOCK-dot/pmock-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X $(MODULE)/internal/build.Version=$(VERSION) \
	-X $(MODULE)/internal/build.Commit=$(COMMIT) \
	-X $(MODULE)/internal/build.Date=$(DATE)
PREFIX  ?= /usr/local

.PHONY: build test vet fmt check install clean snapshot run update-spec package npm skills skills-check

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

test:
	go test -race -count=1 ./...

vet:
	go vet ./...

fmt:
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

check: fmt vet test

install: build
	install -d $(PREFIX)/bin
	install -m755 bin/$(BINARY) $(PREFIX)/bin/$(BINARY)

snapshot:
	goreleaser release --snapshot --clean

run: build
	./bin/$(BINARY) --help

update-spec:
	./scripts/update-spec.sh

# Regenerate the committed agent skills tree from the command catalog.
skills:
	go run . skills generate --out skills

# Fail when skills/ no longer matches the catalog (run in CI).
skills-check:
	@go run . skills generate --out skills --check

package:
	./scripts/package.sh

npm:
	VERSION=$${VERSION:-0.1.4} ./scripts/build-npm.sh

clean:
	rm -rf bin dist
