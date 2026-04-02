# Henry Clone — CRE Deal Deck Generator

A Go backend that generates professional CRE deal decks from raw financial data.

## Quick Start

```bash
cd henry-clone
go mod tidy
go run ./cmd/server
# Then: ./scripts/test_deal.sh
```

Set `OPENAI_API_KEY` for real LLM narratives, or leave unset for stub mode.

## API

- `POST /api/deals` — Create deal (JSON or multipart with CSV uploads)
- `GET /api/deals/{id}` — Get deal JSON
- `GET /api/deals/{id}/deck` — Get rendered HTML deck
