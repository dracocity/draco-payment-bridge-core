#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
OUT_DIR="$ROOT_DIR/plugins/dist"

mkdir -p "$OUT_DIR"

go build -buildmode=plugin -o "$OUT_DIR/nowpayments.so" "$ROOT_DIR/plugins/src/nowpayments"
go build -buildmode=plugin -o "$OUT_DIR/coinbasecommerce.so" "$ROOT_DIR/plugins/src/coinbasecommerce"
go build -buildmode=plugin -o "$OUT_DIR/bitpay.so" "$ROOT_DIR/plugins/src/bitpay"
go build -buildmode=plugin -o "$OUT_DIR/coingate.so" "$ROOT_DIR/plugins/src/coingate"
go build -buildmode=plugin -o "$OUT_DIR/coinpayments.so" "$ROOT_DIR/plugins/src/coinpayments"

echo "Plugins built in $OUT_DIR"
