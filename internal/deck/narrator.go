// Package deck handles the assembly of deal decks from analyzed data.
// The key abstraction here is the Narrator interface — it decouples
// narrative generation (LLM) from deck assembly (templating).
// This means we can swap OpenAI for Anthropic, or use a stub for testing,
// without touching the deck builder.
package deck

import (
	"context"

	"github.com/henry-clone/internal/domain"
)

// Narrator generates human-readable narratives from structured deal data.
// This is the LLM boundary — everything behind this interface can be
// an API call, a local model, or a hardcoded stub.
type Narrator interface {
	ExecutiveSummary(ctx context.Context, deal *domain.Deal) (string, error)
	PropertyOverview(ctx context.Context, deal *domain.Deal) (string, error)
	MarketOverview(ctx context.Context, deal *domain.Deal) (string, error)
	DealThesis(ctx context.Context, deal *domain.Deal) (string, error)
}
