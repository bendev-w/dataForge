#!/bin/sh
# Install DataForge: curl -fsSL https://raw.githubusercontent.com/bendev-w/dataForge/main/install.sh | sh
# Optional env: VERSION=v1.0.0  INSTALL_DIR=/custom/bin
set -eu

REPO="bendev-w/dataForge"
BIN="dataforge"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  linux|darwin) ;;
  *) echo "Unsupported OS: $os (on Windows use install.ps1)" >&2; exit 1 ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
esac

if [ -z "${VERSION:-}" ]; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n1)
  [ -n "$VERSION" ] || { echo "Could not determine latest version" >&2; exit 1; }
fi

archive="${BIN}_${VERSION#v}_${os}_${arch}.tar.gz"
base="https://github.com/$REPO/releases/download/$VERSION"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading $archive ..."
curl -fsSL "$base/$archive" -o "$tmp/$archive"
curl -fsSL "$base/checksums.txt" -o "$tmp/checksums.txt"

expected=$(grep " $archive\$" "$tmp/checksums.txt" | awk '{print $1}')
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$tmp/$archive" | awk '{print $1}')
else
  actual=$(shasum -a 256 "$tmp/$archive" | awk '{print $1}')
fi
[ -n "$expected" ] && [ "$expected" = "$actual" ] || { echo "Checksum verification failed" >&2; exit 1; }

tar -xzf "$tmp/$archive" -C "$tmp" "$BIN"

dest="${INSTALL_DIR:-/usr/local/bin}"
if [ ! -w "$dest" ] && [ -z "${INSTALL_DIR:-}" ]; then
  if command -v sudo >/dev/null 2>&1 && [ -d "$dest" ]; then
    sudo install -m 755 "$tmp/$BIN" "$dest/$BIN"
    echo "Installed $BIN $VERSION to $dest"
    exit 0
  fi
  dest="$HOME/.local/bin"
fi
mkdir -p "$dest"
install -m 755 "$tmp/$BIN" "$dest/$BIN"
echo "Installed $BIN $VERSION to $dest"
case ":$PATH:" in *":$dest:"*) ;; *) echo "Add $dest to your PATH" ;; esac
