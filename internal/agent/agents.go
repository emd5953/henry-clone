package agent

import (
	"context"
	"time"

	"github.com/henry-clone/internal/domain"
	"github.com/henry-clone/internal/enrichment"
)

// State keys used by agents to communicate through PipelineState.
const (
	KeyDeal           = "deal"
	KeyAnalysis       = "analysis"
	KeyComps          = "comps"
	KeyMarketData     = "market_data"
	KeyLocation       = "location"
	KeyNarratives     = "narratives"
	KeySections       = "sections"
	KeyDeck           = "deck"
)

// Narratives is defined in the deck package (deck.Narratives).
// Agents use that type via the pipeline state.

// DataExtractionAgent validates and normalizes the uploaded data.
// In Henry's system, this is where unstructured docs get parsed.
func DataExtractionAgent() Agent {
	return Agent{
		Name:    "data_extraction",
		Retries: 1,
		Timeout: 30 * time.Second,
		Fn: func(ctx context.Context, state *PipelineState) error {
			deal, _ := GetTyped[*domain.Deal](state, KeyDeal)
			// Data is already parsed at upload time for now.
			// This agent validates and normalizes.
			if len(deal.RentRoll.Units) == 0 {
				// Not fatal — some deck types don't need a rent roll
			}
			return nil
		},
	}
}

// FinancialAnalysisAgent runs deterministic financial calculations.
// No LLM needed — pure math. Fast and reliable.
func FinancialAnalysisAgent() Agent {
	return Agent{
		Name:      "financial_analysis",
		DependsOn: []string{"data_extraction"},
		Timeout:   5 * time.Second,
		Fn: func(ctx context.Context, state *PipelineState) error {
			deal, _ := GetTyped[*domain.Deal](state, KeyDeal)
			analysis := domain.Analyze(deal)
			deal.Analysis = analysis
			state.Set(KeyAnalysis, analysis)
			return nil
		},
	}
}

// CompsAgent fetches comparable sales/leases from external sources.
// Runs in parallel with financial analysis — no dependency between them.
func CompsAgent(provider enrichment.CompsProvider) Agent {
	return Agent{
		Name:      "comps_fetch",
		DependsOn: []string{"data_extraction"},
		Retries:   2,
		Timeout:   15 * time.Second,
		Fn: func(ctx context.Context, state *PipelineState) error {
			deal, _ := GetTyped[*domain.Deal](state, KeyDeal)
			comps, err := provider.FetchComps(ctx, deal.Property, 5)
			if err != nil {
				return err
			}
			deal.Comps = comps
			state.Set(KeyComps, comps)
			return nil
		},
	}
}

// MarketDataAgent fetches demographics and economic data.
// Also parallel with financial analysis.
func MarketDataAgent(provider enrichment.MarketDataProvider) Agent {
	return Agent{
		Name:      "market_data_fetch",
		DependsOn: []string{"data_extraction"},
		Retries:   2,
		Timeout:   15 * time.Second,
		Fn: func(ctx context.Context, state *PipelineState) error {
			deal, _ := GetTyped[*domain.Deal](state, KeyDeal)
			data, err := provider.FetchMarketData(ctx, deal.Property.Address)
			if err != nil {
				return err
			}
			deal.MarketData = data
			state.Set(KeyMarketData, data)
			return nil
		},
	}
}

// GeoAgent fetches location and map data.
func GeoAgent(provider enrichment.GeoProvider) Agent {
	return Agent{
		Name:      "geo_fetch",
		DependsOn: []string{"data_extraction"},
		Retries:   2,
		Timeout:   10 * time.Second,
		Fn: func(ctx context.Context, state *PipelineState) error {
			deal, _ := GetTyped[*domain.Deal](state, KeyDeal)
			loc, err := provider.FetchLocation(ctx, deal.Property.Address)
			if err != nil {
				return err
			}
			deal.Location = loc
			state.Set(KeyLocation, loc)
			return nil
		},
	}
}
