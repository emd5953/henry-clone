package llm

import (
	"context"
	"fmt"

	"google.golang.org/genai"
	"github.com/henry-clone/internal/domain"
)

// GeminiNarrator implements deck.Narrator using Google's Gemini API.
type GeminiNarrator struct {
	client *genai.Client
	model  string
}

func NewGeminiNarrator(ctx context.Context, apiKey string) (*GeminiNarrator, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("creating gemini client: %w", err)
	}
	return &GeminiNarrator{
		client: client,
		model:  "gemini-2.5-flash",
	}, nil
}

func (n *GeminiNarrator) generate(ctx context.Context, system, user string) (string, error) {
	result, err := n.client.Models.GenerateContent(ctx, n.model,
		genai.Text(user),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: system}},
			},
			Temperature: genai.Ptr(float32(0.7)),
			MaxOutputTokens: 800,
		},
	)
	if err != nil {
		return "", fmt.Errorf("gemini: %w", err)
	}
	return result.Text(), nil
}

func (n *GeminiNarrator) ExecutiveSummary(ctx context.Context, d *domain.Deal) (string, error) {
	return n.generate(ctx,
		`You are a senior CRE analyst writing an executive summary for an offering memorandum. Professional tone, 2-3 paragraphs. Output HTML paragraphs only (no markdown, no code fences).`,
		fmt.Sprintf("Property: %s\nAddress: %s\nAsset Class: %s | %d units\nNOI: $%.0f\nOccupancy: %.1f%%\nAvg Rent: $%.0f\nThesis: %s",
			d.Property.Name, d.Property.Address.OneLiner(), d.Property.AssetClass, d.Property.Units,
			d.Analysis.NOI, d.Analysis.OccupancyRate*100, d.Analysis.AvgMonthlyRent, d.Thesis))
}

func (n *GeminiNarrator) PropertyOverview(ctx context.Context, d *domain.Deal) (string, error) {
	return n.generate(ctx,
		`You are a CRE analyst writing a property overview. Describe the asset, location, features. 2 paragraphs. HTML paragraphs only.`,
		fmt.Sprintf("Property: %s\nAddress: %s\nAsset Class: %s\nUnits: %d | SqFt: %d | Built: %d",
			d.Property.Name, d.Property.Address.OneLiner(), d.Property.AssetClass,
			d.Property.Units, d.Property.SqFt, d.Property.YearBuilt))
}

func (n *GeminiNarrator) MarketOverview(ctx context.Context, d *domain.Deal) (string, error) {
	return n.generate(ctx,
		`You are a CRE analyst writing a market overview. Submarket dynamics, demand drivers. 2 paragraphs. HTML paragraphs only.`,
		fmt.Sprintf("Location: %s\nAsset Class: %s\nOccupancy: %.1f%%",
			d.Property.Address.OneLiner(), d.Property.AssetClass, d.Analysis.OccupancyRate*100))
}

func (n *GeminiNarrator) DealThesis(ctx context.Context, d *domain.Deal) (string, error) {
	return n.generate(ctx,
		`You are a CRE analyst writing the investment thesis. Expand on the broker's thesis with data. 2-3 paragraphs. HTML paragraphs only.`,
		fmt.Sprintf("Thesis: %s\nNOI: $%.0f | Occupancy: %.1f%% | Expense Ratio: %.1f%%\nAvg Rent: $%.0f | Units: %d",
			d.Thesis, d.Analysis.NOI, d.Analysis.OccupancyRate*100, d.Analysis.ExpenseRatio*100,
			d.Analysis.AvgMonthlyRent, d.Analysis.TotalUnits))
}
