#!/usr/bin/env bash
set -euo pipefail

# bump-version.sh updates every submodule go.mod in this repo so that
# cross-module requires reference the newly released version, and drops
# the internal `replace` directives that only exist for local development.
#
# Usage:
#   scripts/bump-version.sh <version>
#
# Example:
#   scripts/bump-version.sh v0.0.12

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
  echo "Error: version argument required (e.g. v0.0.12)" >&2
  exit 1
fi

if [[ "$VERSION" != v* ]]; then
  echo "Error: version must start with 'v' (e.g. v0.0.12)" >&2
  exit 1
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

# Collect every go.mod in the repo (excluding vendor dirs).
MODFILES=()
while IFS= read -r mf; do
  MODFILES+=("$mf")
done < <(find . -name go.mod -not -path '*/vendor/*' | sort)

# Build the set of internal module paths (everything under our module root)
# that can be a cross-module dependency. The repo-root module itself is
# excluded so we never try to require the top-level module from itself.
MODULE_PREFIX="github.com/davidcollom/terratest-utils"
INTERNAL=()
for mf in "${MODFILES[@]}"; do
  path="$(awk '$1=="module" {print $2; exit}' "$mf")"
  if [[ "$path" == "$MODULE_PREFIX/"* ]]; then
    INTERNAL+=("$path")
  fi
done

if [[ "${#INTERNAL[@]}" -eq 0 ]]; then
  echo "No internal modules found (nothing to update)."
  exit 0
fi

echo "Updating all modules to reference version $VERSION"
echo

changed=0
for mf in "${MODFILES[@]}"; do
  dir="$(dirname "$mf")"
  own="$(awk '$1=="module" {print $2; exit}' "$mf")"
  for pkg in "${INTERNAL[@]}"; do
    # Skip the module's own path, and modules that don't reference this dep.
    if [[ "$pkg" == "$own" ]]; then
      continue
    fi
    if ! grep -q "$pkg" "$mf" 2>/dev/null; then
      continue
    fi
    echo "-- $mf"
    echo "   require $pkg $VERSION"
    (
      cd "$dir"
      go mod edit -require="$pkg@$VERSION"
      # Local dev `replace` directives must not ship in a release.
      go mod edit -dropreplace="$pkg"
    )
    changed=1
  done
done

if [[ "$changed" -eq 0 ]]; then
  echo "No submodules reference internal dependencies; nothing changed."
fi

echo
echo "Done. Cross-module requires now reference $VERSION."
echo "Tip: run 'go mod tidy' in each changed module if you need go.sum updates."
