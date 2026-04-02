package figma

import (
	"context"
	"fmt"

	"github.com/henry-clone/internal/domain"
)

// Bridge connects the deck generation pipeline to Figma.
// It manages the lifecycle of a deck's Figma file:
// - Creating/linking a Figma file for a deal
// - Pushing section content to Figma
// - Pulling edits back from Figma
// - Exporting the final PDF from Figma
type Bridge struct {
	client *Client
}

func NewBridge(token string) *Bridge {
	return &Bridge{client: NewClient(token)}
}

// DeckFigmaLink stores the connection between a deal and its Figma file.
type DeckFigmaLink struct {
	FileKey  string            `json:"file_key"`
	FileURL  string            `json:"file_url"`
	NodeMap  map[string]string `json:"node_map"` // sectionType -> figma nodeID
	LastSync string            `json:"last_sync"`
}

// LinkDealToFile associates a deal with an existing Figma file.
// The QC team would have a template file per deck type in Figma.
// When a deck is generated, we link it to a copy of that template.
func (b *Bridge) LinkDealToFile(deal *domain.Deal, fileKey string) *DeckFigmaLink {
	return &DeckFigmaLink{
		FileKey: fileKey,
		FileURL: fmt.Sprintf("https://www.figma.com/design/%s", fileKey),
		NodeMap: make(map[string]string),
	}
}

// GetFileStructure retrieves the Figma file's page/frame structure.
// Used to map deck sections to Figma frames.
func (b *Bridge) GetFileStructure(ctx context.Context, fileKey string) (*FileResponse, error) {
	return b.client.GetFile(ctx, fileKey)
}

// ExportPDFFromFigma exports the entire Figma file as a PDF.
// This is the "final" PDF after QC edits in Figma.
func (b *Bridge) ExportPDFFromFigma(ctx context.Context, fileKey string, pageNodeIDs []string) (map[string]string, error) {
	resp, err := b.client.ExportNodes(ctx, fileKey, pageNodeIDs, "pdf")
	if err != nil {
		return nil, fmt.Errorf("exporting PDF from Figma: %w", err)
	}
	return resp.Images, nil
}

// AddReviewComment posts a QC review comment to the Figma file.
func (b *Bridge) AddReviewComment(ctx context.Context, fileKey string, comment string) error {
	return b.client.PostComment(ctx, fileKey, comment)
}
