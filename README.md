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
go test -timeout 30m terraform-provider-snapcd/internal/tests/runner_pool_assignment -v -tags=all -args -test.v
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

# Generate docs

```
make generate
```

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
