package domain

import "time"

// ReviewStatus tracks where a deck is in the QC pipeline.
// Henry has a human QC team that spends ~15 min reviewing every deck
// before it's delivered to the broker. This models that workflow.
type ReviewStatus string

const (
	ReviewPending  ReviewStatus = "pending_review"
	ReviewInProgress ReviewStatus = "in_review"
	ReviewApproved ReviewStatus = "approved"
	ReviewRejected ReviewStatus = "needs_revision"
)

// Review represents a QC review of a generated deck.
type Review struct {
	ID         string       `json:"id"`
	DealID     string       `json:"deal_id"`
	ReviewerID string       `json:"reviewer_id"`
	Status     ReviewStatus `json:"status"`
	StartedAt  *time.Time   `json:"started_at,omitempty"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
	Notes      string       `json:"notes,omitempty"`
	Edits      []SectionEdit `json:"edits,omitempty"`
}

// SectionEdit records a change made by a reviewer to a specific section.
type SectionEdit struct {
	SectionIdx  int       `json:"section_idx"`
	Field       string    `json:"field"` // "title" or "content"
	OldValue    string    `json:"old_value"`
	NewValue    string    `json:"new_value"`
	EditedAt    time.Time `json:"edited_at"`
	EditedBy    string    `json:"edited_by"`
}
