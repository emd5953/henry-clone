// Package export handles converting generated decks to distributable formats.
// Henry delivers decks as weblink + PDF. The PDF needs to be pixel-perfect
// because brokers send these to institutional investors.
// We use headless Chrome (chromedp) for high-fidelity HTML→PDF conversion.
package export

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// PDFExporter converts HTML deck content to PDF bytes.
type PDFExporter struct {
	timeout time.Duration
}

func NewPDFExporter() *PDFExporter {
	return &PDFExporter{timeout: 30 * time.Second}
}

// GeneratePDF takes raw HTML and returns PDF bytes.
// Uses headless Chrome for pixel-perfect rendering — same approach
// any serious document generation platform would use.
func (e *PDFExporter) GeneratePDF(ctx context.Context, html string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()

	var pdfBuf []byte
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.Sleep(500*time.Millisecond), // let CSS render
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				WithPaperWidth(8.5).
				WithPaperHeight(11).
				Do(ctx)
			if err != nil {
				return fmt.Errorf("printing to PDF: %w", err)
			}
			pdfBuf = buf
			return nil
		}),
	); err != nil {
		return nil, fmt.Errorf("chromedp: %w", err)
	}

	return pdfBuf, nil
}
