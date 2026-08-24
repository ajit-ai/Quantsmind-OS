# Quantis

An AI-native operating system forked from [9front](https://9front.org)
(the actively developed Plan 9): every resource is a file, an API endpoint,
and a stream of events.

- **[DESIGN.md](DESIGN.md)** — pillars: resource protocol, namespaces as sandboxes,
  vertical+horizontal scheduling, memory reuse engine, event fabric, AI subsystem
- **[ROADMAP.md](ROADMAP.md)** — phased plan from first boot to cluster inference

## Quick Start (Windows host)

```powershell
# 1. One-time: QEMU + Hypervisor Platform acceleration
winget install SoftwareFreedomConservancy.QEMU
wsl --install -d Ubuntu        # optional, for Linux-side tooling

# 2. Get 9front ISO and install to a virtual disk
.\tools\get-iso.ps1
.\tools\run-qemu.ps1 -Install   # boots installer; follow FQA4 prompts

# 3. Daily boot
.\tools\run-qemu.ps1
```

After install, inside 9front you build from source with `git/pull` + `mk`
(self-hosting workflow — see ROADMAP Phase 1).

## Runs on Windows, Linux, macOS, BSD

`quantisd` (the hosted Quantis resource server) is pure Go stdlib —
no cgo, no syscalls. Build official binaries for every platform:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools\build-hosted.ps1
# -> hosted/bin/{windows,linux,darwin,freebsd,openbsd,netbsd}-<arch>/quantisd
```

Run anywhere: `quantisd`, then browse `http://localhost:9090/ns`.

## Status

Phase 0 — environment ready; hosted prototype running on 8 platforms.
