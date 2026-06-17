#!/usr/bin/env bash
# Resolve a dogfood corpus repo name to the "--repo owner/repo" flag for actiondoc, taking the
# identity from the canonical GitHub URL pinned in dogfood/manifest.txt. This gives the output
# snapshots the same repository identity a real user has (a git remote or $GITHUB_REPOSITORY),
# so a repo's references to its own reusable workflows resolve as self-calls rather than
# external. Prints nothing when the name is not in the manifest (caller then omits --repo).
set -euo pipefail

name="${1:?usage: repo_identity.sh <corpus-repo-name>}"
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
manifest="$root/dogfood/manifest.txt"

url="$(grep -v '^#' "$manifest" | awk -F'\t' -v n="$name" '$1 == n { print $2; exit }')"
[ -n "$url" ] || exit 0

repo="$(printf '%s' "$url" | sed -E 's,^https?://github.com/,,; s,/+$,,')"
[ -n "$repo" ] && printf -- '--repo %s' "$repo"
