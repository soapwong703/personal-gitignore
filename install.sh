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

mkdir -p "$BIN_DIR"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

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

download "${BASE_URL}/pgi" "${TMP_DIR}/pgi"
download "${BASE_URL}/personal-gitignore" "${TMP_DIR}/personal-gitignore"

cp "${TMP_DIR}/pgi" "${BIN_DIR}/pgi"
cp "${TMP_DIR}/personal-gitignore" "${BIN_DIR}/personal-gitignore"
chmod u+x "${BIN_DIR}/personal-gitignore" "${BIN_DIR}/pgi"

echo "Successfully installed:"
echo "${BIN_DIR}/pgi"
echo "Use \`pgi\` as the default command."
echo "${BIN_DIR}/personal-gitignore"
