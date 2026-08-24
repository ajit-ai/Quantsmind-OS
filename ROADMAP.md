# Quantis Roadmap (9front fork path)

Reality check: 9front is self-hosted. You do not build it on Windows/Linux —
you boot an ISO once, install to disk, then `git/pull` + `mk` *inside* 9front.
Development loop: edit code → push to your git server → pull inside VM.

## Phase 0 — Environment (day 1)
- [ ] Install WSL2 Ubuntu (`wsl --install -d Ubuntu`) for QEMU + tooling
- [ ] Download latest 9front ISO from https://9front.org/releases/
- [ ] Boot ISO in QEMU (`tools/run-qemu.ps1 -Install`), install to virtual disk
- [ ] Snapshot the working disk image
- Learn: rc shell, mk build tool, Plan 9 C dialect, rio, acme

## Phase 1 — Fork & Daily Loop (week 1–2)
- [ ] Fetch upstream history: `git remote add 9front https://github.com/9front/9front.git && git fetch 9front front`
- [ ] Create branch `quantis/main` from `front`
- [ ] Set up drawterm or serial console workflow into the VM
- [ ] First change: add a custom `/dev/qmos` kernel device; rebuild kernel inside VM; boot it

## Phase 2 — Master the Internals (month 1–3)
Read + annotate these subsystems (the whole point of forking):
- [ ] `sys/src/9/port/` — portable kernel: proc, sched, page allocator
- [ ] `sys/src/9/pc64/` — amd64 machine layer: paging, interrupts, APIC
- [ ] `sys/src/lib9p/` — 9P servers library (foundation of everything)
- [ ] `sys/src/cmd/ip/` — network stack
- Deliverable: internal wiki of subsystem notes in `docs/`

## Phase 3 — Event Fabric (month 3–6)          ← first real novelty
- [ ] Kernel watchpoints on file/vnode operations
- [ ] `watch` syscall + stream delivery over existing channels
- [ ] Userspace demo: subscribe to /proc changes, live dashboard in acme

## Phase 4 — Resource Protocol Extensions (month 6–10)
- [ ] HTTP/REST gateway translating to/from 9P (userspace server first)
- [ ] Resource schema introspection (self-describing namespaces)
- [ ] Streaming semantics: chunked reads as event subscriptions

## Phase 5 — Memory Reuse Engine (month 10–14)
- [ ] Arena/pool allocator behind existing malloc(2) interfaces
- [ ] Explicit lifetime annotations for buffers; zero-copy channel transport
- [ ] Unified pager covering RAM + device memory

## Phase 6 — AI Subsystem (year 1.5+)
- [ ] `/dev/accel*` abstraction; first backend = CPU SIMD, then Vulkan/GPU passthrough
- [ ] Model cache filesystem; weights mapped as pages
- [ ] Inference daemon as userspace 9P server speaking our extensions

## Phase 7 — Horizontal Scaling (year 2+)
- [ ] Mount remote accelerators (`mount tcp!node/gpu0`)
- [ ] Cluster scheduler exposing remote CPUs through namespaces
- [ ] Multi-node demo: inference request served across machines transparently

## Rules of Engagement
1. Never break the self-hosted build — every commit boots in the VM snapshot.
2. Upstream 9front stays fetchable; keep merge discipline monthly.
3. Every pillar lands userspace-first before kernel work.
4. Small commits; Plan 9 style: simplicity over features.
