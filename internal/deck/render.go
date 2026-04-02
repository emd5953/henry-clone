package deck

import (
	"fmt"
	"strings"

	"github.com/henry-clone/internal/domain"
)

// renderCover builds the cover page HTML.
func (b *Builder) renderCover(deal *domain.Deal) string {
	ac := string(deal.Property.AssetClass)
	if len(ac) > 0 {
		ac = strings.ToUpper(ac[:1]) + ac[1:]
	}
	return fmt.Sprintf(`<div class="cover">
		<h1>%s</h1>
		<p class="address">%s</p>
		<p class="asset-class">%s | %d Units</p>
	</div>`, deal.Property.Name, deal.Property.Address.OneLiner(), ac, deal.Property.Units)
}

// renderFinancials builds the financial summary table.
// This is deterministic — no LLM needed. Pure data → HTML.
func (b *Builder) renderFinancials(deal *domain.Deal) string {
	a := deal.Analysis
	if a == nil {
		return "<p>Financial analysis not available.</p>"
	}

	return fmt.Sprintf(`<div class="financials">
		<table>
			<tr><th colspan="2">Income</th></tr>
			<tr><td>Gross Rental Income</td><td>$%s</td></tr>
			<tr><td>Other Income</td><td>$%s</td></tr>
			<tr><td>Vacancy Loss</td><td>($%s)</td></tr>
			<tr class="subtotal"><td>Effective Gross Income</td><td>$%s</td></tr>
			<tr><th colspan="2">Expenses</th></tr>
			<tr><td>Taxes</td><td>$%s</td></tr>
			<tr><td>Insurance</td><td>$%s</td></tr>
			<tr><td>Utilities</td><td>$%s</td></tr>
			<tr><td>Maintenance</td><td>$%s</td></tr>
			<tr><td>Management</td><td>$%s</td></tr>
			<tr><td>Other</td><td>$%s</td></tr>
			<tr class="subtotal"><td>Total Expenses</td><td>$%s</td></tr>
			<tr class="total"><td>Net Operating Income (NOI)</td><td>$%s</td></tr>
		</table>
		<div class="metrics">
			<div class="metric"><span class="label">Occupancy</span><span class="value">%.1f%%</span></div>
			<div class="metric"><span class="label">Expense Ratio</span><span class="value">%.1f%%</span></div>
			<div class="metric"><span class="label">Avg Monthly Rent</span><span class="value">$%s</span></div>
		</div>
	</div>`,
		formatMoney(deal.T12.Income.GrossRentalIncome),
		formatMoney(deal.T12.Income.OtherIncome),
		formatMoney(deal.T12.Income.VacancyLoss),
		formatMoney(a.EffectiveGrossIncome),
		formatMoney(deal.T12.Expenses.Taxes),
		formatMoney(deal.T12.Expenses.Insurance),
		formatMoney(deal.T12.Expenses.Utilities),
		formatMoney(deal.T12.Expenses.Maintenance),
		formatMoney(deal.T12.Expenses.Management),
		formatMoney(deal.T12.Expenses.Other),
		formatMoney(a.TotalExpenses),
		formatMoney(a.NOI),
		a.OccupancyRate*100,
		a.ExpenseRatio*100,
		formatMoney(a.AvgMonthlyRent),
	)
}

// renderRentRoll builds the rent roll table.
func (b *Builder) renderRentRoll(deal *domain.Deal) string {
	var rows strings.Builder
	for _, u := range deal.RentRoll.Units {
		tenant := u.Tenant
		if u.IsVacant() {
			tenant = `<span class="vacant">Vacant</span>`
		}
		rows.WriteString(fmt.Sprintf(
			"<tr><td>%s</td><td>%s</td><td>%d</td><td>$%s</td></tr>\n",
			u.UnitID, tenant, u.SqFt, formatMoney(u.MonthlyRent),
		))
	}

	return fmt.Sprintf(`<div class="rent-roll">
		<table>
			<thead><tr><th>Unit</th><th>Tenant</th><th>Sq Ft</th><th>Monthly Rent</th></tr></thead>
			<tbody>%s</tbody>
		</table>
	</div>`, rows.String())
}

func formatMoney(amount float64) string {
	if amount >= 1_000_000 {
		return fmt.Sprintf("%.2fM", amount/1_000_000)
	}
	// Go doesn't have %,d — format with commas manually
	raw := fmt.Sprintf("%.0f", amount)
	if len(raw) <= 3 {
		return raw
	}
	var result []byte
	for i, c := range raw {
		if i > 0 && (len(raw)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
