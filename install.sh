#!/usr/bin/env sh
set -eu

BASE_URL="${PGI_INSTALL_BASE_URL:-https://raw.githubusercontent.com/soapwong703/personal-gitignore/main}"
BIN_DIR="${HOME}/.local/bin"

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

if ! command -v go >/dev/null 2>&1; then
  echo "Error: Go is required to install pgi." >&2
  exit 1
fi

mkdir -p "$BIN_DIR"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
SRC_DIR="$TMP_DIR/src"
mkdir -p "$SRC_DIR/cmd/pgi"

download() {
  url="$1"
  destination="$2"
  case "$url" in
    file://*)
      cp "${url#file://}" "$destination"
      ;;
    *)
      if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" -o "$destination"
      elif command -v wget >/dev/null 2>&1; then
        wget -qO "$destination" "$url"
      else
        echo "Error: curl or wget is required to download files." >&2
        exit 1
      fi
      ;;
  esac
}

download "${BASE_URL}/go.mod" "${SRC_DIR}/go.mod"
download "${BASE_URL}/cmd/pgi/main.go" "${SRC_DIR}/cmd/pgi/main.go"

(
  cd "$SRC_DIR"
  go build -o "$BIN_DIR/pgi" ./cmd/pgi
)
cp "$BIN_DIR/pgi" "$BIN_DIR/personal-gitignore"
chmod u+x "$BIN_DIR/pgi" "$BIN_DIR/personal-gitignore"

echo "Successfully installed:"
echo "${BIN_DIR}/pgi"
echo 'Use `pgi` as the default command.'
echo "${BIN_DIR}/personal-gitignore"
