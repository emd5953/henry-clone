# Agent Pipeline Design

The multi-agent pipeline is the core of both Henry.ai and this clone. It orchestrates parallel, independent units of work that transform raw deal inputs into a polished deck.

## Why Agents?

Traditional approach: sequential function calls. If step 3 fails, everything fails.

Agent approach: independent units that run in parallel where possible, retry on failure, and degrade gracefully. This is critical for a system that depends on multiple external services (LLM APIs, data providers, map services) that can each fail independently.

## Pipeline DAG

```
data_extraction
    ├── financial_analysis
    ├── comps_fetch          ← parallel
    ├── market_data_fetch    ← parallel
    ├── geo_fetch            ← parallel
    └── aesthetic_analysis   ← parallel (if photos uploaded)
            │
            ▼
    narrative_generation     ← depends on financial_analysis + market_data
            │
            ▼
    deck_assembly            ← depends on financial_analysis + narrative + comps + geo
```

## Agent Definition

```go
type Agent struct {
    Name      string
    Fn        func(ctx context.Context, state *PipelineState) error
    DependsOn []string
    Retries   int
    Timeout   time.Duration
}
```

## Shared State

Agents communicate through `PipelineState` — a thread-safe key-value store:

```go
state.Set("analysis", analysisResult)
// Later, in another agent:
analysis, ok := GetTyped[*FinancialAnalysis](state, "analysis")
```

This avoids tight coupling between agents. Each agent reads what it needs and writes what it produces.

## Execution Flow

1. Pipeline resolves the dependency graph
2. Agents with no dependencies start immediately (in parallel goroutines)
3. Each agent waits for its dependencies to complete
4. If a dependency failed, the dependent agent is skipped
5. Retries happen within each agent (configurable)
6. Timeouts are per-agent (not global)
7. Results are collected after all agents finish

## Failure Handling

| Scenario | Behavior |
|---|---|
| Comps API timeout | Comps agent retries 2x, then fails. Deck generates without comps section. |
| Gemini rate limit | Narrative agent retries 1x. If still fails, deck generation fails (critical). |
| Photo analysis fails | Aesthetic agent falls back to default design tokens. Non-fatal. |
| Financial analysis fails | All downstream agents skip. Deck generation fails (critical). |

Critical agents: `financial_analysis`, `narrative_generation`, `deck_assembly`
Non-critical agents: `comps_fetch`, `market_data_fetch`, `geo_fetch`, `aesthetic_analysis`

## Agent Inventory

| Agent | Dependencies | Retries | Timeout | Critical? |
|---|---|---|---|---|
| `data_extraction` | none | 1 | 30s | Yes |
| `financial_analysis` | data_extraction | 0 | 5s | Yes |
| `comps_fetch` | data_extraction | 2 | 15s | No |
| `market_data_fetch` | data_extraction | 2 | 15s | No |
| `geo_fetch` | data_extraction | 2 | 10s | No |
| `aesthetic_analysis` | data_extraction | 1 | 30s | No |
| `narrative_generation` | financial_analysis, market_data | 1 | 60s | Yes |
| `deck_assembly` | financial_analysis, narrative, comps, geo | 0 | 10s | Yes |

## Comparison to Henry

Henry's system (from Roman Martynenko's interview):

> "One agent retrieves zoning data. Another summarizes it. A third ranks the risk. If one part fails, the system reroutes and flags it. Each step is testable, and each output is traceable and auditable."

Our implementation follows the same principles:
- Independent agents with clear responsibilities
- Parallel execution where dependencies allow
- Graceful degradation on non-critical failures
- Traceable results (each agent reports status, duration, attempts, errors)

The main difference: Henry's agents are likely more specialized (separate agents for zoning, risk ranking, etc.) and may use LLM-based routing to decide which agents to invoke based on the deal type and available data.
