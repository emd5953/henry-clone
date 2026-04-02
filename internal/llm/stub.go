package llm

import (
	"context"
	"fmt"

	"github.com/henry-clone/internal/domain"
)

// StubNarrator generates placeholder narratives without any LLM calls.
// Useful for development, testing, and demos without an API key.
// In a real system, you'd also use this for integration tests
// where you want deterministic output.
type StubNarrator struct{}

func NewStubNarrator() *StubNarrator {
	return &StubNarrator{}
}

func (s *StubNarrator) ExecutiveSummary(_ context.Context, deal *domain.Deal) (string, error) {
	return fmt.Sprintf(`<p>%s presents a compelling investment opportunity in the %s %s market. 
	The property features %d units with a current occupancy rate of %.1f%% and generates 
	a net operating income of $%.0f annually.</p>
	<p>With strong fundamentals and favorable market dynamics, this asset offers investors 
	an attractive risk-adjusted return profile with meaningful upside potential through 
	operational improvements and strategic repositioning.</p>`,
		deal.Property.Name,
		deal.Property.Address.City,
		deal.Property.AssetClass,
		deal.Property.Units,
		deal.Analysis.OccupancyRate*100,
		deal.Analysis.NOI,
	), nil
}

func (s *StubNarrator) PropertyOverview(_ context.Context, deal *domain.Deal) (string, error) {
	return fmt.Sprintf(`<p>%s is a %d-unit %s property located at %s. 
	Built in %d, the property comprises approximately %d square feet of rentable space 
	across a well-maintained physical plant.</p>
	<p>The property benefits from its strategic location with proximity to major 
	transportation corridors, employment centers, and retail amenities that drive 
	sustained tenant demand in the submarket.</p>`,
		deal.Property.Name,
		deal.Property.Units,
		deal.Property.AssetClass,
		deal.Property.Address.OneLiner(),
		deal.Property.YearBuilt,
		deal.Property.SqFt,
	), nil
}

func (s *StubNarrator) MarketOverview(_ context.Context, deal *domain.Deal) (string, error) {
	return fmt.Sprintf(`<p>The %s, %s %s market continues to demonstrate strong fundamentals 
	driven by population growth, job creation, and limited new supply. 
	The submarket has experienced steady rent growth over the trailing twelve months, 
	with vacancy rates remaining well below historical averages.</p>
	<p>Demand drivers include proximity to major employers, favorable demographic trends, 
	and ongoing infrastructure investment that supports long-term value appreciation 
	in the immediate trade area.</p>`,
		deal.Property.Address.City,
		deal.Property.Address.State,
		deal.Property.AssetClass,
	), nil
}

func (s *StubNarrator) DealThesis(_ context.Context, deal *domain.Deal) (string, error) {
	return fmt.Sprintf(`<p>%s</p>
	<p>The property's current NOI of $%.0f with a %.1f%% expense ratio suggests 
	meaningful operational upside. At %.1f%% occupancy with an average monthly rent 
	of $%.0f, there is clear runway to drive revenue growth through lease-up 
	of vacant units and mark-to-market rent adjustments on upcoming renewals.</p>`,
		deal.Thesis,
		deal.Analysis.NOI,
		deal.Analysis.ExpenseRatio*100,
		deal.Analysis.OccupancyRate*100,
		deal.Analysis.AvgMonthlyRent,
	), nil
}
