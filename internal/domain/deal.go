package domain

import "time"

// Deal is the top-level aggregate that ties everything together.
// A deal = property + financials + thesis + generated deck.
// This is the unit of work in Henry's world.
type Deal struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Status    DealStatus `json:"status"`

	// Inputs — what the broker provides
	Property Property `json:"property"`
	RentRoll RentRoll `json:"rent_roll"`
	T12      T12      `json:"t12"`
	Thesis   string   `json:"thesis"` // 2-3 bullet points on why this deal

	// Configuration
	DeckType DeckType `json:"deck_type"`
	Brand    Brand    `json:"brand"`

	// Enrichment — pulled from external sources
	Comps      []Comp      `json:"comps,omitempty"`
	MarketData *MarketData `json:"market_data,omitempty"`
	Location   *LocationMap `json:"location,omitempty"`
	PhotoURLs  []string    `json:"photo_urls,omitempty"`

	// Derived — what we compute
	Analysis *FinancialAnalysis `json:"analysis,omitempty"`

	// Output — the generated deck
	Deck *Deck `json:"deck,omitempty"`

	// QC Review — human review before delivery
	Review *Review `json:"review,omitempty"`

	// Figma integration — deck editing in Figma
	FigmaFileKey string `json:"figma_file_key,omitempty"`
	FigmaFileURL string `json:"figma_file_url,omitempty"`
}

type DealStatus string

const (
	StatusPending    DealStatus = "pending"
	StatusAnalyzing  DealStatus = "analyzing"
	StatusGenerating DealStatus = "generating"
	StatusReady      DealStatus = "ready"
	StatusInReview   DealStatus = "in_review"
	StatusApproved   DealStatus = "approved"
	StatusFailed     DealStatus = "failed"
)

// FinancialAnalysis contains the key metrics any CRE buyer/seller cares about.
// These are derived from the T12 and rent roll — not user-provided.
type FinancialAnalysis struct {
	// Income metrics
	EffectiveGrossIncome float64 `json:"effective_gross_income"`
	TotalExpenses        float64 `json:"total_expenses"`
	NOI                  float64 `json:"noi"` // Net Operating Income = EGI - Expenses

	// Per-unit / per-sqft metrics (depends on asset class)
	NOIPerUnit float64 `json:"noi_per_unit,omitempty"`
	NOIPerSqFt float64 `json:"noi_per_sqft,omitempty"`

	// Occupancy
	TotalUnits    int     `json:"total_units"`
	OccupiedUnits int     `json:"occupied_units"`
	OccupancyRate float64 `json:"occupancy_rate"` // 0.0 - 1.0

	// Expense ratio
	ExpenseRatio float64 `json:"expense_ratio"` // expenses / EGI

	// Rent analysis
	AvgMonthlyRent  float64 `json:"avg_monthly_rent"`
	TotalAnnualRent float64 `json:"total_annual_rent"`

	// Investment metrics
	CapRate      float64 `json:"cap_rate,omitempty"`      // NOI / purchase price
	PricePerUnit float64 `json:"price_per_unit,omitempty"`
	PricePerSqFt float64 `json:"price_per_sqft,omitempty"`
	GRM          float64 `json:"grm,omitempty"` // Gross Rent Multiplier
}

// Deck is the generated output — the thing the broker sends to win deals.
type Deck struct {
	HTML      string    `json:"html"`
	GeneratedAt time.Time `json:"generated_at"`
	Sections  []Section `json:"sections"`
}

// Section represents a logical block of the deck.
// Decks aren't monolithic — they're composed of discrete sections
// that can be reordered, customized, or swapped per brokerage style.
type Section struct {
	Type    SectionType `json:"type"`
	Title   string      `json:"title"`
	Content string      `json:"content"` // HTML content for this section
}

type SectionType string

const (
	SectionCover            SectionType = "cover"
	SectionExecutiveSummary SectionType = "executive_summary"
	SectionPropertyOverview SectionType = "property_overview"
	SectionFinancials       SectionType = "financials"
	SectionRentRoll         SectionType = "rent_roll"
	SectionMarketOverview   SectionType = "market_overview"
	SectionDealThesis       SectionType = "deal_thesis"
	SectionComps            SectionType = "comparables"
	SectionLocationMap      SectionType = "location_map"
	SectionDemographics     SectionType = "demographics"
	SectionPhotos           SectionType = "photos"
	SectionValuation        SectionType = "valuation"
	SectionDebtSummary      SectionType = "debt_summary"
	SectionLeaseAbstract    SectionType = "lease_abstract"
)
