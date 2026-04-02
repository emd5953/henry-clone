package deck

import (
	"fmt"
	"strings"

	"github.com/henry-clone/internal/domain"
)

// assembleHTML wraps all sections into a complete, styled HTML document.
// Now driven by per-client brand (fonts, colors, logos) from the deal.
func (b *Builder) assembleHTML(deal *domain.Deal, sections []domain.Section) string {
	brand := deal.Brand
	if brand.ID == "" {
		brand = domain.DefaultBrand()
	}

	var body strings.Builder

	// Logo header if brand has one
	if brand.LogoURL != "" {
		body.WriteString(fmt.Sprintf(`<header class="deck-header">
			<img src="%s" alt="%s" class="brand-logo" />
		</header>`, brand.LogoURL, brand.Name))
	}

	for _, s := range sections {
		body.WriteString(fmt.Sprintf(`<section class="deck-section %s">
			<h2>%s</h2>
			<div class="section-content">%s</div>
		</section>`, s.Type, s.Title, s.Content))
		body.WriteString("\n")
	}

	css := b.brandedCSS(brand)

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s — Deal Deck</title>
<style>
%s
</style>
</head>
<body>
<div class="deck">
%s
</div>
</body>
</html>`, deal.Property.Name, css, body.String())
}

// brandedCSS generates CSS with the client's brand colors and fonts.
func (b *Builder) brandedCSS(brand domain.Brand) string {
	return fmt.Sprintf(`
* { margin: 0; padding: 0; box-sizing: border-box; }

body {
	font-family: %s;
	background: #f5f5f5;
	color: #1a1a1a;
	line-height: 1.6;
}

.deck {
	max-width: 1000px;
	margin: 0 auto;
	background: white;
	box-shadow: 0 2px 20px rgba(0,0,0,0.08);
}

.deck-header {
	padding: 24px 80px;
	border-bottom: 3px solid %s;
}

.brand-logo {
	height: 40px;
	width: auto;
}

.deck-section {
	padding: 60px 80px;
	border-bottom: 1px solid #eee;
	page-break-inside: avoid;
}

.deck-section h2 {
	font-family: %s;
	font-size: 14px;
	text-transform: uppercase;
	letter-spacing: 2px;
	color: %s;
	margin-bottom: 24px;
}

.cover {
	text-align: center;
	padding: 40px 0;
}

.cover h1 {
	font-family: %s;
	font-size: 42px;
	font-weight: 700;
	margin-bottom: 12px;
	color: %s;
}

.cover .address {
	font-size: 18px;
	color: %s;
	margin-bottom: 8px;
}

.cover .asset-class {
	font-size: 14px;
	text-transform: uppercase;
	letter-spacing: 1.5px;
	color: #999;
}

.section-content {
	font-size: 16px;
	color: #333;
}

.section-content p {
	margin-bottom: 16px;
}

table {
	width: 100%%;
	border-collapse: collapse;
	margin: 16px 0;
}

th, td {
	padding: 10px 16px;
	text-align: left;
	border-bottom: 1px solid #eee;
	font-size: 14px;
}

th {
	font-weight: 600;
	color: %s;
	background: #fafafa;
	text-transform: uppercase;
	font-size: 11px;
	letter-spacing: 1px;
}

td:last-child {
	text-align: right;
}

.subtotal td {
	font-weight: 600;
	border-top: 2px solid #ddd;
}

.total td {
	font-weight: 700;
	font-size: 16px;
	border-top: 3px solid %s;
	padding-top: 16px;
}

.metrics {
	display: flex;
	gap: 40px;
	margin-top: 32px;
	padding-top: 24px;
	border-top: 1px solid #eee;
}

.metric {
	display: flex;
	flex-direction: column;
}

.metric .label {
	font-size: 11px;
	text-transform: uppercase;
	letter-spacing: 1px;
	color: #888;
}

.metric .value {
	font-size: 28px;
	font-weight: 700;
	color: %s;
}

.vacant {
	color: #e74c3c;
	font-style: italic;
}

.rent-roll table {
	font-size: 13px;
}

a { color: %s; }

@media print {
	body { background: white; }
	.deck { box-shadow: none; }
	.deck-section { padding: 40px 60px; }
}
`,
		brand.FontBody,
		brand.PrimaryColor,
		brand.FontHeading,
		brand.SecondaryColor,
		brand.FontHeading,
		brand.PrimaryColor,
		brand.SecondaryColor,
		brand.PrimaryColor,
		brand.PrimaryColor,
		brand.PrimaryColor,
		brand.AccentColor,
	)
}
