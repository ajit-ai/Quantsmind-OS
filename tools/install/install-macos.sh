#!/usr/bin/env bash
# One-click install for Quantis on macOS (Intel or Apple Silicon).
# Usage:
#   bash tools/install/install-macos.sh [amd64|arm64]
set -euo pipefail

ARCH="${1:-$(uname -m)}"          # amd64 or arm64
[ "$ARCH" = "x86_64" ] && ARCH="amd64"
[ "$ARCH" = "aarch64" ] && ARCH="arm64"

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
DEST="/usr/local/quantis"
ARCH_DIR="darwin-$ARCH"

echo "== Quantis one-click installer (macOS/$ARCH) =="

sudo mkdir -p "$DEST"
if [ ! -x "$ROOT/hosted/bin/$ARCH_DIR/quantisd" ]; then
    echo "building binaries..."
    ( cd "$ROOT/hosted" && GOOS=darwin GOARCH=$ARCH go build -o "bin/$ARCH_DIR/quantisd" . )
fi
sudo cp "$ROOT/hosted/bin/$ARCH_DIR/quantisd" "$DEST/quantisd"

( cd "$ROOT/hosted/reconciler" && GOOS=darwin GOARCH=$ARCH go build -o "$DEST/reconciler" . )

PLIST="$HOME/Library/LaunchAgents/com.quantis.server.plist"
echo "installing LaunchAgent..."
cat > "$PLIST" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>Label</key><string>com.quantis.server</string>
  <key>ProgramArguments</key>
  <array><string>$DEST/quantisd</string></array>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><true/>
  <key>StandardOutPath</key><string>/tmp/quantis.log</string>
  <key>StandardErrorPath</key><string>/tmp/quantis.err</string>
</dict></plist>
EOF

launchctl unload "$PLIST" 2>/dev/null || true
launchctl load "$PLIST"

echo ""
echo "Quantis installed to $DEST (LaunchAgent: com.quantis.server)"
echo "  resources:  http://localhost:9090/ns"
echo "  reconciler: http://localhost:9091/decl"
