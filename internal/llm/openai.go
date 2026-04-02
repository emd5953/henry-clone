package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"github.com/henry-clone/internal/domain"
)

// OpenAINarrator implements deck.Narrator using the OpenAI API.
// Each method crafts a specific prompt for its section type.
// We keep prompts focused and section-specific rather than trying
// to generate the entire deck in one shot — this gives us more
// control over quality and lets us retry individual sections.
type OpenAINarrator struct {
	client *openai.Client
	model  string
}

func NewOpenAINarrator(apiKey string) *OpenAINarrator {
	return &OpenAINarrator{
		client: openai.NewClient(apiKey),
		model:  openai.GPT4o,
	}
}

func (n *OpenAINarrator) generate(ctx context.Context, system, user string) (string, error) {
	resp, err := n.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: n.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
		Temperature: 0.7,
		MaxTokens:   800,
	})
	if err != nil {
		return "", fmt.Errorf("openai completion: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}

func (n *OpenAINarrator) ExecutiveSummary(ctx context.Context, deal *domain.Deal) (string, error) {
	system := `You are a senior commercial real estate analyst writing an executive summary 
for a property offering memorandum. Write in a professional, confident tone. 
Be concise — 2-3 paragraphs max. Focus on the investment highlights.
Output clean HTML paragraphs only (no markdown).`

	user := fmt.Sprintf(`Property: %s
Address: %s
Asset Class: %s | %d units
NOI: $%.0f
Occupancy: %.1f%%
Avg Monthly Rent: $%.0f
Broker's thesis: %s`,
		deal.Property.Name,
		deal.Property.Address.OneLiner(),
		deal.Property.AssetClass, deal.Property.Units,
		deal.Analysis.NOI,
		deal.Analysis.OccupancyRate*100,
		deal.Analysis.AvgMonthlyRent,
		deal.Thesis,
	)
	return n.generate(ctx, system, user)
}

func (n *OpenAINarrator) PropertyOverview(ctx context.Context, deal *domain.Deal) (string, error) {
	system := `You are a commercial real estate analyst writing a property overview section 
for a deal deck. Describe the physical asset, location advantages, and key features.
2 paragraphs max. Output clean HTML paragraphs only.`

	user := fmt.Sprintf(`Property: %s
Address: %s
Asset Class: %s
Units: %d | Sq Ft: %d | Year Built: %d`,
		deal.Property.Name,
		deal.Property.Address.OneLiner(),
		deal.Property.AssetClass,
		deal.Property.Units, deal.Property.SqFt, deal.Property.YearBuilt,
	)
	return n.generate(ctx, system, user)
}

func (n *OpenAINarrator) MarketOverview(ctx context.Context, deal *domain.Deal) (string, error) {
	system := `You are a commercial real estate analyst writing a market overview for a deal deck.
Discuss the submarket dynamics, demand drivers, and competitive landscape.
2 paragraphs max. Output clean HTML paragraphs only.`

	user := fmt.Sprintf(`Property location: %s
Asset Class: %s
Current occupancy: %.1f%%`,
		deal.Property.Address.OneLiner(),
		deal.Property.AssetClass,
		deal.Analysis.OccupancyRate*100,
	)
	return n.generate(ctx, system, user)
}

func (n *OpenAINarrator) DealThesis(ctx context.Context, deal *domain.Deal) (string, error) {
	system := `You are a commercial real estate analyst writing the investment thesis section 
of a deal deck. Expand on the broker's thesis with supporting financial data.
Make a compelling case. 2-3 paragraphs. Output clean HTML paragraphs only.`

	user := fmt.Sprintf(`Broker's thesis: %s
NOI: $%.0f | Occupancy: %.1f%% | Expense Ratio: %.1f%%
Avg Monthly Rent: $%.0f | Total Units: %d`,
		deal.Thesis,
		deal.Analysis.NOI,
		deal.Analysis.OccupancyRate*100,
		deal.Analysis.ExpenseRatio*100,
		deal.Analysis.AvgMonthlyRent,
		deal.Analysis.TotalUnits,
	)
	return n.generate(ctx, system, user)
}
