package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/henry-clone/internal/domain"
	"github.com/xuri/excelize/v2"
)

// ParseRentRollExcel reads a rent roll from an Excel file.
// Brokers love Excel — this is the most common format in practice.
// Handles messy headers via fuzzy column matching.
func ParseRentRollExcel(r io.Reader) (*domain.RentRoll, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("opening excel: %w", err)
	}
	defer f.Close()

	// Use the first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found")
	}
	sheetName := sheets[0]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("reading rows: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("need at least a header row and one data row")
	}

	// Find the header row — sometimes brokers have title rows above the actual headers.
	// We look for the first row that has at least 2 fuzzy-matched columns.
	headerIdx := -1
	var colMap map[string]int
	for i, row := range rows {
		cm := FuzzyMapColumns(row)
		if len(cm) >= 2 {
			headerIdx = i
			colMap = cm
			break
		}
	}
	if headerIdx == -1 {
		return nil, fmt.Errorf("could not identify header row")
	}

	var units []domain.UnitLine
	for i := headerIdx + 1; i < len(rows); i++ {
		row := rows[i]
		if isEmptyRow(row) {
			continue
		}

		unit := domain.UnitLine{
			UnitID: getExcelField(row, colMap, "unit_id"),
			Tenant: getExcelField(row, colMap, "tenant"),
		}

		if sqft := getExcelField(row, colMap, "sq_ft"); sqft != "" {
			unit.SqFt, _ = strconv.Atoi(cleanNumber(sqft))
		}
		if rent := getExcelField(row, colMap, "monthly_rent"); rent != "" {
			unit.MonthlyRent, _ = strconv.ParseFloat(cleanNumber(rent), 64)
		}

		units = append(units, unit)
	}

	return &domain.RentRoll{
		AsOfDate: time.Now(),
		Units:    units,
	}, nil
}

// ParseT12Excel reads a T12 from an Excel file.
// T12s in Excel are typically laid out as rows of categories with amounts.
func ParseT12Excel(r io.Reader) (*domain.T12, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("opening excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("reading rows: %w", err)
	}

	values := make(map[string]float64)
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		category := FuzzyMatchT12Category(row[0])
		if category == "" {
			continue
		}
		val, _ := strconv.ParseFloat(cleanNumber(row[len(row)-1]), 64)
		values[category] = val
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

func getExcelField(row []string, colMap map[string]int, key string) string {
	idx, ok := colMap[key]
	if !ok || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

func cleanNumber(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "$", "")
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	return s
}

func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
