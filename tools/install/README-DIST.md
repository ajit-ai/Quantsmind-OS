# Quantis — Binary Distribution

This folder contains prebuilt Quantis binaries for this platform:

- `quantisd` — the resource server (one object = file read + API write + event stream)
- `reconciler` — declaration → self-heal orchestration loop
- `tools/demo.ps1` — optional scripted usability demo

## One-click install

Pick your platform installer from `tools/install/`:

| OS | Command |
|---|---|
| Windows | `powershell -ExecutionPolicy Bypass -File tools\install\install-windows.ps1` |
| Debian/Ubuntu | `bash tools/install/install-linux.sh` |
| macOS | `bash tools/install/install-macos.sh` |
| FreeBSD/OpenBSD/NetBSD | `sudo sh tools/install/install-bsd.sh` |

## Manual run (no install)

```
./quantisd      # resource server on :9090
./reconciler    # orchestration loop on :9091
```

## Verify

- Resources (namespace):  `curl http://localhost:9090/ns`
- Live kernel-ish stats:  `curl http://localhost:9090/r/sys/mem`
- Self-heal demo:         `curl -X POST -H "Content-Type: application/json" \
                              -d '{"name":"svc","desired":3}' \
                              http://localhost:9091/decl`

## Note: these are unsigned community builds
macOS Gatekeeper and Windows SmartScreen may prompt. For production,
sign binaries (see docs/INSTALL.md) or build from source.
