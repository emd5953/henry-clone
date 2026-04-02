package deck

import (
	"fmt"
	"strings"
	"time"

	"github.com/henry-clone/internal/domain"
)

// renderComps builds the comparable sales table.
// Henry auto-generates comp tables from external data sources.
func (b *Builder) renderComps(deal *domain.Deal) string {
	if len(deal.Comps) == 0 {
		return "<p>Comparable sales data not available.</p>"
	}

	var rows strings.Builder
	for _, c := range deal.Comps {
		rows.WriteString(fmt.Sprintf(
			"<tr><td>%s</td><td>%s</td><td>$%s</td><td>$%.0f</td><td>%.1f%%</td><td>%.1f mi</td></tr>\n",
			c.Address.Street,
			c.SaleDate.Format("Jan 2006"),
			formatMoney(c.SalePrice),
			c.PricePerSF,
			c.CapRate*100,
			c.Distance,
		))
	}

	// Calculate averages
	var avgCap, avgPSF float64
	for _, c := range deal.Comps {
		avgCap += c.CapRate
		avgPSF += c.PricePerSF
	}
	n := float64(len(deal.Comps))
	avgCap /= n
	avgPSF /= n

	return fmt.Sprintf(`<div class="comps">
		<table>
			<thead><tr>
				<th>Address</th><th>Sale Date</th><th>Sale Price</th>
				<th>$/SF</th><th>Cap Rate</th><th>Distance</th>
			</tr></thead>
			<tbody>%s</tbody>
			<tfoot><tr class="subtotal">
				<td colspan="3">Average</td>
				<td>$%.0f</td><td>%.1f%%</td><td>—</td>
			</tr></tfoot>
		</table>
	</div>`, rows.String(), avgPSF, avgCap*100)
}

// renderMap builds the location/map section.
func (b *Builder) renderMap(deal *domain.Deal) string {
	if deal.Location == nil {
		return "<p>Location data not available.</p>"
	}

	loc := deal.Location
	mapHTML := fmt.Sprintf(`<div class="location-map">
		<div class="map-container">
			<img src="%s" alt="Property location map" style="width:100%%;max-width:800px;border-radius:8px;" />
		</div>
		<div class="location-details">
			<p><strong>Coordinates:</strong> %.4f, %.4f</p>
			<p><strong>Address:</strong> %s</p>`,
		loc.StaticMapURL,
		loc.Latitude, loc.Longitude,
		deal.Property.Address.OneLiner(),
	)

	// Add walk/transit scores if available
	if deal.MarketData != nil {
		if deal.MarketData.WalkScore > 0 {
			mapHTML += fmt.Sprintf(`<div class="metrics" style="margin-top:16px;">
				<div class="metric"><span class="label">Walk Score</span><span class="value">%d</span></div>
				<div class="metric"><span class="label">Transit Score</span><span class="value">%d</span></div>
			</div>`, deal.MarketData.WalkScore, deal.MarketData.TransitScore)
		}
	}

	mapHTML += `</div></div>`
	return mapHTML
}

// renderDemographics builds the demographics section from market data.
func (b *Builder) renderDemographics(deal *domain.Deal) string {
	if deal.MarketData == nil {
		return "<p>Demographic data not available.</p>"
	}

	md := deal.MarketData
	return fmt.Sprintf(`<div class="demographics">
		<table>
			<tr><td>Population</td><td>%s</td></tr>
			<tr><td>Population Growth (YoY)</td><td>%.1f%%</td></tr>
			<tr><td>Median Household Income</td><td>$%s</td></tr>
			<tr><td>Unemployment Rate</td><td>%.1f%%</td></tr>
			<tr><td>Median Rent</td><td>$%s</td></tr>
			<tr><td>Submarket Vacancy Rate</td><td>%.1f%%</td></tr>
			<tr><td>Rent Growth (YoY)</td><td>%.1f%%</td></tr>
		</table>
		<div class="metrics" style="margin-top:24px;">
			<div class="metric"><span class="label">Walk Score</span><span class="value">%d</span></div>
			<div class="metric"><span class="label">Transit Score</span><span class="value">%d</span></div>
		</div>
	</div>`,
		formatMoney(float64(md.Population)),
		md.PopulationGrowth,
		formatMoney(md.MedianIncome),
		md.UnemploymentRate,
		formatMoney(md.MedianRent),
		md.VacancyRate,
		md.RentGrowth,
		md.WalkScore,
		md.TransitScore,
	)
}

// renderPhotos builds the photo gallery section.
func (b *Builder) renderPhotos(deal *domain.Deal) string {
	if len(deal.PhotoURLs) == 0 {
		return "<p>Property photos not available.</p>"
	}

	var photos strings.Builder
	for i, url := range deal.PhotoURLs {
		photos.WriteString(fmt.Sprintf(
			`<div class="photo-item"><img src="%s" alt="Property photo %d" style="width:100%%;border-radius:8px;" /></div>`,
			url, i+1,
		))
	}

	return fmt.Sprintf(`<div class="photo-gallery" style="display:grid;grid-template-columns:repeat(2,1fr);gap:16px;">%s</div>`, photos.String())
}

// renderValuation builds the valuation analysis for BOVs.
func (b *Builder) renderValuation(deal *domain.Deal) string {
	if deal.Analysis == nil {
		return "<p>Valuation analysis not available.</p>"
	}

	a := deal.Analysis

	// Derive valuation from comp cap rates if available
	var avgCapRate float64
	if len(deal.Comps) > 0 {
		for _, c := range deal.Comps {
			avgCapRate += c.CapRate
		}
		avgCapRate /= float64(len(deal.Comps))
	} else {
		avgCapRate = 0.055 // default assumption
	}

	estimatedValue := a.NOI / avgCapRate
	var pricePerUnit, pricePerSF float64
	if a.TotalUnits > 0 {
		pricePerUnit = estimatedValue / float64(a.TotalUnits)
	}
	if deal.Property.SqFt > 0 {
		pricePerSF = estimatedValue / float64(deal.Property.SqFt)
	}

	_ = time.Now() // used for date display

	return fmt.Sprintf(`<div class="valuation">
		<table>
			<tr><th colspan="2">Income Approach Valuation</th></tr>
			<tr><td>Net Operating Income</td><td>$%s</td></tr>
			<tr><td>Applied Cap Rate</td><td>%.2f%%</td></tr>
			<tr class="total"><td>Estimated Value</td><td>$%s</td></tr>
		</table>
		<div class="metrics" style="margin-top:24px;">
			<div class="metric"><span class="label">Price / Unit</span><span class="value">$%s</span></div>
			<div class="metric"><span class="label">Price / SF</span><span class="value">$%s</span></div>
		</div>
	</div>`,
		formatMoney(a.NOI),
		avgCapRate*100,
		formatMoney(estimatedValue),
		formatMoney(pricePerUnit),
		formatMoney(pricePerSF),
	)
}
