// Package domain contains the core business types for CRE deal modeling.
// These types are deliberately decoupled from any API or storage concerns.
// The goal: model the messy real world of CRE into clean, composable abstractions.
package domain

import "time"

// Property represents a commercial real estate asset.
// In the real world, properties are wildly inconsistent — a 200-unit multifamily
// looks nothing like a single-tenant NNN retail deal. We keep this flexible
// by using AssetClass to drive downstream behavior (analysis, deck layout, etc.)
type Property struct {
	Name       string     `json:"name"`
	Address    Address    `json:"address"`
	AssetClass AssetClass `json:"asset_class"`
	Units      int        `json:"units,omitempty"`
	SqFt       int        `json:"sq_ft,omitempty"`
	YearBuilt  int        `json:"year_built,omitempty"`
	ImageURLs  []string   `json:"image_urls,omitempty"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

func (a Address) OneLiner() string {
	return a.Street + ", " + a.City + ", " + a.State + " " + a.Zip
}

type AssetClass string

const (
	Multifamily AssetClass = "multifamily"
	Office      AssetClass = "office"
	Retail      AssetClass = "retail"
	Industrial  AssetClass = "industrial"
	MixedUse    AssetClass = "mixed_use"
)

// RentRoll is a snapshot of all tenants/units at a point in time.
// This is one of the messiest inputs — brokers send everything from
// polished Excel files to hand-typed CSVs with inconsistent columns.
type RentRoll struct {
	AsOfDate time.Time  `json:"as_of_date"`
	Units    []UnitLine `json:"units"`
}

type UnitLine struct {
	UnitID     string  `json:"unit_id"`
	Tenant     string  `json:"tenant,omitempty"` // empty = vacant
	SqFt       int     `json:"sq_ft,omitempty"`
	MonthlyRent float64 `json:"monthly_rent"`
	LeaseStart *time.Time `json:"lease_start,omitempty"`
	LeaseEnd   *time.Time `json:"lease_end,omitempty"`
}

func (u UnitLine) IsVacant() bool {
	return u.Tenant == "" && u.MonthlyRent == 0
}

func (u UnitLine) AnnualRent() float64 {
	return u.MonthlyRent * 12
}

// T12 represents trailing 12 months of income and expenses.
// This is the financial heartbeat of any CRE deal — it tells you
// what the property actually earned vs. what it cost to operate.
type T12 struct {
	PeriodEnd time.Time    `json:"period_end"`
	Income    IncomeItems  `json:"income"`
	Expenses  ExpenseItems `json:"expenses"`
}

type IncomeItems struct {
	GrossRentalIncome float64 `json:"gross_rental_income"`
	OtherIncome       float64 `json:"other_income"`
	VacancyLoss       float64 `json:"vacancy_loss"` // stored as positive number
}

func (i IncomeItems) EffectiveGrossIncome() float64 {
	return i.GrossRentalIncome + i.OtherIncome - i.VacancyLoss
}

type ExpenseItems struct {
	Taxes       float64 `json:"taxes"`
	Insurance   float64 `json:"insurance"`
	Utilities   float64 `json:"utilities"`
	Maintenance float64 `json:"maintenance"`
	Management  float64 `json:"management"`
	Other       float64 `json:"other"`
}

func (e ExpenseItems) Total() float64 {
	return e.Taxes + e.Insurance + e.Utilities + e.Maintenance + e.Management + e.Other
}
