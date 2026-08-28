#!/bin/sh
# One-click install for Quantis on FreeBSD (also OpenBSD/NetBSD with rc.conf var swap).
# Usage:
#   sudo sh tools/install/install-bsd.sh [amd64]
set -e

ARCH="${1:-amd64}"
ROOT=$(CDPATH= cd -- "$(dirname "$0")/../.." && pwd)
OS_SYS=$(uname)
# map OS to our binary dir + go target
case "$OS_SYS" in
  FreeBSD) GOTARGET="freebsd/$ARCH" ;;
  OpenBSD) GOTARGET="openbsd/$ARCH" ;;
  NetBSD)  GOTARGET="netbsd/$ARCH" ;;
  *) echo "unsupported: $OS_SYS"; exit 1 ;;
esac

DEST="/usr/local/quantis"
mkdir -p "$DEST"
goos=$(echo "$GOTARGET" | cut -d/ -f1)

echo "== Quantis one-click installer ($OS_SYS/$ARCH) =="
[ ! -x "$ROOT/hosted/bin/$goos-$ARCH/quantisd" ] && \
  ( cd "$ROOT/hosted" && GOOS=$goos GOARCH=$ARCH go build -o "bin/$goos-$ARCH/quantisd" . )
cp "$ROOT/hosted/bin/$goos-$ARCH/quantisd" "$DEST/quantisd"
( cd "$ROOT/hosted/reconciler" && GOOS=$goos GOARCH=$ARCH go build -o "$DEST/reconciler" . )

# rc.d service scripts
cat > /usr/local/etc/rc.d/quantis <<EOF
#!/bin/sh
# PROVIDE: quantis
# REQUIRE: NETWORKING
. /etc/rc.subr
name="quantis"
rcvar="quantis_enable"
command="$DEST/quantisd"
pidfile="/var/run/quantis.pid"
command_args="&"
load_rc_config \$name
: \${quantis_enable:="NO"}
run_rc_command "\$1"
EOF
chmod +x /usr/local/etc/rc.d/quantis

# enable at boot
sysrc -f /etc/rc.conf quantis_enable="YES" >/dev/null

echo ""
echo "Quantis installed to $DEST (rc.d: quantis)"
echo "  resources:  http://localhost:9090/ns"
echo "  reconciler: http://localhost:9091/decl"
