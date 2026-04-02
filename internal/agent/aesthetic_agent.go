package agent

import (
	"context"
	"time"

	"google.golang.org/genai"
	"github.com/henry-clone/internal/domain"
	"github.com/henry-clone/internal/llm"
)

const KeyAesthetic = "aesthetic"

// AestheticAgent analyzes property photos to extract a visual identity
// that drives the deck's design — colors, typography, mood.
// This is what makes each deck feel bespoke rather than templated.
func AestheticAgent(client *genai.Client) Agent {
	return Agent{
		Name:      "aesthetic_analysis",
		DependsOn: []string{"data_extraction"},
		Retries:   1,
		Timeout:   30 * time.Second,
		Fn: func(ctx context.Context, state *PipelineState) error {
			deal, _ := GetTyped[*domain.Deal](state, KeyDeal)

			if len(deal.PhotoURLs) == 0 {
				// No photos — use default aesthetic
				aesthetic := &llm.PropertyAesthetic{
					PrimaryColor:   "#1a1a2e",
					SecondaryColor: "#4a4a5a",
					AccentColor:    "#2563eb",
					BackgroundTone: "light",
					Style:          "professional",
					Mood:           "professional",
					FontSuggestion: "Inter",
					Description:    "Clean, professional design.",
				}
				applyAestheticToBrand(deal, aesthetic)
				state.Set(KeyAesthetic, aesthetic)
				return nil
			}

			aesthetic, err := llm.AnalyzePropertyPhotos(ctx, client, deal.PhotoURLs)
			if err != nil {
				// Non-fatal — fall back to defaults
				aesthetic = &llm.PropertyAesthetic{
					PrimaryColor:   "#1a1a2e",
					SecondaryColor: "#4a4a5a",
					AccentColor:    "#2563eb",
					BackgroundTone: "light",
					Style:          "professional",
					Mood:           "professional",
					FontSuggestion: "Inter",
					Description:    "Clean, professional design.",
				}
			}

			applyAestheticToBrand(deal, aesthetic)
			state.Set(KeyAesthetic, aesthetic)
			return nil
		},
	}
}

// applyAestheticToBrand updates the deal's brand with photo-derived colors.
// Only overrides if the brand is still the default (not custom-set by user).
func applyAestheticToBrand(deal *domain.Deal, a *llm.PropertyAesthetic) {
	if deal.Brand.ID == "default" || deal.Brand.ID == "" {
		deal.Brand.PrimaryColor = a.PrimaryColor
		deal.Brand.SecondaryColor = a.SecondaryColor
		deal.Brand.AccentColor = a.AccentColor
		if a.FontSuggestion != "" {
			deal.Brand.FontHeading = "'" + a.FontSuggestion + "', Georgia, serif"
		}
	}
}
