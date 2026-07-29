default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

# The vendored OpenAPI spec (schemas/openapi.yaml) is the source of truth for attribute
# descriptions and "Required permissions" blocks. Two ways to refresh it (run
# `make generate` afterwards so the generated code and docs pick up the new document):
#
#   make sync        — download from a snapcd GitHub release (the version comes from
#                      versions.env, or pass VERSION=1.9.0)
#   make sync-local  — copy from a local snapcd checkout, for previewing changes that
#                      have not been released yet
SNAPCD_REPO ?= ../../../applications/snapcd

sync:
	@scripts/fetch-openapi-spec.sh $(VERSION)
	@git --no-pager diff --stat -- schemas || true

sync-local:
	@test -d "$(SNAPCD_REPO)/schemas" || { \
	    echo "snapcd repo not found at $(SNAPCD_REPO)"; \
	    echo "pass SNAPCD_REPO=/path/to/snapcd, e.g. make sync-local SNAPCD_REPO=~/code/snapcd"; \
	    exit 1; \
	}
	cp "$(SNAPCD_REPO)/schemas/openapi.yaml" schemas/openapi.yaml
	@git --no-pager diff --stat -- schemas || true

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate sync sync-local
