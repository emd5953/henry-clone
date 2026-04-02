package parser

import "strings"

// fuzzyColumnMap maps messy broker column names to canonical keys.
// Real brokers send "Unit #", "unit_number", "Unit ID", "Ste",
// "Suite", "Space #" — all meaning the same thing.
// Henry uses LLM-assisted extraction for the hardest cases;
// this handles the common ones deterministically.
var fuzzyColumnMap = map[string][]string{
	"unit_id": {
		"unit", "unit_id", "unit_number", "unit_no", "unit_num",
		"unit #", "ste", "suite", "space", "space #", "apt", "apt #",
	},
	"tenant": {
		"tenant", "tenant_name", "lessee", "occupant", "renter",
		"company", "business", "name",
	},
	"sq_ft": {
		"sq_ft", "sqft", "sf", "square_feet", "square_footage",
		"area", "size", "rsf", "rentable_sf", "nrsf",
	},
	"monthly_rent": {
		"monthly_rent", "rent", "monthly", "base_rent", "contract_rent",
		"current_rent", "in_place_rent", "scheduled_rent", "mo_rent",
	},
	"lease_start": {
		"lease_start", "start_date", "commencement", "lease_commencement",
		"move_in", "start",
	},
	"lease_end": {
		"lease_end", "end_date", "expiration", "lease_expiration",
		"lease_exp", "exp_date", "maturity",
	},
}

// FuzzyMatchColumn takes a raw header string and returns the canonical
// column name, or empty string if no match.
func FuzzyMatchColumn(raw string) string {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, ".", "")

	for canonical, variants := range fuzzyColumnMap {
		for _, v := range variants {
			if normalized == v {
				return canonical
			}
		}
	}

	// Substring matching as fallback
	for canonical, variants := range fuzzyColumnMap {
		for _, v := range variants {
			if strings.Contains(normalized, v) || strings.Contains(v, normalized) {
				return canonical
			}
		}
	}

	return ""
}

// FuzzyMapColumns takes a header row and returns a map of canonical
// column names to their index positions.
func FuzzyMapColumns(headers []string) map[string]int {
	result := make(map[string]int)
	for i, h := range headers {
		if canonical := FuzzyMatchColumn(h); canonical != "" {
			// Don't overwrite if we already found a better match
			if _, exists := result[canonical]; !exists {
				result[canonical] = i
			}
		}
	}
	return result
}

// T12 category fuzzy matching
var fuzzyT12Map = map[string][]string{
	"gross_rental_income": {
		"gross_rental_income", "rental_income", "gross_rent", "base_rental_income",
		"scheduled_rent", "gross_potential_rent", "gpr",
	},
	"other_income": {
		"other_income", "misc_income", "ancillary_income", "laundry",
		"parking", "fee_income", "utility_reimbursement",
	},
	"vacancy_loss": {
		"vacancy_loss", "vacancy", "vacancy_&_credit_loss", "credit_loss",
		"concessions", "loss_to_lease",
	},
	"taxes": {
		"taxes", "real_estate_taxes", "property_taxes", "tax",
	},
	"insurance": {
		"insurance", "property_insurance", "hazard_insurance",
	},
	"utilities": {
		"utilities", "utility", "electric", "gas", "water", "sewer",
		"water_&_sewer", "trash",
	},
	"maintenance": {
		"maintenance", "repairs", "repairs_&_maintenance", "r&m",
		"building_maintenance", "grounds",
	},
	"management": {
		"management", "management_fee", "property_management",
		"mgmt", "mgmt_fee",
	},
	"other_expense": {
		"other_expense", "other", "general_&_admin", "g&a",
		"administrative", "legal", "professional_fees", "marketing",
	},
}

// FuzzyMatchT12Category maps a raw T12 line item to a canonical category.
func FuzzyMatchT12Category(raw string) string {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	normalized = strings.ReplaceAll(normalized, "-", "_")

	for canonical, variants := range fuzzyT12Map {
		for _, v := range variants {
			if normalized == v {
				return canonical
			}
		}
	}

	for canonical, variants := range fuzzyT12Map {
		for _, v := range variants {
			if strings.Contains(normalized, v) || strings.Contains(v, normalized) {
				return canonical
			}
		}
	}

	return ""
}
