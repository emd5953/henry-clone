package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/henry-clone/internal/api"
	"github.com/henry-clone/internal/deck"
	"github.com/henry-clone/internal/enrichment"
	"github.com/henry-clone/internal/figma"
	"github.com/henry-clone/internal/llm"
)

func main() {
	// Wire up dependencies
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	var narrator deck.Narrator
	if apiKey != "" {
		narrator = llm.NewAnthropicNarrator(apiKey)
		log.Println("Using Claude for narrative generation")
	} else {
		narrator = llm.NewStubNarrator()
		log.Println("No ANTHROPIC_API_KEY set — using stub narrator")
	}

	builder := deck.NewBuilder(narrator)

	// Enrichment providers — stubs for now
	compsProvider := enrichment.NewStubCompsProvider()
	marketProvider := enrichment.NewStubMarketDataProvider()
	geoProvider := enrichment.NewStubGeoProvider()

	handler := api.NewHandler(api.HandlerConfig{
		Builder:  builder,
		Narrator: narrator,
		Comps:    compsProvider,
		Market:   marketProvider,
		Geo:      geoProvider,
	})

	// Figma integration (optional — needs FIGMA_TOKEN env var)
	var figmaHandler *api.FigmaHandler
	if figmaToken := os.Getenv("FIGMA_TOKEN"); figmaToken != "" {
		bridge := figma.NewBridge(figmaToken)
		figmaHandler = api.NewFigmaHandler(bridge, handler)
		log.Println("Figma integration enabled")
	} else {
		log.Println("No FIGMA_TOKEN set — Figma integration disabled")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	// Deal CRUD
	r.Post("/api/deals", handler.CreateDeal)
	r.Get("/api/deals", handler.ListDeals)
	r.Get("/api/deals/{dealID}", handler.GetDeal)
	r.Get("/api/deals/{dealID}/deck", handler.GetDeck)
	r.Get("/api/deals/{dealID}/deck.pdf", handler.GetDeckPDF)
	r.Get("/api/deals/{dealID}/sections", handler.GetSections)
	r.Put("/api/deals/{dealID}/sections/{sectionIdx}", handler.UpdateSection)

	// QC Review workflow
	r.Get("/api/reviews", handler.GetReviewQueue)
	r.Post("/api/deals/{dealID}/review/start", handler.StartReview)
	r.Post("/api/deals/{dealID}/review/complete", handler.CompleteReview)
	r.Post("/api/deals/{dealID}/review/edit", handler.ReviewEdit)

	// Figma integration
	if figmaHandler != nil {
		r.Post("/api/deals/{dealID}/figma/link", figmaHandler.LinkFigmaFile)
		r.Get("/api/deals/{dealID}/figma", figmaHandler.GetFigmaFile)
		r.Get("/api/deals/{dealID}/figma/export", figmaHandler.ExportFigmaPDF)
		r.Post("/api/deals/{dealID}/figma/comment", figmaHandler.PostFigmaComment)
	}

	// Serve React frontend (static files)
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fsHandler := http.FileServer(http.Dir("frontend/dist"))
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		if _, err := http.Dir("frontend/dist").Open(path); err != nil {
			http.ServeFile(w, r, "frontend/dist/index.html")
			return
		}
		fsHandler.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Henry Clone running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
