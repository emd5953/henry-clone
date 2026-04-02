package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/henry-clone/internal/agent"
	"github.com/henry-clone/internal/deck"
	"github.com/henry-clone/internal/domain"
	"github.com/henry-clone/internal/enrichment"
	"github.com/henry-clone/internal/export"
	"github.com/henry-clone/internal/parser"
)

// Handler serves the HTTP API.
// Now supports the multi-agent pipeline, deck types, and branding.
type Handler struct {
	builder       *deck.Builder
	narrator      deck.Narrator
	comps         enrichment.CompsProvider
	market        enrichment.MarketDataProvider
	geo           enrichment.GeoProvider
	pdfExporter   *export.PDFExporter
	deals         map[string]*domain.Deal
	mu            sync.RWMutex
}

type HandlerConfig struct {
	Builder  *deck.Builder
	Narrator deck.Narrator
	Comps    enrichment.CompsProvider
	Market   enrichment.MarketDataProvider
	Geo      enrichment.GeoProvider
}

func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		builder:     cfg.Builder,
		narrator:    cfg.Narrator,
		comps:       cfg.Comps,
		market:      cfg.Market,
		geo:         cfg.Geo,
		pdfExporter: export.NewPDFExporter(),
		deals:       make(map[string]*domain.Deal),
	}
}

// CreateDealRequest accepts either JSON or multipart form data.
type CreateDealRequest struct {
	Property domain.Property `json:"property"`
	RentRoll domain.RentRoll `json:"rent_roll"`
	T12      domain.T12      `json:"t12"`
	Thesis   string          `json:"thesis"`
	DeckType domain.DeckType `json:"deck_type"`
	Brand    domain.Brand    `json:"brand"`
}

// CreateDeal handles POST /api/deals
// Now runs the full multi-agent pipeline instead of sequential builder.
func (h *Handler) CreateDeal(w http.ResponseWriter, r *http.Request) {
	var deal domain.Deal

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := h.parseDealFromForm(r, &deal); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		var req CreateDealRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		deal.Property = req.Property
		deal.RentRoll = req.RentRoll
		deal.T12 = req.T12
		deal.Thesis = req.Thesis
		deal.DeckType = req.DeckType
		deal.Brand = req.Brand
	}

	// Defaults
	if deal.DeckType == "" {
		deal.DeckType = domain.DeckTypeOM
	}
	if deal.Brand.ID == "" {
		deal.Brand = domain.DefaultBrand()
	}

	deal.ID = generateID()
	deal.CreatedAt = time.Now()
	deal.Status = domain.StatusPending

	// Run the multi-agent pipeline
	pipeline := agent.NewPipeline(
		agent.DataExtractionAgent(),
		agent.FinancialAnalysisAgent(),
		agent.CompsAgent(h.comps),
		agent.MarketDataAgent(h.market),
		agent.GeoAgent(h.geo),
		agent.NarrativeAgent(h.narrator),
		agent.AssemblyAgent(h.builder),
	)

	state := agent.NewPipelineState()
	state.Set(agent.KeyDeal, &deal)

	deal.Status = domain.StatusAnalyzing
	results := pipeline.Run(r.Context(), state)

	// Check for critical failures
	var failed []string
	for _, res := range results {
		if res.Status == agent.StatusFailed {
			failed = append(failed, res.AgentName)
			log.Printf("agent %s failed after %d attempts: %v", res.AgentName, res.Attempts, res.Error)
		}
	}

	// Assembly failure is fatal; enrichment failures are tolerable
	for _, f := range failed {
		if f == "deck_assembly" || f == "narrative_generation" || f == "financial_analysis" {
			deal.Status = domain.StatusFailed
			http.Error(w, "deck generation failed: "+f, http.StatusInternalServerError)
			return
		}
	}

	h.mu.Lock()
	h.deals[deal.ID] = &deal
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deal)
}

// GetDeal handles GET /api/deals/{dealID}
func (h *Handler) GetDeal(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	h.mu.RLock()
	deal, ok := h.deals[dealID]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deal)
}

// GetDeck handles GET /api/deals/{dealID}/deck — returns the HTML deck.
func (h *Handler) GetDeck(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	h.mu.RLock()
	deal, ok := h.deals[dealID]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.Deck == nil {
		http.Error(w, "deck not yet generated", http.StatusAccepted)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(deal.Deck.HTML))
}

// ListDeals handles GET /api/deals
func (h *Handler) ListDeals(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	deals := make([]*domain.Deal, 0, len(h.deals))
	for _, d := range h.deals {
		deals = append(deals, d)
	}
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deals)
}

func (h *Handler) parseDealFromForm(r *http.Request, deal *domain.Deal) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return err
	}

	deal.Property.Name = r.FormValue("property_name")
	deal.Property.Address = domain.Address{
		Street: r.FormValue("street"),
		City:   r.FormValue("city"),
		State:  r.FormValue("state"),
		Zip:    r.FormValue("zip"),
	}
	deal.Property.AssetClass = domain.AssetClass(r.FormValue("asset_class"))
	deal.Thesis = r.FormValue("thesis")

	if dt := r.FormValue("deck_type"); dt != "" {
		deal.DeckType = domain.DeckType(dt)
	}

	if units := r.FormValue("units"); units != "" {
		deal.Property.Units, _ = strconv.Atoi(units)
	}
	if sqft := r.FormValue("sq_ft"); sqft != "" {
		deal.Property.SqFt, _ = strconv.Atoi(sqft)
	}
	if yb := r.FormValue("year_built"); yb != "" {
		deal.Property.YearBuilt, _ = strconv.Atoi(yb)
	}

	// Parse rent roll — CSV or Excel
	if file, header, err := r.FormFile("rent_roll"); err == nil {
		defer file.Close()
		docType, _ := parser.DetectDocumentType(header.Filename)
		switch docType {
		case parser.DocTypeExcel:
			rr, err := parser.ParseRentRollExcel(file)
			if err != nil {
				return err
			}
			deal.RentRoll = *rr
		default: // CSV fallback
			rr, err := parser.ParseRentRoll(file)
			if err != nil {
				return err
			}
			deal.RentRoll = *rr
		}
		if deal.Property.Units == 0 {
			deal.Property.Units = len(deal.RentRoll.Units)
		}
	}

	// Parse T12 — CSV or Excel
	if file, header, err := r.FormFile("t12"); err == nil {
		defer file.Close()
		docType, _ := parser.DetectDocumentType(header.Filename)
		switch docType {
		case parser.DocTypeExcel:
			t12, err := parser.ParseT12Excel(file)
			if err != nil {
				return err
			}
			deal.T12 = *t12
		default:
			t12, err := parser.ParseT12(file)
			if err != nil {
				return err
			}
			deal.T12 = *t12
		}
	}

	return nil
}

// generateID creates a crypto-random ID instead of a predictable timestamp.
func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GetDeckPDF handles GET /api/deals/{dealID}/deck.pdf
// Converts the HTML deck to PDF using headless Chrome.
// Henry delivers as weblink + PDF — this is the PDF path.
func (h *Handler) GetDeckPDF(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	h.mu.RLock()
	deal, ok := h.deals[dealID]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.Deck == nil {
		http.Error(w, "deck not yet generated", http.StatusAccepted)
		return
	}

	pdfBytes, err := h.pdfExporter.GeneratePDF(r.Context(), deal.Deck.HTML)
	if err != nil {
		log.Printf("PDF generation failed for deal %s: %v", dealID, err)
		http.Error(w, "PDF generation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filename := strings.ReplaceAll(deal.Property.Name, " ", "_") + "_Deck.pdf"
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	w.Write(pdfBytes)
}

// UpdateSection handles PUT /api/deals/{dealID}/sections/{sectionIdx}
// This powers the deck editor — brokers can tweak individual sections.
// Henry's React editor lets users adjust content without losing AI benefits.
type UpdateSectionRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *Handler) UpdateSection(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")
	sectionIdxStr := chi.URLParam(r, "sectionIdx")
	sectionIdx, err := strconv.Atoi(sectionIdxStr)
	if err != nil {
		http.Error(w, "invalid section index", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	deal, ok := h.deals[dealID]
	if !ok {
		h.mu.Unlock()
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.Deck == nil || sectionIdx < 0 || sectionIdx >= len(deal.Deck.Sections) {
		h.mu.Unlock()
		http.Error(w, "invalid section", http.StatusBadRequest)
		return
	}

	var req UpdateSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.mu.Unlock()
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title != "" {
		deal.Deck.Sections[sectionIdx].Title = req.Title
	}
	if req.Content != "" {
		deal.Deck.Sections[sectionIdx].Content = req.Content
	}

	// Rebuild HTML with updated sections
	deal.Deck.HTML = h.builder.RebuildHTML(deal)
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deal.Deck.Sections[sectionIdx])
}

// GetSections handles GET /api/deals/{dealID}/sections
// Returns the structured section data for the deck editor.
func (h *Handler) GetSections(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	h.mu.RLock()
	deal, ok := h.deals[dealID]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.Deck == nil {
		http.Error(w, "deck not yet generated", http.StatusAccepted)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deal.Deck.Sections)
}
