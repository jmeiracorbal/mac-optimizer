#!/usr/bin/env sh
set -e

REPO="jmeiracorbal/mac-optimizer"
BIN="mac-optimizer"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

OS=$(uname -s)
if [ "$OS" != "Darwin" ]; then
  echo "This installer only supports macOS" >&2
  exit 1
fi

ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64)  ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"v\([^"]*\)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Could not determine latest version" >&2
  exit 1
fi

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

TARBALL="${BIN}_${VERSION}_darwin_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${TARBALL}"

echo "Downloading mac-optimizer v${VERSION} (darwin/${ARCH})..."
curl -fsSL "$URL" -o "${TMP}/${TARBALL}"
tar -xzf "${TMP}/${TARBALL}" -C "$TMP"

echo "Installing to ${INSTALL_DIR}..."
install -d "$INSTALL_DIR"
install -m 755 "${TMP}/${BIN}" "${INSTALL_DIR}/${BIN}"

echo ""
echo "mac-optimizer v${VERSION} installed."
"${INSTALL_DIR}/${BIN}" --version
