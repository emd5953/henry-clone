package domain

// Brand represents a brokerage's visual identity.
// Henry learns each client's fonts, colors, logos, and tone.
// Every deck is rendered against a brand so output feels bespoke.
type Brand struct {
	ID             string `json:"id"`
	Name           string `json:"name"` // e.g. "Colliers", "CBRE"
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	AccentColor    string `json:"accent_color"`
	FontHeading    string `json:"font_heading"`
	FontBody       string `json:"font_body"`
	LogoURL        string `json:"logo_url,omitempty"`
	Tone           string `json:"tone,omitempty"` // e.g. "formal", "conversational"
}

// DefaultBrand is used when no client brand is configured.
func DefaultBrand() Brand {
	return Brand{
		ID:             "default",
		Name:           "Default",
		PrimaryColor:   "#111111",
		SecondaryColor: "#555555",
		AccentColor:    "#2563eb",
		FontHeading:    "'Helvetica Neue', Arial, sans-serif",
		FontBody:       "'Helvetica Neue', Arial, sans-serif",
	}
}

// DeckType determines the structure and content of the generated deck.
// Henry produces OMs, BOVs, loan packages, syndication decks, and flyers.
type DeckType string

const (
	DeckTypeOM          DeckType = "offering_memorandum"
	DeckTypeBOV         DeckType = "broker_opinion_of_value"
	DeckTypeLoanPackage DeckType = "loan_package"
	DeckTypeSyndication DeckType = "syndication_deck"
	DeckTypeFlyer       DeckType = "leasing_flyer"
	DeckTypeTeaser      DeckType = "investment_teaser"
)
