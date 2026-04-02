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
	"github.com/henry-clone/internal/llm"
)

func main() {
	// Wire up dependencies
	apiKey := os.Getenv("OPENAI_API_KEY")
	var narrator deck.Narrator
	if apiKey != "" {
		narrator = llm.NewOpenAINarrator(apiKey)
		log.Println("Using OpenAI for narrative generation")
	} else {
		narrator = llm.NewStubNarrator()
		log.Println("No OPENAI_API_KEY set — using stub narrator")
	}

	builder := deck.NewBuilder(narrator)

	// Enrichment providers — stubs for now, swap for real APIs later
	// (CoStar, Reonomy, Census, Google Maps, etc.)
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
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

	// Serve React frontend (static files)
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve from frontend/dist, fall back to index.html for SPA routing
		fsHandler := http.FileServer(http.Dir("frontend/dist"))
		// Check if file exists
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		if _, err := http.Dir("frontend/dist").Open(path); err != nil {
			// SPA fallback — serve index.html for client-side routing
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
