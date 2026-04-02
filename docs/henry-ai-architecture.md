# Henry.ai — Technical Architecture Research

Research compiled from publicly available sources: YC launch page, founder interviews, press releases, founding engineer interviews, and the henry.ai website.

## Company Overview

- **Founded:** 2024 by Sammy Greenwall (CEO) and Adam Pratt (CTO)
- **Funding:** $4M seed (Susa Ventures + 1Sharpe Ventures), YC S24
- **Customers:** 5 of the top 10 US commercial brokerages, 80+ customers
- **Team:** Flat engineering structure, no managers. Engineers own everything from whiteboard to production.

## What Henry Does

Henry takes messy, unstructured CRE data (rent rolls, T12s, comps, underwriting models, photos) and produces polished, professional marketing materials in minutes instead of weeks.

### Document Types Generated
- Offering Memorandums (OMs) — the primary product, typically 40-60 pages
- Broker Opinions of Value (BOVs)
- Loan packages
- Syndication decks
- Leasing flyers
- Investment teasers
- LinkedIn deal announcements

### Delivery Format
- Weblink (shareable URL)
- PDF (downloadable)

## Technical Architecture

### Core Engine: Golang Orchestration

From founding engineer Roman Martynenko's interviews:

> "A Golang orchestration engine that chains AI agents to process up to five unstructured documents per deal, transforming Excel, PDF, or Word files into a structured, presentation-ready deck."

The engine is built around **autonomous agents** that:
- Run independently and can fail without taking down the pipeline
- Reroute around failures automatically
- Produce testable, traceable, auditable output at each step

### Agent Pipeline

A typical pipeline for one deal:

```
Document Upload (Excel, PDF, Word — up to 5 files)
    ↓
Agent 1: Data Extraction
    → Parse rent rolls, T12s, underwriting models
    → Handle messy formats (merged cells, inconsistent columns)
    → LLM-assisted extraction for unstructured docs
    ↓
Agent 2: Financial Analysis
    → NOI, cap rates, occupancy, expense ratios
    → Per-unit and per-sqft metrics
    → Underwriting model validation
    ↓
Agent 3: External Data Enrichment (parallel)
    ├── Comps (comparable sales/leases)
    ├── Demographics (population, income, employment)
    ├── Zoning data
    ├── Market trends
    └── Location/maps
    ↓
Agent 4: Narrative Generation (LLM)
    → Executive summary
    → Property overview
    → Market overview
    → Investment thesis
    → Section-specific narratives
    ↓
Agent 5: Deck Assembly
    → Per-client branding (fonts, colors, logos, tone)
    → Section layout per deck type
    → Photo placement
    → Map generation
    → Comp tables
    ↓
Human QC Review (~15 minutes per deck)
    → Real analysts and designers review
    → Visual polish in design tool (likely Figma)
    → Content accuracy check
    ↓
Delivery (weblink + PDF)
```

### Failure Recovery

> "If one part fails, the system reroutes and flags it. Each step is testable, and each output is traceable and auditable."

This means:
- If the comps agent fails, the deck still generates without comps
- If one narrative section fails, it retries independently
- Each agent's output is logged for debugging
- The system degrades gracefully rather than failing entirely

### Frontend: React Deck Editor

From the interviews:
- React-based deck builder that "cuts pitch creation time from hours to under just minutes"
- Users can adjust structure, content, and layout
- The AI handles the heavy lifting, users have "just enough control to confidently finalize"
- The editor evolved based on user feedback — they noticed users were spending time fine-tuning layouts

### Per-Client Customization

> "Henry builds a deck in your style and voice."

The system learns each brokerage's:
- Visual brand (fonts, colors, logos)
- Writing tone and style preferences
- Preferred section ordering
- Formatting conventions

This is likely implemented as a per-client configuration that's refined over time, possibly with few-shot examples from previous decks.

### Human-in-the-Loop QC

From the website:
> "AI handles the first draft. Real humans fill the gaps."
> "We have a quality control team that spends at least 15 minutes reviewing every deck before it's sent to you."

This is a critical differentiator — Henry doesn't just ship raw AI output. Every deck goes through human review for:
- Content accuracy
- Visual polish
- Brand consistency
- Data validation

## Infrastructure

### Security
- SOC 2 Type I compliant
- Encrypted by default
- Enterprise-grade security

### Scale
- Serves 80+ customers including 5 of the top 10 US brokerages
- Lean team (no traditional sales or extensive operations teams)
- Production-grade, audited systems

## Key Technical Decisions

1. **Go for the orchestration engine** — Performance, concurrency, and reliability for the agent pipeline
2. **React for the frontend** — Interactive deck editor with real-time editing
3. **Agent-based architecture** — Independent, retryable, traceable units of work
4. **Human QC as a feature** — Not a bug. The 15-minute review is part of the product promise.
5. **Per-client branding** — Not templates. Each deck is "uniquely built."

## Sources

- [YC Launch Page](https://ycombinator.com/launches/Lcp-henry-ai-copilot-for-commercial-real-estate-brokers) (Aug 2024)
- [Commercial Observer Interview](https://commercialobserver.com/2024/06/sammy-greenwall-adam-pratt-henry-ai/) (Jun 2024)
- [KeyCrew Interview with Sammy Greenwall](https://keycrew.co/journal/the-ai-revolution-in-commercial-real-estate-how-henry-ai-is-improving-deal-making-workflows/) (Jun 2025)
- [TechBullion Interview with Roman Martynenko](https://techbullion.com/interview-with-roman-martynenko-fullstack-software-engineer-founding-engineer-henry-ai/) (Jun 2025)
- [TechBullion: Roman Martynenko's AI Engineering Playbook](https://techbullion.com/from-api-failures-to-autonomous-recovery-roman-martynenkos-ai-engineering-playbook/) (May 2025)
- [Henry AI Substack: $4M Seed Announcement](https://henryai.substack.com/p/how-henry-ai-is-transforming-commercial) (Feb 2025)
- [PRNewswire: Seed Funding Announcement](https://www.prnewswire.com/news-releases/henry-ai-raises-4-million-in-seed-funding-to-automate-commercial-real-estate-transactions-302380986.html) (Mar 2025)
- [henry.ai website](https://henry.ai) (accessed Apr 2026)
