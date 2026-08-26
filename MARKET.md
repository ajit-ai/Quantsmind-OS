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
| Licensing | Open-source core + subscription + enterprise tiers |

---

## Product Tiers

| Tier | Price | What's in it | Who buys |
|---|---|---|---|
| **Community** | Free | Full Quantis OS, single-node inference, community support (GitHub) | Hobbyists, eval, education, startups |
| **Pro** | $50/node/mo | Multi-node clustering, web dashboard, TLS, backup/restore, email support | Small-mid businesses in production |
| **Enterprise** | $200/node/mo + custom | HIPAA/SOC2/FedRAMP compliance, air-gapped builds, 24/7 SLA, on-site support | Hospitals, defense, finance, industrial |

**Open source core** = kernel, resource protocol, inference daemon, CLI
tools. Builds community, eliminates "why trust you" barrier.

**Pro/Enterprise** = operational layer on top (dashboards, compliance,
support). Pays the bills.

---

## Hardware Model — You Don't Touch It

The product is a disk image. Partners make the boxes. Just like pfSense
sells software, Netgate makes the hardware.

| Method | How it works |
|---|---|
| **DIY** | Customer flashes Quantis ISO onto any x86-64 box with a GPU |
| **Reference HW** | You publish "Recommended Hardware List" — tested configs. Customer buys whatever they want, you guarantee it works. |
| **Certified Partner** | HW vendors (NVIDIA Jetson reseller, Dell) pre-install Quantis. They pay a per-unit license fee or co-sell. |
| **Branded Appliance** (year 3+) | Partner manufactures Quantis-branded boxes (like Netgate/pfSense). Margin per unit. |

Year 1–2: DIY + reference list only. Zero inventory. Zero manufacturing.
Year 3+: Partner program + branded appliances if volume justifies.

---

## Concrete Use Case: MedScan Hospital

**Customer:** MedScan, 200 MRI scanners across 12 hospitals.
**Before Quantis:** Linux server room running Docker + K8s + Triton.
Three DevOps engineers at $150K/yr each. Each new model deployment
takes two weeks of YAML debugging.

### With Quantis

```
Day 1  — MedScan IT receives a Quantis flash drive
Day 1  — Flash onto NUC (Intel + NVIDIA T4) next to MRI
Day 1  — Boot → CLI → 5 commands:

         modelctl load /n/models/mri-seg-v4 -to /dev/accel0
         listen -r tcp!*!8080 export /ai/mri
         watch /ai/mri/stats > /dev/null &

Day 1  — MRI streams to Quantis box → results back to PACS
Day 1  — IT opens https://quantis-console.medscan.local
         Dashboard: GPU util, latency, model version — live, zero config
```

**Time to production: 2 hours** instead of 2 weeks. That's the sale.

### What MedScan Pays

| | Linux stack | Quantis |
|---|---|---|
| DevOps labor | 3 × $150K = $450K/yr | 0 |
| Software licenses | K8s+Docker (free) + Triton (free) | $200/node/mo × 200 = $480K/yr |
| **Total cost** | **$450K+ (labor) + cloud bills** | **$480K (all-in)** |
| **Time to deploy** | 2 weeks per model | 2 hours |
| **Ongoing maintenance** | Full-time DevOps team | Automated self-healing |

At 200+ nodes the math tips decisively. At 500+ nodes Quantis saves
$1M+/yr vs Linux+K8s stack. The CFO signs the deal.

### The "Why Not Just Use Linux?" Answer

> "You can. That's what you're doing today. How many DevOps engineers
> does it take to keep it running?"

The answer is usually 1–3 people at $150K+/year. That's $150K–$450K/year
in labor alone — more than Quantis Enterprise for 200 nodes.

**Quantis eliminates the Ops in MLOps.**

---

## Revenue Projection (Conservative)

| Year | Customers (nodes) | ARPU (annual) | ARR |
|---|---|---|---|
| 1 | 50 (10 customers × 5 nodes) | $600/yr | $30K |
| 2 | 500 (50 customers × 10 nodes) | $1,200/yr | $600K |
| 3 | 3,000 (200 customers × 15 nodes) | $1,800/yr | $5.4M |
| 4 | 10,000+ (enterprise deals) | $2,400/yr | $24M |

Breakeven at ~500 nodes (Year 2), solo dev + one support hire.

---

## License Structure

```
Quantis Community     MIT/BSL    — full OS, single-node
Quantis Pro           BSL 1.1    — multi-node, dashboard, TLS
Quantis Enterprise    Proprietary — compliance, SLA, custom builds
Quantis Model Hub     Proprietary — certified model zoo, auto-updates
```

BSL 1.1 (Business Source License) for Pro: source-available, converts
to MIT after 4 years. Prevents cloud providers from reselling your work
while keeping community trust.
