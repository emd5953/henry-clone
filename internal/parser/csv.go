// Package parser handles ingestion of messy broker-provided data files.
// In the real world, this is one of the hardest parts of the system.
// Brokers send wildly inconsistent formats — different column names,
// merged cells, extra header rows, mixed date formats, etc.
// This implementation handles clean CSVs; a production system would
// need fuzzy column matching, LLM-assisted extraction, and validation.
package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/henry-clone/internal/domain"
)

// ParseRentRoll reads a CSV with columns: unit_id, tenant, sq_ft, monthly_rent
func ParseRentRoll(r io.Reader) (*domain.RentRoll, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	// Read header row
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	colMap := mapColumns(header)
	required := []string{"unit_id", "monthly_rent"}
	for _, col := range required {
		if _, ok := colMap[col]; !ok {
			return nil, fmt.Errorf("missing required column: %s", col)
		}
	}

	var units []domain.UnitLine
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading row: %w", err)
		}

		unit := domain.UnitLine{
			UnitID: getField(record, colMap, "unit_id"),
			Tenant: getField(record, colMap, "tenant"),
		}

		if sqft := getField(record, colMap, "sq_ft"); sqft != "" {
			unit.SqFt, _ = strconv.Atoi(sqft)
		}

		if rent := getField(record, colMap, "monthly_rent"); rent != "" {
			unit.MonthlyRent, _ = strconv.ParseFloat(strings.ReplaceAll(rent, ",", ""), 64)
		}

		units = append(units, unit)
	}

	return &domain.RentRoll{
		AsOfDate: time.Now(),
		Units:    units,
	}, nil
}

// ParseT12 reads a CSV with rows: category, amount
// Expected categories: gross_rental_income, other_income, vacancy_loss,
// taxes, insurance, utilities, maintenance, management, other_expense
func ParseT12(r io.Reader) (*domain.T12, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	values := make(map[string]float64)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading row: %w", err)
		}
		if len(record) < 2 {
			continue
		}

		key := strings.TrimSpace(strings.ToLower(record[0]))
		val, _ := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(record[1]), ",", ""), 64)
		values[key] = val
	}

	return &domain.T12{
		PeriodEnd: time.Now(),
		Income: domain.IncomeItems{
			GrossRentalIncome: values["gross_rental_income"],
			OtherIncome:       values["other_income"],
			VacancyLoss:       values["vacancy_loss"],
		},
		Expenses: domain.ExpenseItems{
			Taxes:       values["taxes"],
			Insurance:   values["insurance"],
			Utilities:   values["utilities"],
			Maintenance: values["maintenance"],
			Management:  values["management"],
			Other:       values["other_expense"],
		},
	}, nil
}

// mapColumns normalizes header names to lowercase snake_case keys.
// This is a simple version — production would use fuzzy matching
// to handle "Unit #", "unit_number", "Unit ID", etc.
func mapColumns(header []string) map[string]int {
	m := make(map[string]int)
	for i, h := range header {
		key := strings.TrimSpace(strings.ToLower(h))
		key = strings.ReplaceAll(key, " ", "_")
		m[key] = i
	}
	return m
}

func getField(record []string, colMap map[string]int, key string) string {
	idx, ok := colMap[key]
	if !ok || idx >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[idx])
}
