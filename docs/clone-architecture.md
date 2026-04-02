# Clone Architecture

How this Henry.ai clone is built, the design decisions made, and how each component maps to Henry's production system.

## System Overview

```
┌─────────────────────────────────────────────────────┐
│                    React Frontend                    │
│  (Deal Creator, Deck Editor, QC Review, Figma)      │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP/JSON
┌──────────────────────▼──────────────────────────────┐
│                   Go HTTP Server                     │
│              (Chi router, handlers)                  │
├─────────────────────────────────────────────────────┤
│                 Agent Pipeline                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐            │
│  │  Data     │ │Financial │ │ Comps    │            │
│  │Extraction │ │Analysis  │ │ Fetch    │  parallel  │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘            │
│       │            │            │                    │
│  ┌────▼─────┐ ┌────▼─────┐ ┌────▼─────┐            │
│  │ Market   │ │   Geo    │ │Aesthetic │  parallel   │
│  │  Data    │ │  Fetch   │ │Analysis  │             │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘            │
│       └────────────┼────────────┘                    │
│                    ▼                                 │
│            ┌──────────────┐                          │
│            │  Narrative   │  (depends on analysis    │
│            │  Generation  │   + market data)         │
│            └──────┬───────┘                          │
│                   ▼                                  │
│            ┌──────────────┐                          │
│            │    Deck      │  (depends on all above)  │
│            │  Assembly    │                          │
│            └──────────────┘                          │
├─────────────────────────────────────────────────────┤
│              External Services                       │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐       │
│  │ Gemini │ │ Figma  │ │Comps   │ │ Chrome │       │
│  │  API   │ │  API   │ │(stub)  │ │(PDF)   │       │
│  └────────┘ └────────┘ └────────┘ └────────┘       │
└─────────────────────────────────────────────────────┘
```

## Backend (Go)

### Why Go?

Henry uses Go for their orchestration engine. We chose Go for the same reasons:
- Goroutines make parallel agent execution trivial
- Strong typing catches data modeling errors at compile time
- Fast compilation and startup
- Single binary deployment
- The `chromedp` library gives us headless Chrome for PDF export

### Package Structure

```
internal/
├── agent/          Pipeline engine
│   ├── pipeline.go       DAG executor with dependency resolution
│   ├── executor.go       Parallel execution, retry, timeout
│   ├── agents.go         Data, financial, comps, market, geo agents
│   ├── narrative_agent.go  LLM narrative + deck assembly agents
│   └── aesthetic_agent.go  Photo → design token analysis
│
├── api/            HTTP layer
│   ├── handler.go        Deal CRUD, sections, PDF export
│   ├── review.go         QC review workflow
│   └── figma.go          Figma integration endpoints
│
├── deck/           Deck generation
│   ├── builder.go        Orchestrates section assembly
│   ├── narrator.go       Narrator interface (LLM boundary)
│   ├── render.go         Cover, financials, rent roll renderers
│   ├── render_enriched.go  Comps, maps, demographics, photos, valuation
│   └── template.go       Branded HTML/CSS assembly
│
├── domain/         Core business types
│   ├── deal.go           Deal aggregate (property + financials + deck)
│   ├── property.go       Property, RentRoll, T12, Address types
│   ├── brand.go          Per-client branding + deck types
│   ├── comps.go          Comparable sales, market data, location
│   ├── review.go         QC review with audit trail
│   └── analyzer.go       Deterministic financial calculations
│
├── enrichment/     External data
│   └── comps.go          Provider interfaces + stubs
│
├── export/         Output formats
│   └── pdf.go            Headless Chrome HTML→PDF
│
├── figma/          Figma API
│   ├── client.go         REST API client
│   ├── types.go          Response types
│   └── bridge.go         Deal ↔ Figma file bridge
│
├── llm/            AI integration
│   ├── gemini.go         Gemini narrator (4 section types)
│   ├── vision.go         Photo aesthetic analysis
│   └── stub.go           Deterministic stub for testing
│
└── parser/         Document ingestion
    ├── csv.go            CSV rent roll + T12 parsing
    ├── excel.go          Excel parsing via excelize
    ├── fuzzy.go          Fuzzy column name matching
    └── document.go       File type detection
```

### Agent Pipeline Design

The pipeline is a DAG (directed acyclic graph) of agents. Each agent:
- Has a name, function, dependencies, retry count, and timeout
- Communicates through shared `PipelineState` (thread-safe key-value store)
- Runs in its own goroutine
- Waits for dependencies before executing
- Retries on failure (configurable per agent)
- Skips if a dependency failed

This mirrors Henry's architecture where "if one part fails, the system reroutes and flags it."

### Narrator Interface

The `Narrator` interface is the LLM boundary:

```go
type Narrator interface {
    ExecutiveSummary(ctx, deal) (string, error)
    PropertyOverview(ctx, deal) (string, error)
    MarketOverview(ctx, deal) (string, error)
    DealThesis(ctx, deal) (string, error)
}
```

Implementations: `GeminiNarrator` (production), `StubNarrator` (testing). Easy to add Claude, GPT-4, or any other model.

### Fuzzy Column Matching

Real broker data has inconsistent headers. Our parser handles:
- "Unit #", "unit_id", "Unit Number", "Ste", "Suite", "Space #" → `unit_id`
- "Monthly Rent", "Base Rent", "Contract Rent", "In Place Rent" → `monthly_rent`
- "Repairs & Maintenance", "R&M", "Building Maintenance" → `maintenance`

This is deterministic fuzzy matching. Henry likely also uses LLM-assisted extraction for the hardest cases.

## Frontend (React)

### Design System

Uses the LeaseIQ color palette:
- Background: `#F9F8F4` (warm off-white)
- Foreground: `#2D3A31` (dark green-gray)
- Primary: `#8C9A84` (sage green)
- Accent: `#C27B66` (terracotta)
- Secondary: `#DCCFC2` (warm beige)
- Fonts: Playfair Display (headings) + Source Sans 3 (body)

### Views

1. **Deal List** — Stats row + deal cards with status badges and financial metrics
2. **Deal Creator** — Form with drag-drop file uploads (rent roll, T12, photos)
3. **Deck Editor** — Section navigator + WYSIWYG contentEditable editor + formatting toolbar
4. **QC Review Queue** — Pending decks with claim/continue actions
5. **Review Editor** — Same as deck editor but with approve/reject + Figma panel

## How It Maps to Henry

| Henry Component | Our Implementation |
|---|---|
| Golang orchestration engine | `internal/agent/` pipeline with DAG execution |
| AI agents (5 per deal) | 7 agents: data, financial, comps, market, geo, aesthetic, narrative |
| Unstructured doc processing | CSV + Excel with fuzzy matching (no PDF/Word yet) |
| LLM narratives | Gemini 2.5 Flash via `google.golang.org/genai` |
| Per-client branding | `domain.Brand` struct applied to CSS generation |
| React deck editor | contentEditable WYSIWYG with section navigation |
| Human QC workflow | Review queue with claim, edit, approve/reject, audit trail |
| Figma integration | REST API client for file linking, comments, export |
| PDF delivery | chromedp headless Chrome rendering |
