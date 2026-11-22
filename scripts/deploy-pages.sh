#!/usr/bin/env bash
set -euo pipefail

# Build the static site + unified MyGo WASM playground and publish to the orphan gh-pages branch.
TMP_DIR="$(mktemp -d /tmp/mygox-ui-ghpages-XXXXXX)"
trap 'rm -rf "$TMP_DIR"' EXIT

bash ./scripts/build-site.sh "$TMP_DIR"

REMOTE_URL="$(git remote get-url origin)"
COMMIT_SHA="$(git rev-parse --short HEAD)"
USER_NAME="$(git config user.name || echo margbug01)"
USER_EMAIL="$(git config user.email || echo margbug01@users.noreply.github.com)"

(
  cd "$TMP_DIR"
  git init -b gh-pages >/dev/null
  git config user.name "$USER_NAME"
  git config user.email "$USER_EMAIL"
  git add -A
  git commit -q -m "deploy: mygox-ui site & wasm playground from $COMMIT_SHA"
  git push -f "$REMOTE_URL" gh-pages
)

echo "Published gh-pages from $COMMIT_SHA"
