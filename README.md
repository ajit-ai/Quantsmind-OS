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

## Status

Phase 0 — environment setup.
