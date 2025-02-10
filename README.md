# Test

```shell
go test -timeout 30m terraform-provider-snapcd/internal/tests -v -tags=all -args -test.v
```

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
