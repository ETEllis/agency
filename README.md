# Agency

**Your intent, multiplied.**

Agency is a terminal-first AI working organization. Staff a persistent office of autonomous agents — each with a distinct voice, personality, and role — that run indefinitely on a schedule or one-off, build in public, and coordinate through a shared ledger. Watch them think. Hear them speak. Approve or reject their proposals in real time.

---

## What's Built

The full 5-stage core runtime is complete and ships in this repo.

### Architecture

```
┌─────────────────────────────────────────────────────────┐
│  USER LAYER                                              │
│  TUI iMessage bubbles · Voice (kokoro/say) · Approval    │
│  lane · Bulletin board · Genesis wizard                  │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│  DETERMINISTIC LAYER  (Nested Temporal Cron Tree)        │
│  ScheduleNode tree · prompt_injection per node           │
│  Fires enriched WakeSignals into reactive layer          │
└────────────────────────┬────────────────────────────────┘
                         │ enriched WakeSignal
┌────────────────────────▼────────────────────────────────┐
│  REACTIVE LAYER  (Stateful Agent Runtime)                │
│  Redis event bus · Actor daemons · Ledger · Consensus    │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│  GIST COGNITIVE LAYER  (per-agent)                       │
│  Causal compression · ElasticStretch · Lattice state     │
│  Outputs: GISTVerdict + execution_intent                 │
└────────────────────────┬────────────────────────────────┘
                         │ ActionIntent
┌────────────────────────▼────────────────────────────────┐
│  MODEL ROUTING LAYER                                     │
│  ModelRouter · CredentialBroker · 4 provider adapters    │
│  5 hard gates (capability/auth/privacy/tools/budget)     │
│  Ollama-first soft scoring                               │
└────────────────────────┬────────────────────────────────┘
                         │ Result
┌────────────────────────▼────────────────────────────────┐
│  LEDGER  (append-only, single source of truth)           │
│  CommitCertificate · Quorum consensus · Snapshots        │
└─────────────────────────────────────────────────────────┘
```

**Full request chain:**
```
WakeSignal → GIST/ReasoningCore → ActionIntent → ModelRouter → ProviderAdapter → Result → Ledger
```

### Completed Stages

| Stage | What shipped |
|-------|-------------|
| **1 — Live Agent Foundation** | Voice (kokoro-onnx + macOS `say` fallback), DB poll scheduler, broadcast→TUI pipeline, env config |
| **2 — GIST Cognitive Layer** | `GISTAgentCore`, causal compression, `ElasticStretch`, `LatticeStore`, per-wake lattice persistence |
| **3 — Model Routing Layer** | `ModelRouter` (5 hard gates + soft scoring), `CredentialBroker`, Anthropic/Ollama/OpenAI/Gemini adapters, routing audit log |
| **4 — Core TUI Experience** | iMessage-style bubbles (per-actor color + avatar + timestamp), TTS voice on broadcast, `ApprovalCmp` panel (a/r keys, auto right-rail), approval channel + vote relay |
| **5 — Nested Temporal Orchestration** | `ScheduleNode` tree with `prompt_injection`, `NestedScheduler`, `PerformanceRecord`, bulletin timeline (directive→output→score), daemon wired: directive → 1.5-weight GIST atom + performance publish |

---

## Quick Start

### Prerequisites

- Go 1.22+
- Redis (for multi-agent mode — optional for solo)
- Python 3.9+ with `kokoro-onnx` (for high-quality voice — falls back to `say`)
- An API key for at least one provider (Anthropic, OpenAI, Gemini) **or** Ollama running locally

### Install voice (optional)

```bash
scripts/install-voice
```

### Boot the office

```bash
# Build all daemons
scripts/build-daemons

# Set at least one provider key
export ANTHROPIC_API_KEY=sk-...   # or OPENAI_API_KEY, GEMINI_API_KEY, OLLAMA_API_BASE

# Start the office (TUI + scheduler + actor daemons)
overmind start
```

### TUI only (no daemons)

```bash
go build -o agency .
./agency
```

### Key commands inside the TUI

```
/agency genesis   — natural-language intent → structured org config
/agency bootstrap — boot a staffed office from a constitution
/agency status    — inspect running office, actors, schedules
/agency stop      — graceful shutdown
```

### Approval lane

When agents propose actions, the approval panel appears in the right rail automatically. Press `a` to approve, `r` to reject, `↑↓` to navigate.

### Bulletin board

Performance records (directive→output→score) stream into the messages viewport as agents complete inference cycles. Color-coded score badges shift green→yellow→red.

---

## Configuration

Primary config: `~/.agency.json` or `.agency.json` in your project root.

```jsonc
{
  "agency": {
    "productName": "Agency",
    "office": {
      "mode": "staffed",
      "sharedWorkplace": ".agency/workplace",
      "autoBoot": true
    },
    "redis": {
      "enabled": true,
      "address": "localhost:6379"
    },
    "schedules": {
      "defaultCadence": "@every 5m",
      "timezone": "America/New_York"
    },
    "currentConstitution": "coding-office"
  }
}
```

Legacy `.teamcode.json` / `.opencode.json` config files are still read as fallbacks for existing installs.

---

## Architecture — Key Files

| File | Role |
|------|------|
| `internal/agency/daemon_actor.go` | Actor main loop: GIST → routing → proposals → ledger → bus |
| `internal/agency/types.go` | All domain types including `ScheduleNode`, `ActionProposal`, `WakeSignal` |
| `internal/agency/gist_core.go` | GIST subprocess manager + elastic stretch |
| `internal/agency/nested_scheduler.go` | Cron tree with prompt injection |
| `internal/agency/routing.go` | `ModelRouter`, `CredentialBroker`, 5-gate scoring |
| `internal/agency/performance.go` | `PerformanceRecord`, `BulletinChannel`, `PublishPerformance` |
| `internal/agency/runtime.go` | `RuntimeManager`, channel helpers |
| `internal/tui/components/chat/approval.go` | Approval panel component |
| `internal/tui/components/chat/bulletin.go` | Bulletin timeline renderer |
| `internal/app/agency.go` | `AgencyService` — subscriptions, votes, genesis |
| `internal/db/migrations/` | 6 migrations (schema + agency runtime) |
| `AGENCY_BLUEPRINT.md` | Full architecture reference (canonical) |

---

## What's Next

### IPC Transport (in progress)
Unix socket server exposing the live office event stream to local clients (desktop app, CLI tools). Same `WakeSignal` / `LedgerEntry` schema as Redis but over a local socket.

### macOS Desktop App (in progress)
Native SwiftUI companion. Real-time office view: iMessage bubbles, bulletin board, approval lanes, agent status. Connects to the Go runtime via IPC socket. macOS-first, iPad companion after.

### WebSocket Transport
Remote client event stream — same schema as IPC, enabling web dashboard and mobile companion.

---

## Providers

Agency routes to whichever provider passes its gates. Set any subset of these:

```bash
export ANTHROPIC_API_KEY=...
export OPENAI_API_KEY=...
export GEMINI_API_KEY=...
export OLLAMA_API_BASE=http://localhost:11434   # Ollama preferred first (local-first)
```

---

## License

MIT
