# Quantis Design Document

> An AI-native operating system descended from Plan 9 / 9front:
> every resource is a file, an API endpoint, and a stream of events.

## Lineage

```
Bell Labs Plan 9 ──► 9front (active fork, self-hosted)
                        │
                        └──► Quantis (this fork)
                             AI-native resources, event fabric,
                             hierarchical scheduling, cluster-as-mount
```

We fork 9front because its core ideas are exactly what modern AI platforms
need, and its ~100k-line codebase is small enough for one person to master.

## Core Pillars

### P1 — Resource Protocol (9P++)
One protocol serves files, REST semantics, IPC, and events.
Every kernel object is simultaneously:
- a file in a namespace (`/ai/models/llama/infer`)
- an endpoint (`GET` = read, `POST` = write, native streaming)
- a subscribable event source (`watch /dev/net0`)

No separate web stack, no separate RPC layer, no separate message bus.

### P2 — Namespaces as Sandboxes
Per-process namespaces (Plan 9's signature feature) isolate AI agents:
each agent sees only the capabilities mounted into its namespace.
No root user concept; authority flows through capability servers.

### P3 — Vertical Threading (within one machine)
Hierarchical scheduler: per-core run queues, NUMA-aware placement,
priority inheritance, deadline classes for real-time inference paths.

### P4 — Horizontal Scaling (across machines)
Plan 9's distributed heritage revived: remote CPUs/GPUs/NPUs are
*mounted resources*, not cluster middleware.

    mount -a tcp!node2/gpu0 /dev/gpu0

A cluster is a filesystem operation, not Kubernetes.

### P5 — Memory Reuse Engine
Arena/pool allocators with explicit lifetimes.
Tensor memory shared zero-copy between processes.
GPU/NPU memory managed by the same pager as RAM.
Pressure-driven reclaim ordered by workload priority class.

### P6 — Event Fabric
Kernel-native publish/subscribe. Any resource can be watched;
deliveries arrive as streams through the resource protocol.
Services self-describe by exposing their schema in their namespace.

### P7 — AI Subsystem
`/dev/accel*` unified accelerator interface, model cache as filesystem,
inference daemons as ordinary userspace 9P servers. Models are files;
weights are pages; inference is a write followed by a streaming read.

## Architecture Sketch

```
        userland agents / shells / services
   ─────────────────────────────────────────────
        namespaces │ resource protocol (9P++)
   ─────────────────────────────────────────────
   scheduler(v+h) │ memory engine │ event fabric
   ─────────────────────────────────────────────
     drivers: accel, net, disk, input, display
```

## Orchestration Model (displaces Docker + Kubernetes)

K8s exists because Linux lacks native resource primitives. Quantis supplies
them in-kernel, so the entire orchestrator collapses into three daemons:
`reconciler` (desired state), `factotum` (auth/secrets), `importd`
(resource import/export) — plus the Venti content-addressed store as
registry. No apiserver, no etcd, no kubelet, no CNI, no sidecars.

| Docker/K8s function | Quantis mechanism |
|---|---|
| OCI image / registry | Venti snapshots served over 9P++ (`/n/registry/app@v12`) |
| Container runtime | process with private namespace (primary, not bolt-on) |
| cgroups / quotas | scheduler classes + memory-engine priority pools |
| kubelet | `importd` capability-aware daemon |
| etcd desired-state | declarative namespace script + event-driven reconciler |
| kube-scheduler | kernel scheduler; remote CPUs/GPUs visible as mounts |
| CoreDNS discovery | `/srv` publish + namespace resolution |
| CNI overlay network | none: 9P++ over TCP/TLS is the fabric |
| PV/PVC volumes | native network filesystem mounts |
| Secrets | factotum, scoped per namespace |
| Rolling update | bind new tree over old; revoke capability to drain |
| HPA autoscaler | pressure/queue-depth driven spawn and remote import |
| Jobs / batch | deadline-class scheduler queues |
| Prometheus metrics | `watch` on any resource file — zero agents |
| Ingress | built-in REST↔9P++ gateway |

Beyond all existing orchestrators, Quantis adds AI-native primitives:
cross-machine GPU/NPU mounting, models-as-paged-memory shared zero-copy,
tensor IPC via shared arenas, latency-SLO (deadline) scheduling for token
streams, kernel pub/sub event fabric, and capability security with no root.

## Language Stack

| Layer | Language | Share |
|---|---|---|
| Kernel, drivers, memory engine, inference cores | Plan 9 C dialect | ~65% |
| Userspace services: reconciler, importd, REST gateway | Go (`9fans.net/go`) | ~25% |
| Deployment manifests (= executable namespace scripts), build | rc + mk | ~8% |
| Boot/trap assembly (mostly read-only) | x86-64 | ~2% |

Hosted-side only (never inside the OS): Python for model
conversion/quantization tooling. No Rust/C++/JS in-tree: 9front's
self-hosted ken-C toolchain is the contract; POSIX idioms are unlearned,
not ported.

## Safety Model

Plan 9 C provides no language-level memory safety. Quantis compensates in
layers rather than switching languages:

1. **Language safety at the edges** — all network-facing services (gateway,
   reconciler, importd) are Go; untrusted input never reaches C first.
2. **Auditability at the core** — modified kernel surface stays within tens
   of thousands of lines; one person can review every changed file.
3. **Sanitizer CI** — identical C sources compiled under clang ASan/UBSan +
   fuzzing on the hosted side before entering the VM tree.
4. **Architectural containment** — per-process namespaces and factotum
   credential custody bound the blast radius of any surviving bug.

Rejected alternatives: Go kernel (GC pauses contradict P3 latency goals;
Biscuit OS showed viability only with heavy caveats); Rust rewrite
(breaks the self-hosted build contract; Redox precedent shows decade-scale
timelines). C++ rejected outright: no frontend in the 9front toolchain,
kernel-hostile runtime.

## Non-Goals
- POSIX compliance as a goal (APE bridge exists, not a target)
- Competing as a general-purpose desktop
- Backwards compatibility above correctness
