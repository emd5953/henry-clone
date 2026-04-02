# Henry Clone

A reverse-engineered clone of [Henry.ai](https://henry.ai) — the AI-powered deal deck generator for commercial real estate brokers. Built to understand how Henry's technology works and replicate its core architecture.

## The Problem

Commercial real estate brokers spend 3-4 hours a day building presentations. An Offering Memorandum — the document that sells a property — is 40-60 pages, takes weeks to produce, and requires:

- Gathering financial data from messy spreadsheets (rent rolls, T12s, underwriting models)
- Researching comparable sales, demographics, and market trends
- Writing professional narratives for each section
- Designing a polished, branded document with maps, photos, and tables
- Coordinating with analysts, designers, and marketing teams

Every day a deck isn't done is a day the broker's competitor takes their deal. Brokerages spend heavily on back-office staff (analysts, designers) just to produce these materials, and the process is still slow, manual, and error-prone.

## The Solution

Henry.ai (and this clone) automates the entire pipeline:

1. **Upload** — Broker drops in a rent roll, T12, photos, and a 2-3 sentence thesis
2. **Analyze** — AI agents run financial analysis, fetch comps, pull demographics, and generate maps — all in parallel
3. **Generate** — LLM writes professional narratives (executive summary, market overview, investment thesis) tailored to the deal
4. **Design** — The deck is assembled with per-client branding (fonts, colors, logos) and photo-driven aesthetics
5. **Review** — Human QC team reviews and polishes in Figma before delivery
6. **Deliver** — Broker gets a weblink + PDF in minutes instead of weeks

The result: brokers close more deals faster, with less back-office overhead, and every deck looks like someone spent weeks on it.

## How It Works

```
Broker uploads files (rent roll, T12, photos)
        ↓
  Multi-agent pipeline (parallel execution)
  ├── Data extraction & validation
  ├── Financial analysis (NOI, occupancy, cap rates)
  ├── Comparable sales enrichment
  ├── Market demographics enrichment
  ├── Location/map data enrichment
  ├── Photo aesthetic analysis (Gemini Vision)
  └── AI narrative generation (Gemini)
        ↓
  Deck assembly (branded HTML sections)
        ↓
  QC review workflow → Figma editing → PDF export
```

## Henry.ai vs This Clone

| Capability | Henry.ai | This Clone |
|---|---|---|
| Multi-agent orchestration | Golang engine, 5+ agents per deal | ✅ Go pipeline with parallel agents, retry, dependency resolution |
| Document ingestion | Excel, PDF, Word with LLM extraction | ✅ CSV + Excel with fuzzy column matching |
| AI narratives | LLM-generated (model unknown) | ✅ Gemini 2.5 Flash |
| Financial analysis | Full underwriting engine | ✅ NOI, occupancy, expense ratio, per-unit metrics |
| Deck types | OM, BOV, loan package, syndication, flyer, teaser | ✅ All six types with different section layouts |
| Per-client branding | Learns fonts, colors, logos, tone per brokerage | ✅ Brand struct with colors, fonts, logos |
| Photo-driven design | Unknown | ✅ Gemini Vision analyzes photos → design tokens |
| Comparable sales | Real data (likely CoStar/Reonomy) | ⚠️ Stub data |
| Demographics | Real data (Census, economic sources) | ⚠️ Stub data |
| Maps | Auto-generated (Google Maps/Mapbox) | ⚠️ Stub data |
| PDF export | Weblink + PDF delivery | ✅ Headless Chrome HTML→PDF |
| Deck editor | React-based, users tweak AI output | ✅ WYSIWYG contentEditable editor |
| QC review | Human team reviews every deck ~15 min | ✅ Claim, edit, approve/reject with audit trail |
| Figma integration | Likely used for design polish | ✅ Link files, comments, export |
| Database | Persistent (likely PostgreSQL) | ❌ In-memory (lost on restart) |
| Auth | Multi-tenant, per-brokerage | ❌ None |
| Async processing | Job queue, 202 + polling | ❌ Synchronous |
| SOC 2 compliance | Type I certified | ❌ None |

## Tech Stack

**Backend:** Go 1.26, Chi router, Gemini API, chromedp (PDF), excelize (Excel)

**Frontend:** React 18, TypeScript, Vite, Tailwind CSS v4

**AI:** Google Gemini 2.5 Flash (narratives + vision)

**Design:** Figma API integration, LeaseIQ-inspired UI

## Quick Start

```bash
# 1. Clone and install
git clone <repo>
cd henry-clone
go mod tidy
cd frontend && npm install && npm run build && cd ..

# 2. Set environment variables
cp .env.example .env
# Edit .env with your API keys

# 3. Run
bash scripts/run.sh
# Open http://localhost:8080
```

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `GEMINI_API_KEY` | Yes | Google Gemini API key for narratives + vision |
| `FIGMA_TOKEN` | Optional | Figma personal access token for deck editing |
| `PORT` | No | Server port (default: 8080) |

## Project Structure

```
cmd/server/          → Entry point, dependency wiring
internal/
  agent/             → Multi-agent pipeline engine
  api/               → HTTP handlers (deals, reviews, Figma)
  deck/              → Deck builder, renderers, templates
  domain/            → Core business types (Deal, Property, Brand, Review)
  enrichment/        → External data providers (comps, market, geo)
  export/            → PDF generation via headless Chrome
  figma/             → Figma REST API client
  llm/               → Gemini narrator + vision analysis
  parser/            → CSV/Excel parsing with fuzzy column matching
frontend/            → React + TypeScript + Tailwind
sample_data/         → Test CSV and Excel files
scripts/             → Run and test scripts
docs/                → Architecture and research documentation
```

## API Endpoints

| Method | Path | Description |
|---|---|---|
| POST | `/api/deals` | Create deal + generate deck |
| GET | `/api/deals` | List all deals |
| GET | `/api/deals/:id` | Get deal details |
| GET | `/api/deals/:id/deck` | Get HTML deck |
| GET | `/api/deals/:id/deck.pdf` | Download PDF |
| GET | `/api/deals/:id/sections` | Get deck sections |
| PUT | `/api/deals/:id/sections/:idx` | Edit a section |
| GET | `/api/reviews` | QC review queue |
| POST | `/api/deals/:id/review/start` | Claim for review |
| POST | `/api/deals/:id/review/edit` | Edit during review |
| POST | `/api/deals/:id/review/complete` | Approve or reject |
| POST | `/api/deals/:id/figma/link` | Link Figma file |
| GET | `/api/deals/:id/figma` | Get Figma file info |
| POST | `/api/deals/:id/figma/comment` | Post Figma comment |

## Documentation

- [Henry.ai Technical Architecture](docs/henry-ai-architecture.md) — Research on how Henry.ai works
- [Clone Architecture](docs/clone-architecture.md) — How this clone is built
- [Agent Pipeline](docs/agent-pipeline.md) — Multi-agent orchestration design
- [Roadmap](docs/roadmap.md) — What's left to build

## Credits

Built as a technical analysis of [Henry.ai](https://henry.ai) (YC S24). Henry was founded by Sammy Greenwall and Adam Pratt, with Roman Martynenko as founding engineer. All research is based on publicly available information from interviews, press releases, and the Henry.ai website.
