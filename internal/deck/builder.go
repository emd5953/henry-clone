package deck

import (
	"context"
	"fmt"
	"time"

	"github.com/henry-clone/internal/domain"
)

// Builder orchestrates the full deck generation pipeline:
// 1. Run financial analysis
// 2. Generate narratives via LLM (or stub)
// 3. Assemble sections into a complete HTML deck
//
// This is the core "engine" of the product. The Builder now supports
// two modes: the original sequential Build() for backward compat,
// and BuildFromState() for the new multi-agent pipeline.
type Builder struct {
	narrator Narrator
}

func NewBuilder(n Narrator) *Builder {
	return &Builder{narrator: n}
}

// Narratives is a container for LLM-generated text passed from the pipeline.
type Narratives struct {
	ExecutiveSummary string
	PropertyOverview string
	MarketOverview   string
	DealThesis       string
}

// BuildFromState assembles a deck using pre-computed data from the agent pipeline.
// Unlike Build(), this doesn't run analysis or LLM calls — those already happened.
func (b *Builder) BuildFromState(ctx context.Context, deal *domain.Deal, narr *Narratives) error {
	sections := b.sectionsForDeckType(deal, narr)

	deal.Deck = &domain.Deck{
		HTML:        b.assembleHTML(deal, sections),
		GeneratedAt: time.Now(),
		Sections:    sections,
	}
	deal.Status = domain.StatusReady
	return nil
}

// sectionsForDeckType returns the appropriate sections based on deck type.
// Henry produces different layouts for OMs vs BOVs vs flyers.
func (b *Builder) sectionsForDeckType(deal *domain.Deal, narr *Narratives) []domain.Section {
	switch deal.DeckType {
	case domain.DeckTypeFlyer:
		return []domain.Section{
			{Type: domain.SectionCover, Title: deal.Property.Name, Content: b.renderCover(deal)},
			{Type: domain.SectionPropertyOverview, Title: "Property Overview", Content: narr.PropertyOverview},
			{Type: domain.SectionPhotos, Title: "Photos", Content: b.renderPhotos(deal)},
			{Type: domain.SectionFinancials, Title: "Financial Highlights", Content: b.renderFinancials(deal)},
		}
	case domain.DeckTypeBOV:
		return []domain.Section{
			{Type: domain.SectionCover, Title: deal.Property.Name, Content: b.renderCover(deal)},
			{Type: domain.SectionExecutiveSummary, Title: "Executive Summary", Content: narr.ExecutiveSummary},
			{Type: domain.SectionPropertyOverview, Title: "Property Overview", Content: narr.PropertyOverview},
			{Type: domain.SectionFinancials, Title: "Financial Summary", Content: b.renderFinancials(deal)},
			{Type: domain.SectionComps, Title: "Comparable Sales", Content: b.renderComps(deal)},
			{Type: domain.SectionValuation, Title: "Valuation Analysis", Content: b.renderValuation(deal)},
			{Type: domain.SectionMarketOverview, Title: "Market Overview", Content: narr.MarketOverview},
		}
	case domain.DeckTypeTeaser:
		return []domain.Section{
			{Type: domain.SectionCover, Title: deal.Property.Name, Content: b.renderCover(deal)},
			{Type: domain.SectionExecutiveSummary, Title: "Investment Highlights", Content: narr.ExecutiveSummary},
			{Type: domain.SectionPhotos, Title: "Photos", Content: b.renderPhotos(deal)},
			{Type: domain.SectionFinancials, Title: "Financial Summary", Content: b.renderFinancials(deal)},
			{Type: domain.SectionLocationMap, Title: "Location", Content: b.renderMap(deal)},
		}
	default: // OM is the default — the full 50-page beast
		return []domain.Section{
			{Type: domain.SectionCover, Title: deal.Property.Name, Content: b.renderCover(deal)},
			{Type: domain.SectionExecutiveSummary, Title: "Executive Summary", Content: narr.ExecutiveSummary},
			{Type: domain.SectionPropertyOverview, Title: "Property Overview", Content: narr.PropertyOverview},
			{Type: domain.SectionPhotos, Title: "Property Photos", Content: b.renderPhotos(deal)},
			{Type: domain.SectionLocationMap, Title: "Location & Maps", Content: b.renderMap(deal)},
			{Type: domain.SectionFinancials, Title: "Financial Summary", Content: b.renderFinancials(deal)},
			{Type: domain.SectionRentRoll, Title: "Rent Roll", Content: b.renderRentRoll(deal)},
			{Type: domain.SectionComps, Title: "Comparable Sales", Content: b.renderComps(deal)},
			{Type: domain.SectionDemographics, Title: "Demographics", Content: b.renderDemographics(deal)},
			{Type: domain.SectionMarketOverview, Title: "Market Overview", Content: narr.MarketOverview},
			{Type: domain.SectionDealThesis, Title: "Investment Thesis", Content: narr.DealThesis},
		}
	}
}

// Build takes a deal with raw inputs and produces a complete deck.
// It mutates the deal in place (adds Analysis and Deck).
func (b *Builder) Build(ctx context.Context, deal *domain.Deal) error {
	// Step 1: Financial analysis (deterministic, fast)
	deal.Status = domain.StatusAnalyzing
	deal.Analysis = domain.Analyze(deal)

	// Step 2: Generate narratives (LLM, slower)
	deal.Status = domain.StatusGenerating

	execSummary, err := b.narrator.ExecutiveSummary(ctx, deal)
	if err != nil {
		deal.Status = domain.StatusFailed
		return fmt.Errorf("generating executive summary: %w", err)
	}

	propOverview, err := b.narrator.PropertyOverview(ctx, deal)
	if err != nil {
		deal.Status = domain.StatusFailed
		return fmt.Errorf("generating property overview: %w", err)
	}

	marketOverview, err := b.narrator.MarketOverview(ctx, deal)
	if err != nil {
		deal.Status = domain.StatusFailed
		return fmt.Errorf("generating market overview: %w", err)
	}

	thesisNarrative, err := b.narrator.DealThesis(ctx, deal)
	if err != nil {
		deal.Status = domain.StatusFailed
		return fmt.Errorf("generating deal thesis: %w", err)
	}

	// Step 3: Assemble the deck
	sections := []domain.Section{
		{Type: domain.SectionCover, Title: deal.Property.Name, Content: b.renderCover(deal)},
		{Type: domain.SectionExecutiveSummary, Title: "Executive Summary", Content: execSummary},
		{Type: domain.SectionPropertyOverview, Title: "Property Overview", Content: propOverview},
		{Type: domain.SectionFinancials, Title: "Financial Summary", Content: b.renderFinancials(deal)},
		{Type: domain.SectionRentRoll, Title: "Rent Roll", Content: b.renderRentRoll(deal)},
		{Type: domain.SectionMarketOverview, Title: "Market Overview", Content: marketOverview},
		{Type: domain.SectionDealThesis, Title: "Investment Thesis", Content: thesisNarrative},
	}

	deal.Deck = &domain.Deck{
		HTML:        b.assembleHTML(deal, sections),
		GeneratedAt: time.Now(),
		Sections:    sections,
	}
	deal.Status = domain.StatusReady

	return nil
}
