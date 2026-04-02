package llm

import (
	"context"
	"fmt"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/henry-clone/internal/domain"
)

// AnthropicNarrator implements deck.Narrator using Claude.
type AnthropicNarrator struct {
	client anthropic.Client
	model  anthropic.Model
}

func NewAnthropicNarrator(apiKey string) *AnthropicNarrator {
	return &AnthropicNarrator{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
		model:  anthropic.ModelClaudeSonnet4_5_20250929,
	}
}

func (n *AnthropicNarrator) generate(ctx context.Context, system, user string) (string, error) {
	msg, err := n.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     n.model,
		MaxTokens: 800,
		System: []anthropic.TextBlockParam{
			{Text: system},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(user)),
		},
		Temperature: anthropic.Float(0.7),
	})
	if err != nil {
		return "", fmt.Errorf("anthropic: %w", err)
	}
	if len(msg.Content) == 0 {
		return "", fmt.Errorf("anthropic returned no content")
	}
	return msg.Content[0].Text, nil
}

func (n *AnthropicNarrator) ExecutiveSummary(ctx context.Context, d *domain.Deal) (string, error) {
	return n.generate(ctx,
		`You are a senior CRE analyst writing an executive summary for an offering memorandum. Professional tone, 2-3 paragraphs. Output HTML paragraphs only (no markdown).`,
		fmt.Sprintf("Property: %s\nAddress: %s\nAsset Class: %s | %d units\nNOI: $%.0f\nOccupancy: %.1f%%\nAvg Rent: $%.0f\nThesis: %s",
			d.Property.Name, d.Property.Address.OneLiner(), d.Property.AssetClass, d.Property.Units,
			d.Analysis.NOI, d.Analysis.OccupancyRate*100, d.Analysis.AvgMonthlyRent, d.Thesis))
}

func (n *AnthropicNarrator) PropertyOverview(ctx context.Context, d *domain.Deal) (string, error) {
	return n.generate(ctx,
		`You are a CRE analyst writing a property overview. Describe the asset, location, features. 2 paragraphs. HTML only.`,
		fmt.Sprintf("Property: %s\nAddress: %s\nAsset Class: %s\nUnits: %d | SqFt: %d | Built: %d",
			d.Property.Name, d.Property.Address.OneLiner(), d.Property.AssetClass,
			d.Property.Units, d.Property.SqFt, d.Property.YearBuilt))
}

func (n *AnthropicNarrator) MarketOverview(ctx context.Context, d *domain.Deal) (string, error) {
	return n.generate(ctx,
		`You are a CRE analyst writing a market overview. Submarket dynamics, demand drivers. 2 paragraphs. HTML only.`,
		fmt.Sprintf("Location: %s\nAsset Class: %s\nOccupancy: %.1f%%",
			d.Property.Address.OneLiner(), d.Property.AssetClass, d.Analysis.OccupancyRate*100))
}

func (n *AnthropicNarrator) DealThesis(ctx context.Context, d *domain.Deal) (string, error) {
	return n.generate(ctx,
		`You are a CRE analyst writing the investment thesis. Expand on the broker's thesis with data. 2-3 paragraphs. HTML only.`,
		fmt.Sprintf("Thesis: %s\nNOI: $%.0f | Occupancy: %.1f%% | Expense Ratio: %.1f%%\nAvg Rent: $%.0f | Units: %d",
			d.Thesis, d.Analysis.NOI, d.Analysis.OccupancyRate*100, d.Analysis.ExpenseRatio*100,
			d.Analysis.AvgMonthlyRent, d.Analysis.TotalUnits))
}
