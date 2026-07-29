# Contribute

- Install "pre-commit", e.g. `pip3 install pre-commit==2.20.0`. You'll need [Python 3 ](https://www.python.org/downloads/release/python-3110/)
- Run `pre-commit install` in root of folder to activate the pre-commit hooks.

# Test

```shell
go test -timeout 30m terraform-provider-snapcd/internal/tests/core -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/extra_files -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/identity -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/module_env_var -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/module_param -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/namespace_env_var -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/namespace_param -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/role_assignment -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/runner_assignment -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/secret -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/secret_store -v -tags=all -args -test.v
go test -timeout 30m terraform-provider-snapcd/internal/tests/secret_store_assignment -v -tags=all -args -test.v

```


go test -timeout 30m terraform-provider-snapcd/internal/tests/* -v -tags=all -args -test.v

# Build

```shell
$ go build -o terraform-provider-snapcd
```

# Install locally

```shell
export BINARY=terraform-provider-snapcd
export VERSION=1.0.0
export OS_ARCH=linux_amd64
go build -o ${BINARY}
mkdir -p ~/.terraform.d/plugins/schriek/dev/snapcd/${VERSION}/${OS_ARCH}
mv ${BINARY} ~/.terraform.d/plugins/schriek/dev/snapcd/${VERSION}/${OS_ARCH}
```

# Docs

Everything under `docs/` is generated — never edit it by hand. There are two inputs:

1. **The vendored OpenAPI document** (`schemas/openapi.yaml`), owned by the snapcd code repo. It is the source of truth for attribute descriptions, the per-resource "Required permissions" sections, and enum value lists. `tools/openapigen` turns it into `internal/provider/openapidocs/` (description constants referenced by the resource schemas — with "Must be one of …" sentences generated from the enum schemas — permission blocks appended to each schema's `MarkdownDescription`, and per-enum values slices used in `stringvalidator.OneOf`). A handful of validators are intentionally narrower than their spec enum and stay hand-written; each carries a comment saying why.
2. **The Go schemas, templates and examples**, which `tfplugindocs` renders into `docs/`.

## Updating the docs

After changing resource schemas, templates or examples:

```shell
make generate
```

This runs `openapigen` and then `tfplugindocs`, rewriting `internal/provider/openapidocs/` and `docs/`. Commit the result — pre-commit (and CI) re-runs the sync + generate chain and fails on any drift.

## Updating the OpenAPI document

When descriptions or permissions change upstream, refresh the vendored spec first, then regenerate:

```shell
make sync          # download the tagged release pinned in versions.env (or VERSION=1.9.0)
make sync-local    # copy from a local snapcd checkout (SNAPCD_REPO=..., defaults to the monorepo layout)
make generate
```

Pre-commit runs `make sync generate` before every commit, and the CI test workflow runs the same chain and fails on any resulting diff — the vendored spec, the generated code and `docs/` are therefore always pinned to the tagged snapcd release in `versions.env`. `make sync-local` is a preview tool only: committing after it will re-sync the tagged version, so unreleased spec changes land here by shipping a snapcd release and bumping `versions.env` (Renovate does this automatically).

`versions.env` is bumped by Renovate when a new snapcd release appears; the bump PR needs a `make sync && make generate` commit before it can pass CI.

If `make generate` fails to compile because a referenced `openapidocs` constant disappeared, the spec no longer documents that property — fix it in the snapcd repo (the schema-doc coverage gate there should have caught it) rather than hand-writing a description here.

# Provider

```hcl
terraform {
  required_providers {
    snapcd = {
      source  = "schriek/dev/snapcd"
      version = "1.0.0"
    }
  }
}
```
