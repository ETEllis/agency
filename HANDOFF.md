# Agency — Session Handoff
> Written at ~95% context. Start a new session, read this file, say "go".

---

## What This Project Is

**Agency** — terminal-first AI working organization. Binary: `teamcode`. Product: **Agency**. Your intent, multiplied. Each agent has a voice, a bubble, a distinct personality, and runs indefinitely on a schedule or one-off. Full architecture spec: `AGENCY_BLUEPRINT.md v3`.

GitHub: `https://github.com/ETEllis/agency` (renamed from teamcode this session).

---

## Repository

```
/Users/edwardellis/teamcode/
```

```bash
cd /Users/edwardellis/teamcode
export PATH=$PATH:/usr/local/go/bin:/opt/homebrew/bin
go build ./...   # ✅
go test ./...    # ✅
```

Boot:
```bash
scripts/build-daemons && overmind start
```

---

## What's Complete (all verified ✅)

### Stages 1–5 — Full Core Runtime
| Stage | Summary |
|-------|---------|
| 1 | Voice (kokoro/say), DB poll scheduler, broadcast→TUI pipeline |
| 2 | GIST cognitive layer: causal compression, ElasticStretch, LatticeStore |
| 3 | Model routing: ModelRouter, CredentialBroker, 4 provider adapters |
| 4 | TUI: iMessage bubbles, TTS voice, ApprovalCmp panel (a/r keys) |
| 5 | NestedScheduler, PerformanceRecord, bulletin timeline, directive→GIST atom |

### IPC Transport Layer ✅ (shipped this session)
- `internal/agency/ipc.go` — `IPCServer`, newline-delimited JSON over Unix socket
  - Handshake → subscribe to broadcast + approval + bulletin channels
  - Fans out `IPCTypeBroadcast` / `IPCTypeApproval` / `IPCTypeBulletin` events
  - Accepts `IPCTypeVote` from clients → relays as `SignalCorrection` on org channel
  - `IPCSocketPath(baseDir, orgID)` = `{baseDir}/.agency/ipc-{orgID}.sock`
- `internal/agency/cmd/ipc-server/main.go` — standalone daemon (env: `AGENCY_ORG_ID`, `AGENCY_BASE_DIR`, `AGENCY_REDIS_ADDR`)

### Desktop App Scaffold ✅ (in progress — next session)
- `desktop/Agency/Package.swift` — SwiftPM, macOS 14+
- `desktop/Agency/Sources/Agency/Models.swift` — all domain types + IPC payload structs
- `desktop/Agency/Sources/Agency/Services/AgencyConnection.swift` — `@MainActor ObservableObject`, NWConnection Unix socket client, auto-reconnect, full JSON decode pipeline
- `desktop/Agency/Sources/Agency/AgencyApp.swift` — `@main`, `WindowGroup` + secondary `Window("Approvals")`
- `desktop/Agency/Sources/Agency/Views/OfficeView.swift` — `NavigationSplitView`, sidebar with badge counts, `ConnectionBadge`
- **MISSING** (next session): `BubbleListView.swift`, `ApprovalView.swift`, `BulletinView.swift`

---

## START HERE — Next Session Goal

**Complete the macOS desktop app.** Three views remain:

### 1. `BubbleListView.swift`
iMessage-style broadcast bubbles. Mirrors TUI `renderBroadcastBubble` in Swift.

Key data: `connection.broadcasts: [BroadcastMessage]` — each has `.actorID`, `.message`, `.createdAt`, `.initials`, `.roleColor`

Design:
- `ScrollViewReader` + `ScrollView`, auto-scroll to bottom on new message
- Each bubble: colored avatar circle (`.roleColor`) + actor name + timestamp header, message bubble with left thick border in matching color
- Map `RoleColor` cases → SwiftUI `Color` (blue, purple, teal, green, orange, pink)
- Empty state: "Waiting for office signals…" with a muted pulse animation

### 2. `ApprovalView.swift`
Pending proposal list. Key data: `connection.approvals: [ApprovalItem]`

Each row: actor name, action type → target, timestamp. Trailing swipe actions: approve (green checkmark) / reject (red x). Also keyboard: selected row, `return` = approve, `delete` = reject.

Calls: `connection.sendVote(proposalID: ..., approved: true/false)`

Empty state: "No pending approvals" with shield SF symbol.

### 3. `BulletinView.swift`
Performance timeline. Key data: `connection.bulletins: [BulletinEntry]`

Each row:
- Actor name + provider/model tag (muted, right-aligned) + timestamp
- Directive line in italic muted text
- Output text (truncated to 2 lines)
- Score badge: `String(format: "%.0f%%", entry.score * 100)`, color from `entry.scoreColor` (`.good` → green, `.fair` → yellow, `.poor` → red)

Empty state: "No performance records yet."

---

## Wiring It All Up

After the three views are written:

1. **Add `ipc-server` to `Procfile`**:
   ```
   ipc: ./bin/ipc-server
   ```

2. **Add `ipc-server` to `scripts/build-daemons`**:
   ```bash
   go build -o bin/ipc-server ./internal/agency/cmd/ipc-server
   ```

3. **SwiftUI app environment** — set these before launching:
   ```bash
   export AGENCY_ORG_ID=<your-org-id>
   export AGENCY_BASE_DIR=/Users/edwardellis/teamcode
   ```
   Or pass via Xcode scheme environment variables.

4. **Xcode project** — open `desktop/Agency/` as a SwiftPM package (`File → Open → desktop/Agency`) then add a macOS App target that depends on the Agency package target. OR: use the existing `Package.swift` directly (swift run works for dev, Xcode needed for signed .app).

---

## Testing Without Xcode Right Now

You don't need the desktop app to test the runtime. Run the TUI:

```bash
# Terminal 1 — start Redis (if not running)
redis-server

# Terminal 2 — set API key, boot office
export ANTHROPIC_API_KEY=sk-...   # or OPENAI_API_KEY, or start Ollama
scripts/build-daemons
overmind start

# In the TUI:
/agency genesis          # → describe what you want the org to do
/agency bootstrap        # → spins up actors per constitution
```

Watch agents wake on schedule → GIST compress → route to model → broadcast their output. Approval panel auto-shows when proposals arrive. Bulletin entries stream in after each inference cycle.

**To test IPC manually:**
```bash
# Build IPC server
go build -o bin/ipc-server ./internal/agency/cmd/ipc-server

# Run it (in background)
AGENCY_ORG_ID=your-org AGENCY_BASE_DIR=. AGENCY_REDIS_ADDR=localhost:6379 ./bin/ipc-server &

# Connect via nc
nc -U .agency/ipc-your-org.sock
# Then paste (one line):
{"type":"handshake","payload":{"orgId":"your-org","clientType":"cli"}}
# Events stream back as newline-delimited JSON
```

---

## Key Files Reference

| File | Role |
|------|------|
| `internal/agency/daemon_actor.go` | Actor main loop — full pipeline |
| `internal/agency/ipc.go` | Unix socket IPC server |
| `internal/agency/nested_scheduler.go` | Schedule tree with prompt injection |
| `internal/agency/performance.go` | PerformanceRecord + BulletinChannel |
| `internal/agency/routing.go` | ModelRouter + CredentialBroker |
| `internal/tui/components/chat/approval.go` | TUI approval panel |
| `internal/tui/components/chat/bulletin.go` | TUI bulletin renderer |
| `internal/app/agency.go` | AgencyService — all subscriptions + votes |
| `desktop/Agency/Sources/Agency/Services/AgencyConnection.swift` | macOS IPC client |
| `desktop/Agency/Sources/Agency/Views/OfficeView.swift` | Main window + sidebar |
| `AGENCY_BLUEPRINT.md` | Architecture reference (canonical) |

---

## Notes for Next Claude

- `defaultMode: "dontAsk"` in `~/.claude/settings.json` — Write/Edit/Bash all auto-approved.
- `.planning/config.json` disables context monitor hook (prevents 10s freezes). Do not delete.
- Go module path is still `github.com/ETEllis/teamcode` — do NOT change imports. The binary/module name doesn't need to match the GitHub repo name.
- The desktop app uses `NWConnection` (Network framework) for Unix socket — no third-party deps needed.
- `RoleColor` in `Models.swift` maps to SwiftUI colors in `BubbleListView` — make sure the mapping uses the same deterministic hash as Go's `actorColor()` in `list.go`.
- After desktop views are complete, the next phase per blueprint is **WebSocket transport** (item 29) for iPad companion, then **iPad app** (item 31).
- Do NOT redesign. Blueprint v3 is final. Execute the plan.
