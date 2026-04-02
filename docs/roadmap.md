# Roadmap

What's left to build to reach feature parity with Henry.ai, ordered by impact.

## Phase 1: Core Gaps (High Impact)

### Async Job Queue
Currently deck generation blocks the HTTP request. Need to:
- Return 202 Accepted immediately with deal ID
- Run pipeline in background goroutine
- Frontend polls for status updates
- Show live progress as each agent completes

### Database Persistence
Everything is in-memory — deals vanish on restart. Need:
- PostgreSQL or SQLite for deal storage
- Migrate in-memory maps to database queries
- Store generated HTML and section data
- File storage for uploaded documents and photos (S3 or local disk)

### PDF/Word Document Ingestion
Henry handles "up to 5 unstructured documents per deal." We only handle CSV and Excel. Need:
- PDF text extraction (pdftotext or a Go library)
- LLM-assisted parsing for unstructured PDFs (send pages to Gemini, ask it to extract structured data)
- Word document parsing (docx library)
- Intelligent document type detection (is this a rent roll? a T12? a lease abstract?)

### Autonomous Agent Recovery
Current pipeline skips agents on dependency failure. Henry's system reroutes:
- If comps fail, narrative agent should still write about the market without comp data
- If market data fails, use whatever the LLM knows about the submarket
- Agents should adapt their prompts based on what data is available
- Partial results should still produce a usable deck

## Phase 2: Data Quality (Medium Impact)

### Real Comparable Sales
Options:
- ATTOM Data API (free tier available, property transactions)
- Firecrawl scraping of LoopNet/Crexi (already have the API key from LeaseIQ)
- County assessor public records

### Real Demographics
- US Census API (free, requires key)
- BLS employment data
- Walk Score API

### Real Maps
- Google Maps Static API (needs key, ~$2/1000 requests)
- Mapbox (free tier available)
- Embed interactive maps in deck HTML

## Phase 3: Polish (Medium Impact)

### Figma Write-to-Canvas
Push deck content directly into Figma as native frames:
- Create template files per deck type
- Map sections to Figma frames
- Populate text, tables, images programmatically
- QC team edits in Figma, exports final PDF

### Per-Brokerage Learning
Henry "knows how you talk." Need:
- Store previous decks per client
- Use few-shot examples in LLM prompts
- Learn preferred section ordering, tone, terminology
- Brand configuration UI

### Authentication
- JWT-based auth
- Multi-tenant (deals scoped to organizations)
- Role-based access (broker, analyst, QC reviewer, admin)

## Phase 4: Production (Lower Impact, Required for Launch)

### Monitoring & Observability
- Structured logging
- Request tracing
- Agent execution metrics
- Error alerting

### Rate Limiting & Security
- API rate limiting
- Input validation and sanitization
- CORS lockdown (not wildcard)
- Crypto-random deal IDs (already done)

### Email Notifications
- Deck ready notification
- QC review assigned notification
- Review completed notification

### LinkedIn Deal Announcements
Henry generates these. Low effort — just another deck type with a shorter, social-media-optimized format.

## Estimated Effort

| Phase | Effort | Impact |
|---|---|---|
| Phase 1 | 2-3 weeks | Gets us to ~70% of Henry |
| Phase 2 | 1-2 weeks | Real data makes decks actually useful |
| Phase 3 | 2-3 weeks | Gets us to ~85% of Henry |
| Phase 4 | 1-2 weeks | Production-ready |

Total: ~8-10 weeks to approximate Henry's core product.
