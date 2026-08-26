# Quantis — Competitive Market Analysis

## Market Landscape (2026)

| Segment | Size (est.) | Growth | Incumbent | Vulnerability |
|---|---|---|---|---|
| Edge AI Inference | ~$25B by 2028 | 30%+ CAGR | Linux+Docker+K8s+CUDA | Massive stack for simple tasks |
| AI Model Serving | ~$15B by 2028 | 25% CAGR | NVIDIA Triton, vLLM, Ray Serve | Every option on bloated Linux |
| HPC / AI Training Clusters | ~$50B by 2028 | 15% CAGR | Slurm (1995!), PBS Pro | 30-year-old design nobody replaced |
| Embedded / Real-time | ~$8B by 2028 | 8% CAGR | QNX, VxWorks, FreeRTOS | Aging codebases, expensive licensing |
| PaaS / Cloud Deploy | ~$30B by 2028 | 20% CAGR | Heroku, Railway, Fly.io | Container-based, still complex |
| General OS Platform | ~$50B+ | 3% CAGR | Windows, Linux, macOS | Desktop wars — avoid |

---

## Segment #1 — Edge AI Inference (Primary Target)

**What:** Running AI models on physical hardware at factories, hospitals,
stores, autonomous vehicles — not in cloud datacenters.

**Current pain:** Hospital MRI scanner needs one segmentation model.
Today that requires Linux + Docker + NVIDIA toolkit + CUDA + K8s device
plugins + monitoring — minimum 15 components, DevOps team required.

**Incumbents:**
| Player | Weakness |
|---|---|
| NVIDIA JetPack + Linux | Vendor-locked, complex stack |
| AWS Greengrass | Cloud dependency, subscription |
| Azure IoT Edge | Cloud dependency |
| Slurm | Training-focused, not inference |

**Quantis advantage:** The OS IS the inference server. No cloud, no K8s.
Deploy in six commands. Sell the appliance box, not the platform.

**Revenue:** Certified hardware + OS license (per-device SaaS or perpetual).

---

## Segment #2 — AI Model Serving

**What:** Hosting LLMs behind an API developers call.

**Incumbents:**
| Player | Weakness |
|---|---|
| vLLM | Best open option, needs K8s at scale |
| NVIDIA Triton | NVIDIA-only, complex config |
| BentoML / Seldon | Python frameworks, still Linux |
| Replicate / Modal | Cloud-only, vendor lock-in |

**Quantis advantage:** No Python, no Docker, no K8s. Developer hits
`http://quantis-node:8080/v1/chat` — done. Model serving is the OS
itself, not a framework on top.

**Revenue:** Quantis Cloud — inference-as-a-service per token.

---

## Segment #3 — HPC Cluster Orchestration

**What:** Scheduling jobs across hundreds of GPUs for AI training.

**Incumbents:**
| Player | Weakness |
|---|---|
| Slurm | 30 years old, cluster admins retiring |
| PBS Pro | Same generation, niche |
| KubeFlow | K8s complexity ON TOP of HPC complexity |
| Ray | Great for Python, not an OS |

**Quantis advantage:** Plan 9 designed for distributed computing with
CPUs as mounted resources. Quantis adds GPU/NPU mounting — no Slurm
config, no MPI setup, no K8s.

**Revenue:** Annual license + support (Red Hat model). University/research
lab partnerships.

---

## Segment #4 — Embedded / Real-time

**What:** OS for cars, medical devices, industrial controllers, robotics.

**Incumbents:** QNX (automotive), VxWorks (industrial), FreeRTOS (MCU)

**Quantis advantage:** Capability security (no root) + event fabric
(real-time sensor fusion) + modern design. Automotive/robotics desperate
for QNX alternative.

**Revenue:** Per-device licensing ($2-15/unit, same as QNX model).

---

## Segment #5 — Secure Infrastructure

**What:** Government, defense, financial, critical infrastructure.

**Quantis advantage:** Plan 9 namespace isolation + no root + capability
auth. Stronger than anything Linux-based.

**Revenue:** Government contracts, consulting, custom builds.

---

## Segments to Avoid

| Segment | Reason |
|---|---|
| Desktop OS | Wars closed forever |
| General cloud compute | AWS/Azure/GCP infrastructure moat |
| Mobile phones | Android/iOS ecosystem lock-in |
| Consumer devices | No app ecosystem |

---

## Go-to-Market Sequence

```
YEAR 1   Edge AI Appliance — quickest to revenue, smallest market
         Ship box image + "Deploy LLM in 5 minutes" demo

YEAR 2   Model Serving SaaS — cloud inference, API pricing
         Quantis Cloud on rented GPU hardware

YEAR 3   HPC / Training — Slurm replacement
         Partner with one university lab for proof-of-concept

YEAR 4+  Embedded / Secure — longest sales cycles, biggest deals
         Automotive/defense partnerships
```

---

## Why Quantis Wins

| Factor | Edge |
|---|---|
| Complexity debt | Incumbents have 10-30 years of accumulated cruft. You start clean. |
| AI-native design | Nobody built an OS for AI. Everyone bolted AI onto Linux. |
| Developer experience | Six commands vs 200 lines of YAML. Sells itself in demos. |
| Cost structure | No K8s overhead = fewer servers = lower TCO |
| Licensing | Open-source base + support/enterprise tiers |
