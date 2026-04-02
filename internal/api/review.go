package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/henry-clone/internal/domain"
)

// GetReviewQueue handles GET /api/reviews
// Returns all deals that need QC review (status = ready but not yet approved).
func (h *Handler) GetReviewQueue(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	var queue []*domain.Deal
	for _, d := range h.deals {
		if d.Status == domain.StatusReady || d.Status == domain.StatusInReview {
			queue = append(queue, d)
		}
	}
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queue)
}

// StartReview handles POST /api/deals/{dealID}/review/start
// A reviewer claims a deck for QC.
type StartReviewRequest struct {
	ReviewerID string `json:"reviewer_id"`
}

func (h *Handler) StartReview(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	var req StartReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	deal, ok := h.deals[dealID]
	if !ok {
		h.mu.Unlock()
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.Status != domain.StatusReady && deal.Status != domain.StatusInReview {
		h.mu.Unlock()
		http.Error(w, "deal is not ready for review", http.StatusBadRequest)
		return
	}

	now := time.Now()
	deal.Status = domain.StatusInReview
	deal.Review = &domain.Review{
		ID:         generateID(),
		DealID:     dealID,
		ReviewerID: req.ReviewerID,
		Status:     domain.ReviewInProgress,
		StartedAt:  &now,
	}
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deal.Review)
}

// CompleteReview handles POST /api/deals/{dealID}/review/complete
type CompleteReviewRequest struct {
	Status string `json:"status"` // "approved" or "needs_revision"
	Notes  string `json:"notes"`
}

func (h *Handler) CompleteReview(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	var req CompleteReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	deal, ok := h.deals[dealID]
	if !ok {
		h.mu.Unlock()
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.Review == nil {
		h.mu.Unlock()
		http.Error(w, "no active review", http.StatusBadRequest)
		return
	}

	now := time.Now()
	deal.Review.CompletedAt = &now
	deal.Review.Notes = req.Notes

	switch req.Status {
	case "approved":
		deal.Review.Status = domain.ReviewApproved
		deal.Status = domain.StatusApproved
	case "needs_revision":
		deal.Review.Status = domain.ReviewRejected
		deal.Status = domain.StatusReady // goes back to queue
	default:
		h.mu.Unlock()
		http.Error(w, "status must be 'approved' or 'needs_revision'", http.StatusBadRequest)
		return
	}
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deal.Review)
}

// ReviewEdit handles POST /api/deals/{dealID}/review/edit
// Reviewer edits a section during QC — tracked as an audit trail.
type ReviewEditRequest struct {
	SectionIdx int    `json:"section_idx"`
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
}

func (h *Handler) ReviewEdit(w http.ResponseWriter, r *http.Request) {
	dealID := chi.URLParam(r, "dealID")

	var req ReviewEditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	deal, ok := h.deals[dealID]
	if !ok {
		h.mu.Unlock()
		http.Error(w, "deal not found", http.StatusNotFound)
		return
	}

	if deal.Review == nil || deal.Review.Status != domain.ReviewInProgress {
		h.mu.Unlock()
		http.Error(w, "no active review in progress", http.StatusBadRequest)
		return
	}

	if deal.Deck == nil || req.SectionIdx < 0 || req.SectionIdx >= len(deal.Deck.Sections) {
		h.mu.Unlock()
		http.Error(w, "invalid section index", http.StatusBadRequest)
		return
	}

	section := &deal.Deck.Sections[req.SectionIdx]

	if req.Title != "" {
		deal.Review.Edits = append(deal.Review.Edits, domain.SectionEdit{
			SectionIdx: req.SectionIdx,
			Field:      "title",
			OldValue:   section.Title,
			NewValue:   req.Title,
			EditedAt:   time.Now(),
			EditedBy:   deal.Review.ReviewerID,
		})
		section.Title = req.Title
	}

	if req.Content != "" {
		deal.Review.Edits = append(deal.Review.Edits, domain.SectionEdit{
			SectionIdx: req.SectionIdx,
			Field:      "content",
			OldValue:   section.Content,
			NewValue:   req.Content,
			EditedAt:   time.Now(),
			EditedBy:   deal.Review.ReviewerID,
		})
		section.Content = req.Content
	}

	// Rebuild HTML with edits
	deal.Deck.HTML = h.builder.RebuildHTML(deal)
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(section)
}
