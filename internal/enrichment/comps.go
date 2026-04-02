// Package enrichment handles external data fetching.
// Henry pulls comps, demographics, zoning, and maps from
// public sources and proprietary data. These agents run
// in parallel during the pipeline and write to shared state.
package enrichment

import (
	"context"
	"fmt"
	"time"

	"github.com/henry-clone/internal/domain"
)

// CompsProvider fetches comparable sales/leases for a property.
// In production, this would hit CoStar, Reonomy, or public records APIs.
type CompsProvider interface {
	FetchComps(ctx context.Context, prop domain.Property, limit int) ([]domain.Comp, error)
}

// StubCompsProvider returns realistic placeholder comps for development.
type StubCompsProvider struct{}

func NewStubCompsProvider() *StubCompsProvider {
	return &StubCompsProvider{}
}

func (s *StubCompsProvider) FetchComps(_ context.Context, prop domain.Property, limit int) ([]domain.Comp, error) {
	if limit <= 0 {
		limit = 5
	}

	// Generate plausible comps near the subject property
	comps := []domain.Comp{
		{
			Address:    domain.Address{Street: "1300 Market St", City: prop.Address.City, State: prop.Address.State, Zip: prop.Address.Zip},
			SaleDate:   time.Now().AddDate(0, -3, 0),
			SalePrice:  4_200_000,
			PricePerSF: 385,
			CapRate:    0.058,
			SqFt:       10_900,
			AssetClass: prop.AssetClass,
			Distance:   0.4,
		},
		{
			Address:    domain.Address{Street: "850 Harrison Blvd", City: prop.Address.City, State: prop.Address.State, Zip: prop.Address.Zip},
			SaleDate:   time.Now().AddDate(0, -6, 0),
			SalePrice:  6_750_000,
			PricePerSF: 410,
			CapRate:    0.052,
			SqFt:       16_450,
			AssetClass: prop.AssetClass,
			Distance:   0.8,
		},
		{
			Address:    domain.Address{Street: "2100 Third St", City: prop.Address.City, State: prop.Address.State, Zip: prop.Address.Zip},
			SaleDate:   time.Now().AddDate(0, -9, 0),
			SalePrice:  3_100_000,
			PricePerSF: 355,
			CapRate:    0.063,
			SqFt:       8_730,
			AssetClass: prop.AssetClass,
			Distance:   1.2,
		},
		{
			Address:    domain.Address{Street: "475 Brannan St", City: prop.Address.City, State: prop.Address.State, Zip: prop.Address.Zip},
			SaleDate:   time.Now().AddDate(-1, 0, 0),
			SalePrice:  8_900_000,
			PricePerSF: 425,
			CapRate:    0.049,
			SqFt:       20_940,
			AssetClass: prop.AssetClass,
			Distance:   1.5,
		},
		{
			Address:    domain.Address{Street: "990 Illinois St", City: prop.Address.City, State: prop.Address.State, Zip: prop.Address.Zip},
			SaleDate:   time.Now().AddDate(-1, -2, 0),
			SalePrice:  5_400_000,
			PricePerSF: 370,
			CapRate:    0.055,
			SqFt:       14_590,
			AssetClass: prop.AssetClass,
			Distance:   0.6,
		},
	}

	if limit < len(comps) {
		comps = comps[:limit]
	}

	return comps, nil
}

// MarketDataProvider fetches demographic and economic data.
type MarketDataProvider interface {
	FetchMarketData(ctx context.Context, addr domain.Address) (*domain.MarketData, error)
}

// StubMarketDataProvider returns placeholder market data.
type StubMarketDataProvider struct{}

func NewStubMarketDataProvider() *StubMarketDataProvider {
	return &StubMarketDataProvider{}
}

func (s *StubMarketDataProvider) FetchMarketData(_ context.Context, addr domain.Address) (*domain.MarketData, error) {
	return &domain.MarketData{
		Population:       892_280,
		PopulationGrowth: 1.2,
		MedianIncome:     112_449,
		UnemploymentRate: 3.8,
		MedianRent:       3_200,
		VacancyRate:      5.4,
		RentGrowth:       3.8,
		WalkScore:        86,
		TransitScore:     80,
	}, nil
}

// GeoProvider fetches location/map data.
type GeoProvider interface {
	FetchLocation(ctx context.Context, addr domain.Address) (*domain.LocationMap, error)
}

// StubGeoProvider returns placeholder geo data.
type StubGeoProvider struct{}

func NewStubGeoProvider() *StubGeoProvider {
	return &StubGeoProvider{}
}

func (s *StubGeoProvider) FetchLocation(_ context.Context, addr domain.Address) (*domain.LocationMap, error) {
	return &domain.LocationMap{
		Latitude:     37.7599,
		Longitude:    -122.3894,
		StaticMapURL: fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?center=%s&zoom=15&size=800x400&maptype=roadmap", addr.OneLiner()),
	}, nil
}
