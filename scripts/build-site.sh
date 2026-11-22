#!/usr/bin/env bash
set -euo pipefail

OUT="${1:-dist}"
rm -rf "$OUT"
mkdir -p "$OUT/playground" "$OUT/screenshots"

# Copy static landing page, playground shell & screenshots
cp site/index.html "$OUT/index.html"
cp site/playground/index.html "$OUT/playground/index.html"
cp docs/screenshots/*.png "$OUT/screenshots/"

# Copy Go's matching wasm_exec.js runtime
GOROOT_DIR="$(go env GOROOT)"
if [ -f "$GOROOT_DIR/lib/wasm/wasm_exec.js" ]; then
  cp "$GOROOT_DIR/lib/wasm/wasm_exec.js" "$OUT/playground/wasm_exec.js"
elif [ -f "$GOROOT_DIR/misc/wasm/wasm_exec.js" ]; then
  cp "$GOROOT_DIR/misc/wasm/wasm_exec.js" "$OUT/playground/wasm_exec.js"
else
  echo "error: wasm_exec.js not found under GOROOT=$GOROOT_DIR" >&2
  exit 1
fi

# Compile unified MyGo Showcase (UI 530 + Flow + Motion 28) to content-hashed WASM
GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$OUT/playground/mygox-ui.wasm" ./cmd/showcase
WASM_HASH="$(shasum -a 256 "$OUT/playground/mygox-ui.wasm" | cut -c1-10)"
WASM_NAME="mygox-ui.${WASM_HASH}.wasm"
mv "$OUT/playground/mygox-ui.wasm" "$OUT/playground/$WASM_NAME"
sed -i.bak "s|__WASM_FILE__|${WASM_NAME}|g" "$OUT/playground/index.html"
rm -f "$OUT/playground/index.html.bak"

touch "$OUT/.nojekyll"

ls -lh "$OUT/index.html" "$OUT/playground/index.html" "$OUT/playground/wasm_exec.js" "$OUT/playground/$WASM_NAME"
