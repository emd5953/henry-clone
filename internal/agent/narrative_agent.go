package agent

import (
	"context"
	"time"

	"github.com/henry-clone/internal/deck"
	"github.com/henry-clone/internal/domain"
)

// NarrativeAgent generates all LLM-powered text sections.
// Depends on financial analysis and market data so the narratives
// can reference computed metrics and enriched context.
func NarrativeAgent(narrator deck.Narrator) Agent {
	return Agent{
		Name:      "narrative_generation",
		DependsOn: []string{"financial_analysis", "market_data_fetch"},
		Retries:   1,
		Timeout:   60 * time.Second,
		Fn: func(ctx context.Context, state *PipelineState) error {
			deal, _ := GetTyped[*domain.Deal](state, KeyDeal)

			execSummary, err := narrator.ExecutiveSummary(ctx, deal)
			if err != nil {
				return err
			}

			propOverview, err := narrator.PropertyOverview(ctx, deal)
			if err != nil {
				return err
			}

			marketOverview, err := narrator.MarketOverview(ctx, deal)
			if err != nil {
				return err
			}

			thesis, err := narrator.DealThesis(ctx, deal)
			if err != nil {
				return err
			}

			state.Set(KeyNarratives, &deck.Narratives{
				ExecutiveSummary: execSummary,
				PropertyOverview: propOverview,
				MarketOverview:   marketOverview,
				DealThesis:       thesis,
			})
			return nil
		},
	}
}

// AssemblyAgent takes all computed data and narratives and builds
// the final deck HTML. This is the last step in the pipeline.
func AssemblyAgent(builder *deck.Builder) Agent {
	return Agent{
		Name: "deck_assembly",
		DependsOn: []string{
			"financial_analysis",
			"narrative_generation",
			"comps_fetch",
			"geo_fetch",
		},
		Timeout: 10 * time.Second,
		Fn: func(ctx context.Context, state *PipelineState) error {
			deal, _ := GetTyped[*domain.Deal](state, KeyDeal)
			narratives, _ := GetTyped[*deck.Narratives](state, KeyNarratives)

			if err := builder.BuildFromState(ctx, deal, narratives); err != nil {
				return err
			}

			state.Set(KeyDeck, deal.Deck)
			return nil
		},
	}
}
