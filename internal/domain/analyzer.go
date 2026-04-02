package domain

// Analyze computes the financial metrics from raw deal inputs.
// This is deterministic — no LLM, no external calls. Pure math.
// Keeping this in the domain package because it's core business logic,
// not infrastructure.
func Analyze(deal *Deal) *FinancialAnalysis {
	t12 := deal.T12
	rr := deal.RentRoll

	egi := t12.Income.EffectiveGrossIncome()
	totalExp := t12.Expenses.Total()
	noi := egi - totalExp

	totalUnits := len(rr.Units)
	occupied := 0
	var totalMonthlyRent float64
	for _, u := range rr.Units {
		if !u.IsVacant() {
			occupied++
			totalMonthlyRent += u.MonthlyRent
		}
	}

	var occupancyRate float64
	if totalUnits > 0 {
		occupancyRate = float64(occupied) / float64(totalUnits)
	}

	var avgRent float64
	if occupied > 0 {
		avgRent = totalMonthlyRent / float64(occupied)
	}

	var expenseRatio float64
	if egi > 0 {
		expenseRatio = totalExp / egi
	}

	analysis := &FinancialAnalysis{
		EffectiveGrossIncome: egi,
		TotalExpenses:        totalExp,
		NOI:                  noi,
		TotalUnits:           totalUnits,
		OccupiedUnits:        occupied,
		OccupancyRate:        occupancyRate,
		ExpenseRatio:         expenseRatio,
		AvgMonthlyRent:       avgRent,
		TotalAnnualRent:      totalMonthlyRent * 12,
	}

	// Per-unit metrics for multifamily, per-sqft for commercial
	switch deal.Property.AssetClass {
	case Multifamily:
		if totalUnits > 0 {
			analysis.NOIPerUnit = noi / float64(totalUnits)
		}
	default:
		if deal.Property.SqFt > 0 {
			analysis.NOIPerSqFt = noi / float64(deal.Property.SqFt)
		}
	}

	return analysis
}
