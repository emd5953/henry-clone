//go:build ignore

// Quick script to generate sample Excel files for testing.
// Run: go run scripts/create_sample_excel.go
package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func main() {
	// Rent Roll with messy headers (like a real broker would send)
	f := excelize.NewFile()
	sheet := "Sheet1"

	// Intentionally messy headers — tests fuzzy matching
	headers := []string{"Unit #", "Tenant Name", "SF", "Base Rent"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	data := [][]interface{}{
		{"101", "Acme Corp", 1200, 2850},
		{"102", "Bay Area Dental", 950, 2200},
		{"103", "Summit Legal Group", 1400, 3100},
		{"104", "", 1100, 0},
		{"105", "Redwood Financial", 1350, 3000},
		{"106", "Pacific Design Studio", 800, 1900},
		{"107", "Golden Gate Consulting", 1500, 3400},
		{"108", "", 900, 0},
		{"109", "Marina Wellness Center", 1100, 2600},
		{"110", "Coastal Property Mgmt", 1000, 2350},
	}

	for i, row := range data {
		for j, val := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, val)
		}
	}

	if err := f.SaveAs("sample_data/rent_roll.xlsx"); err != nil {
		panic(err)
	}
	fmt.Println("Created sample_data/rent_roll.xlsx")

	// T12
	f2 := excelize.NewFile()
	t12Data := [][]interface{}{
		{"Category", "Annual Amount"},
		{"Gross Rental Income", 295200},
		{"Other Income", 18000},
		{"Vacancy Loss", 24600},
		{"Real Estate Taxes", 42000},
		{"Insurance", 15600},
		{"Utilities", 22800},
		{"Repairs & Maintenance", 18000},
		{"Management Fee", 14400},
		{"General & Admin", 9600},
	}

	for i, row := range t12Data {
		for j, val := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
			f2.SetCellValue("Sheet1", cell, val)
		}
	}

	if err := f2.SaveAs("sample_data/t12.xlsx"); err != nil {
		panic(err)
	}
	fmt.Println("Created sample_data/t12.xlsx")
}
