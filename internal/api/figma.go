package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/henry-clone/internal/figma"
)

// FigmaHandler manages Figma integration endpoints.
type FigmaHandler struct {
	bridge  *figma.Bridge
	handler *Handler // reference back to deal store
}

func NewFigmaHandler(bridge *figma.Bridge, h *Handler) *FigmaHandler {
	return &FigmaHandler{bridge: bridge, handler: h}
}

// LinkFigmaFile handles POST /api/deals/{dealID}/figma/link
// Associates a Figma file with a deal for QC editing.
type LinkFigmaRequest struct {
	FileKey string `json:"file_key"` // from the Figma URL
}

func (fh *FigmaHandler) LinkFigmaFile(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	var req LinkFigmaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	fh.handler.mu.Lock()
	deal, ok := fh.handler.deals[dealID]
	if !ok {
		fh.handler.mu.Unlock()
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	link := fh.bridge.LinkDealToFile(deal, req.FileKey)
	deal.FigmaFileKey = link.FileKey
	deal.FigmaFileURL = link.FileURL
	fh.handler.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

// GetFigmaFile handles GET /api/deals/{dealID}/figma
// Returns the Figma file structure for the linked file.
func (fh *FigmaHandler) GetFigmaFile(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	fh.handler.mu.RLock()
	deal, ok := fh.handler.deals[dealID]
	fh.handler.mu.RUnlock()

	if !ok {
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.FigmaFileKey == "" {
		http.Error(w, "no Figma file linked", http.StatusNotFound)
		return
	}

	file, err := fh.bridge.GetFileStructure(r.Context(), deal.FigmaFileKey)
	if err != nil {
		http.Error(w, "Figma API error: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}

// ExportFigmaPDF handles GET /api/deals/{dealID}/figma/export
// Exports the Figma file as PDF — this is the final deliverable after QC.
func (fh *FigmaHandler) ExportFigmaPDF(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	fh.handler.mu.RLock()
	deal, ok := fh.handler.deals[dealID]
	fh.handler.mu.RUnlock()

	if !ok {
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.FigmaFileKey == "" {
		http.Error(w, "no Figma file linked", http.StatusNotFound)
		return
	}

	// Get the file structure to find page node IDs
	file, err := fh.bridge.GetFileStructure(r.Context(), deal.FigmaFileKey)
	if err != nil {
		http.Error(w, "Figma API error: "+err.Error(), http.StatusBadGateway)
		return
	}

	var pageIDs []string
	for _, page := range file.Document.Children {
		pageIDs = append(pageIDs, page.ID)
	}

	images, err := fh.bridge.ExportPDFFromFigma(r.Context(), deal.FigmaFileKey, pageIDs)
	if err != nil {
		http.Error(w, "PDF export failed: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

// PostFigmaComment handles POST /api/deals/{dealID}/figma/comment
// Adds a QC review comment to the Figma file.
type FigmaCommentRequest struct {
	Message string `json:"message"`
}

func (fh *FigmaHandler) PostFigmaComment(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	var req FigmaCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	fh.handler.mu.RLock()
	deal, ok := fh.handler.deals[dealID]
	fh.handler.mu.RUnlock()

	if !ok {
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.FigmaFileKey == "" {
		http.Error(w, "no Figma file linked", http.StatusNotFound)
		return
	}

	if err := fh.bridge.AddReviewComment(r.Context(), deal.FigmaFileKey, req.Message); err != nil {
		http.Error(w, "comment failed: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "posted"})
}
