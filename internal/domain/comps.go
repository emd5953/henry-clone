package domain

import "time"

// Comp represents a comparable sale or lease transaction.
// Henry pulls these from public sources and proprietary data.
// Comps are critical for BOVs and OMs — they justify pricing.
type Comp struct {
	Address    Address   `json:"address"`
	SaleDate   time.Time `json:"sale_date"`
	SalePrice  float64   `json:"sale_price"`
	PricePerSF float64   `json:"price_per_sf,omitempty"`
	PricePerUnit float64 `json:"price_per_unit,omitempty"`
	CapRate    float64   `json:"cap_rate,omitempty"`
	Units      int       `json:"units,omitempty"`
	SqFt       int       `json:"sq_ft,omitempty"`
	AssetClass AssetClass `json:"asset_class"`
	Distance   float64   `json:"distance_miles,omitempty"` // from subject property
}

// MarketData holds demographic and economic data for the submarket.
// Henry enriches decks with this from external sources.
type MarketData struct {
	Population       int     `json:"population,omitempty"`
	PopulationGrowth float64 `json:"population_growth_pct,omitempty"` // YoY
	MedianIncome     float64 `json:"median_income,omitempty"`
	UnemploymentRate float64 `json:"unemployment_rate,omitempty"`
	MedianRent       float64 `json:"median_rent,omitempty"`
	VacancyRate      float64 `json:"vacancy_rate,omitempty"` // submarket
	RentGrowth       float64 `json:"rent_growth_pct,omitempty"`
	WalkScore        int     `json:"walk_score,omitempty"`
	TransitScore     int     `json:"transit_score,omitempty"`
}

// LocationMap holds generated map data for the deck.
type LocationMap struct {
	StaticMapURL  string `json:"static_map_url,omitempty"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	AerialMapURL  string `json:"aerial_map_url,omitempty"`
}
