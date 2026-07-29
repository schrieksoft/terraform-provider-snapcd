#!/usr/bin/env bash
set -euo pipefail

# Downloads the SnapCD OpenAPI document from a snapcd GitHub release and drops it at
# schemas/openapi.yaml — the vendored source of truth for attribute descriptions and
# "Required permissions" blocks (see tools/openapigen). `make sync-local` is the
# local-development shortcut that copies from a checkout instead.
#
# The version comes from versions.env, which Renovate bumps when a new snapcd release
# appears. After syncing, run `make generate` to refresh the generated code and docs.
#
# Usage:
#   scripts/fetch-openapi-spec.sh              # version from versions.env
#   scripts/fetch-openapi-spec.sh 1.9.0        # explicit version

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

VERSION="${1:-}"
if [[ -z "$VERSION" ]]; then
    VERSION="$(grep -oP '^SNAPCD_VERSION=\K.*' versions.env)"
fi

BASE="https://github.com/schrieksoft/snapcd/releases/download/${VERSION}"

echo "Fetching snapcd ${VERSION} OpenAPI document"

# Download to a staging file first so a failed download cannot truncate the vendored copy.
STAGING="$(mktemp)"
trap 'rm -f "$STAGING"' EXIT

if ! curl -fsSL "${BASE}/openapi.yaml" -o "$STAGING"; then
    echo "  failed to download openapi.yaml from ${VERSION}" >&2
    echo "  (releases before the artifact-publishing change do not carry it)" >&2
    exit 1
fi

mv "$STAGING" schemas/openapi.yaml
trap - EXIT
echo "  schemas/openapi.yaml"
