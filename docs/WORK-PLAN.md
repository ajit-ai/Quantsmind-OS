# Quantis Work Plan — Per-Phase Work Statements

Each phase is a contract: Goal → You do → Deliverables → Done when → Estimate.
Standing rules: never break the self-hosted build · userspace-first ·
small commits · monthly upstream merges from 9front.

---

## Phase 0 — Environment Operational
**Goal:** bootable 9front VM on this machine.
**You do:** enable Windows Hypervisor Platform; download ISO (`tools/get-iso.ps1`);
run installer VM (`tools/run-qemu.ps1 -Install`); complete FQA4 install;
snapshot the disk image.
**Deliverables:** `vm/quantis.qcow2` (installed system) + `vm/SNAPSHOT.md` notes.
**Done when:** `tools/run-qemu.ps1` boots to rio desktop and you can open acme.
**Estimate:** 1 evening.
**Status:** QEMU 11.1.0 installed ✔ · ISO pending download.

## Phase 1 — Fork & Daily Loop
**Goal:** your branch builds *inside* the VM.
**You do:** fetch upstream mirror into this repo; branch `quantis/main` from
`front`; set up drawterm or serial console; add minimal kernel device
`sys/src/9/pc64/devqmos.c`; rebuild kernel in VM; reboot into it.
**Deliverables:** kernel whose `cat /dev/qmos` prints `quantis`.
**Done when:** you are booting a kernel you compiled yourself.
**Estimate:** 2–4 weekends.

## Phase 2 — Internals Mastery
**Goal:** own the mental model before modifying anything deep.
**You do:** read and annotate `sys/src/9/port/proc.c`, scheduler, xalloc/pager,
`dev.c` device table, `lib9p` server library, `cmd/ip` stack; write one doc each.
**Deliverables:** `docs/internals/*.md` (≥6 walkthroughs).
**Done when:** you can trace `write(2)` → dev write handler from memory.
**Estimate:** 6–10 weeks part-time.

## Phase 3 — Event Fabric (P6)
**Goal:** watch() semantics proven userspace-first, then kernel.
**You do:**
3a. extend `/hosted` Go prototype (exists ✔): multi-watcher fan-out,
     replay buffers, backpressure policy;
3b. kernel: vnode watchpoints, `watch(2)` syscall, stream channel delivery;
3c. demo dashboard subscribed to /proc inside acme.
**Deliverables:** `watch(2)` man page + live acme dashboard.
**Done when:** an unrelated process receives /proc events live over 9P.
**Estimate:** 8–12 weeks.
**Status:** 3a scaffold running (REST+SSE resource events) ✔

## Phase 4 — Resource Protocol (P1)
**Goal:** REST↔9P++ gateway and self-describing namespaces.
**You do:** schema/introspection extensions; gateway daemon (Go, hosted first);
streaming semantics (chunked reads = subscriptions); auth via factotum bridging.
**Done when:** plain `curl POST` mutates a namespace resource and fires events.
**Estimate:** 8–12 weeks.

## Phase 5 — Memory Reuse Engine (P5)
**Goal:** arena/pool allocators, explicit lifetimes, zero-copy channels.
**You do:** pool allocator behind existing malloc interfaces; lifetime API for
buffers; benchmark harness; unify pager view over RAM/device memory.
**Done when:** reuse benchmarks beat baseline with zero kernel-build regressions.
**Estimate:** 10–14 weeks.

## Phase 6 — AI Subsystem (P7)
**Goal:** `/dev/accel*`, model cache FS, inference daemon.
**You do:** CPU-SIMD backend first (llama.cpp-style kernels, kencc-compatible C);
weights mapped as pages through Phase 5 engine; inference server as 9P++ daemon.
**Done when:** a quantized small model answers a prompt from a Quantis shell
and streams tokens to a REST client.
**Estimate:** 12–16 weeks.

## Phase 7 — Horizontal Scaling (P4)
**Goal:** remote accelerators mountable; cluster inference transparent.
**You do:** `importd` capability protocol; multi-node QEMU topology;
scheduler placement across imported resources.
**Done when:** node A serves tokens from weights resident on node B.
**Estimate:** 12–16 weeks.

---

**Total: ≈18–24 months part-time to full vision.**
Every phase ends with a bootable snapshot commit.
