#!/usr/bin/env sh
set -eu

BIN_DIR="${HOME}/.local/bin"
RELEASE_BASE="https://github.com/soapwong703/personal-gitignore/releases/latest/download"

while [ "$#" -gt 0 ]; do
  case "$1" in
    --bin-dir)
      if [ "$#" -lt 2 ]; then
        echo "Error: --bin-dir requires a value" >&2
        exit 1
      fi
      BIN_DIR="$2"
      shift 2
      ;;
    *)
      echo "Error: unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
esac

ASSET_EXT="tar.gz"
if [ "$OS" = "windows" ]; then
  ASSET_EXT="zip"
fi

ASSET="pgi_${OS}_${ARCH}.${ASSET_EXT}"
URL="${RELEASE_BASE}/${ASSET}"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
ARCHIVE="$TMP_DIR/$ASSET"

mkdir -p "$BIN_DIR"

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$URL" -o "$ARCHIVE"
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$ARCHIVE" "$URL"
else
  echo "Error: curl or wget is required to download the release asset." >&2
  exit 1
fi

if [ "$ASSET_EXT" = "zip" ]; then
  unzip -q "$ARCHIVE" -d "$TMP_DIR"
else
  tar -C "$TMP_DIR" -xzf "$ARCHIVE"
fi

PKG_DIR="$TMP_DIR/pgi_${OS}_${ARCH}"
install -m 0755 "$PKG_DIR/pgi" "$BIN_DIR/pgi"

echo "Installed: $BIN_DIR/pgi"
