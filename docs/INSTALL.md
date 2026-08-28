# Installing Quantis

Quantis ships as a one-click installer per OS. The hosted services
(`quantisd`, `reconciler`) are pure Go — zero runtime dependencies — so
installation is just copying binaries and registering autostart.

## Build the release package first

```powershell
# requires Go 1.21+ and (for .tar.gz) tar.exe (Win10+ has it)
powershell -NoProfile -ExecutionPolicy Bypass -File tools\package.ps1
# outputs dist/ with per-OS folders + .tar.gz archives
```

## Windows

```powershell
powershell -ExecutionPolicy Bypass -File tools\install\install-windows.ps1
# installs to %ProgramFiles%\Quantis, registers scheduled tasks,
# starts quantisd (:9090) and reconciler (:9091) now and at logon

powershell -ExecutionPolicy Bypass -File tools\install\uninstall-windows.ps1
```

## Debian / Ubuntu (systemd)

```bash
bash tools/install/install-linux.sh amd64      # or arm64
# installs to /opt/quantis, enables quantis.service + quantis-reconciler.service
```

## macOS (Intel / Apple Silicon)

```bash
bash tools/install/install-macos.sh arm64      # or amd64
# installs to /usr/local/quantis, loads LaunchAgent com.quantis.server
```

## FreeBSD / OpenBSD / NetBSD

```sh
sudo sh tools/install/install-bsd.sh amd64
# installs to /usr/local/quantis, enables rc.d service `quantis`
```

## Verify after install

```bash
curl http://localhost:9090/ns                    # namespace introspection
curl http://localhost:9090/r/sys/mem             # real memory (Windows/Linux)
curl -X POST -H "Content-Type: application/json" \
     -d '{"name":"svc","desired":3}' http://localhost:9091/decl
curl http://localhost:9091/decl/svc              # watch it self-heal
```

## Production signing (recommended before public release)

- **Windows** — sign `quantisd.exe` / `reconciler.exe` with a code-signing
  cert (Authenticode) so SmartScreen does not warn.
- **macOS** — sign with an Apple Developer ID and notarize:
  `codesign --options runtime` + `notarytool submit`.
- **Linux** — distribute via your own `.deb` repo or Flatpak store for trust.

## Packaging as the OS appliance (later phases)

When Phase 0–1 complete and Quantis boots on bare metal, `package.ps1`
gains a target that emits a bootable disk image (the current `vm/` ISO
becomes a user-installable `quantis.img` written with `dd`/Rufus). That
is the "final product" appliance distribution, distinct from the hosted
services installers documented above.
