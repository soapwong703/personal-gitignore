#!/usr/bin/env sh
set -eu

BIN_DIR="${HOME}/.local/bin"
API_BASE="https://api.github.com/repos/soapwong703/personal-gitignore/releases/latest"
RELEASE_BASE="https://github.com/soapwong703/personal-gitignore/releases/latest/download"

download_to_file() {
  url="$1"
  output="$2"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$output"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$output" "$url"
  else
    echo "Error: curl or wget is required to download the release asset." >&2
    exit 1
  fi
}

latest_version() {
  if command -v curl >/dev/null 2>&1; then
    response=$(curl -fsSL -H 'Accept: application/vnd.github+json' "$API_BASE")
  elif command -v wget >/dev/null 2>&1; then
    response=$(wget -qO- --header='Accept: application/vnd.github+json' "$API_BASE")
  else
    echo "Error: curl or wget is required to query the latest release." >&2
    exit 1
  fi

  version=$(printf '%s\n' "$response" | sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)
  if [ -z "$version" ]; then
    echo "Error: unable to determine the latest release version." >&2
    exit 1
  fi

  printf '%s\n' "$version"
}

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

case "$(uname -s)" in
  Darwin) OS="darwin" ;;
  Linux) OS="linux" ;;
  *)
    echo "Error: install.sh supports macOS and Linux only. Use install.ps1 on Windows." >&2
    exit 1
    ;;
esac

ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *)
    echo "Error: unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

ASSET="pgi_${OS}_${ARCH}.tar.gz"
URL="${RELEASE_BASE}/${ASSET}"
VERSION="$(latest_version)"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
ARCHIVE="$TMP_DIR/$ASSET"

mkdir -p "$BIN_DIR"

echo "Downloading pgi ${VERSION} for ${OS}/${ARCH}..."
download_to_file "$URL" "$ARCHIVE"

tar -C "$TMP_DIR" -xzf "$ARCHIVE"

PKG_DIR="$TMP_DIR/pgi_${OS}_${ARCH}"
install -m 0755 "$PKG_DIR/pgi" "$BIN_DIR/pgi"

echo "Installed pgi ${VERSION} to $BIN_DIR/pgi"
