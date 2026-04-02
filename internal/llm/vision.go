package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/genai"
)

// PropertyAesthetic is the visual identity extracted from property photos.
type PropertyAesthetic struct {
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	AccentColor    string `json:"accent_color"`
	BackgroundTone string `json:"background_tone"`
	Style          string `json:"style"`
	Mood           string `json:"mood"`
	FontSuggestion string `json:"font_suggestion"`
	Description    string `json:"description"`
}

// AnalyzePropertyPhotos sends photos to Gemini Vision and extracts
// a design aesthetic that drives the deck's visual identity.
func AnalyzePropertyPhotos(ctx context.Context, client *genai.Client, photoPaths []string) (*PropertyAesthetic, error) {
	if len(photoPaths) == 0 {
		return defaultAesthetic(), nil
	}

	prompt := `Analyze these commercial real estate property photos and extract a visual design aesthetic for a professional offering memorandum.

Return a JSON object with these fields:
- primary_color: hex color of the dominant architectural element
- secondary_color: hex complementary color for text and accents
- accent_color: hex highlight color for CTAs and key metrics
- background_tone: "light", "dark", "warm", or "cool"
- style: "modern", "classic", "industrial", "luxury", or "mixed"
- mood: "professional", "inviting", "bold", or "elegant"
- font_suggestion: a Google Font name that matches the property
- description: 1-2 sentences describing the design direction

Return ONLY valid JSON, no markdown fences.`

	// Build content parts
	parts := []*genai.Part{
		{Text: prompt},
	}

	limit := 4
	if len(photoPaths) < limit {
		limit = len(photoPaths)
	}

	for _, path := range photoPaths[:limit] {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		parts = append(parts, &genai.Part{
			InlineData: &genai.Blob{
				MIMEType: detectMIME(path),
				Data:     data,
			},
		})
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash",
		[]*genai.Content{{Parts: parts}},
		&genai.GenerateContentConfig{
			Temperature: genai.Ptr(float32(0.3)),
		},
	)
	if err != nil {
		return defaultAesthetic(), fmt.Errorf("vision analysis: %w", err)
	}

	text := result.Text()

	var aesthetic PropertyAesthetic
	if err := json.Unmarshal([]byte(text), &aesthetic); err != nil {
		return defaultAesthetic(), nil
	}

	return &aesthetic, nil
}

func defaultAesthetic() *PropertyAesthetic {
	return &PropertyAesthetic{
		PrimaryColor:   "#1a1a2e",
		SecondaryColor: "#4a4a5a",
		AccentColor:    "#2563eb",
		BackgroundTone: "light",
		Style:          "professional",
		Mood:           "professional",
		FontSuggestion: "Inter",
		Description:    "Clean, professional design with neutral tones.",
	}
}

func detectMIME(path string) string {
	switch {
	case len(path) > 4 && path[len(path)-4:] == ".png":
		return "image/png"
	case len(path) > 5 && path[len(path)-5:] == ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}
