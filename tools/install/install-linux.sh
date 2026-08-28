#!/usr/bin/env bash
# One-click install for Quantis on Debian/Ubuntu (and most systemd Linux).
# Usage:
#   bash tools/install/install-linux.sh [amd64|arm64]
set -euo pipefail

ARCH="${1:-amd64}"
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"      # repo root
DEST="/opt/quantis"
ARCH_DIR="linux-$ARCH"

echo "== Quantis one-click installer (Linux/$ARCH) =="

# build binaries if not present
if [ ! -x "$ROOT/hosted/bin/$ARCH_DIR/quantisd" ]; then
    echo "building binaries..."
    ( cd "$ROOT/hosted" && GOOS=linux GOARCH=$ARCH go build -o "bin/$ARCH_DIR/quantisd" . )
fi

sudo mkdir -p "$DEST"
echo "copying quantisd..."
sudo cp "$ROOT/hosted/bin/$ARCH_DIR/quantisd" "$DEST/"

echo "building reconciler..."
( cd "$ROOT/hosted/reconciler" && GOOS=linux GOARCH=$ARCH go build -o "$DEST/reconciler" . )

echo "installing systemd units..."
sudo tee /etc/systemd/system/quantis.service >/dev/null <<'EOF'
[Unit]
Description=Quantis resource server
After=network.target

[Service]
ExecStart=/opt/quantis/quantisd
Restart=always
RestartSec=2

[Install]
WantedBy=multi-user.target
EOF

sudo tee /etc/systemd/system/quantis-reconciler.service >/dev/null <<'EOF'
[Unit]
Description=Quantis reconciler (declaration -> self-heal)
After=network.target

[Service]
ExecStart=/opt/quantis/reconciler
Restart=always
RestartSec=2

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now quantis.service
sudo systemctl enable --now quantis-reconciler.service

sleep 1
echo ""
echo "Quantis installed to $DEST (systemd: quantis, quantis-reconciler)"
echo "  resources:  http://localhost:9090/ns"
echo "  reconciler: http://localhost:9091/decl"
