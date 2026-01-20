#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
OUT_DIR="$ROOT_DIR/dist/release/plugins"
GO_BIN="${GO_BIN:-go}"
GOTOOLCHAIN="${GOTOOLCHAIN:-go1.25.6}"
export GOTOOLCHAIN

mkdir -p "$OUT_DIR"

$GO_BIN build -buildmode=plugin -o "$OUT_DIR/nowpayments.so" "$ROOT_DIR/plugins/nowpayments"
$GO_BIN build -buildmode=plugin -o "$OUT_DIR/coinbasecommerce.so" "$ROOT_DIR/plugins/coinbasecommerce"
$GO_BIN build -buildmode=plugin -o "$OUT_DIR/bitpay.so" "$ROOT_DIR/plugins/bitpay"
$GO_BIN build -buildmode=plugin -o "$OUT_DIR/coingate.so" "$ROOT_DIR/plugins/coingate"
$GO_BIN build -buildmode=plugin -o "$OUT_DIR/coinpayments.so" "$ROOT_DIR/plugins/coinpayments"

echo "Plugins built in $OUT_DIR"
