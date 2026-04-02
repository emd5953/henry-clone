package parser

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// DocumentType represents the kind of file uploaded by a broker.
// Henry handles Excel, PDF, Word — not just clean CSVs.
type DocumentType string

const (
	DocTypeCSV   DocumentType = "csv"
	DocTypeExcel DocumentType = "excel"
	DocTypePDF   DocumentType = "pdf"
	DocTypeWord  DocumentType = "word"
)

// DetectDocumentType infers the file type from its name.
func DetectDocumentType(filename string) (DocumentType, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".csv":
		return DocTypeCSV, nil
	case ".xlsx", ".xls":
		return DocTypeExcel, nil
	case ".pdf":
		return DocTypePDF, nil
	case ".docx", ".doc":
		return DocTypeWord, nil
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

// DocumentParser is the unified interface for extracting structured
// data from any broker-provided file. Each implementation handles
// the quirks of its format.
type DocumentParser interface {
	ParseRentRollFromReader(r io.Reader, docType DocumentType) (*ParsedRentRoll, error)
	ParseT12FromReader(r io.Reader, docType DocumentType) (*ParsedT12, error)
}

// ParsedRentRoll is the intermediate representation before domain conversion.
// This lets us normalize messy inputs before they hit the domain model.
type ParsedRentRoll struct {
	Units      []ParsedUnit
	Confidence float64 // 0-1, how confident we are in the extraction
	Warnings   []string
}

type ParsedUnit struct {
	UnitID      string
	Tenant      string
	SqFt        int
	MonthlyRent float64
	LeaseStart  string // raw string, parsed later
	LeaseEnd    string
}

// ParsedT12 is the intermediate T12 before domain conversion.
type ParsedT12 struct {
	Income     map[string]float64
	Expenses   map[string]float64
	Confidence float64
	Warnings   []string
}
