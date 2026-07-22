package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/alexgustafsson/systembolaget-api/v5/systembolaget"
)

func main() {
	verbose := flag.Bool("verbose", false, "verbose logs")
	apiKey := flag.String("api-key", "", "API key")
	flag.Parse()

	if *verbose {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))
	}

	slog.Warn("The proxy API is subject to change, use with caution")

	var authenticatedClient *systembolaget.AuthenticatedClient
	if *apiKey == "" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var err error
		authenticatedClient, err = systembolaget.DefaultClient.GetAuthenticatedClient(ctx)
		cancel()
		if err != nil {
			slog.Error("Failed to get an authenticated client", slog.Any("error", err))
			os.Exit(1)
		}
	} else {
		authenticatedClient = &systembolaget.AuthenticatedClient{
			APIKey: *apiKey,
			Client: http.DefaultClient,
		}
	}

	mux := http.NewServeMux()

	var failures atomic.Int32
	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		if failures.Load() > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/v1/stores/{storeId}/products/{productId}", func(w http.ResponseWriter, r *http.Request) {
		storeID := r.PathValue("storeId")
		productID := r.PathValue("productId")

		products, err := authenticatedClient.Search(r.Context(), nil, systembolaget.FilterByQuery(productID))
		if err != nil {
			failures.Add(1)
			slog.Error("Failed to search for product", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(products.Products) == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if len(products.Products) != 1 {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}

		product := products.Products[0]

		stockStatus, err := authenticatedClient.GetStockStatus(r.Context(), storeID, productID)
		if err != nil {
			failures.Add(1)
			slog.Error("Failed to get product stock status", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var response struct {
			Stock int    `json:"stock"`
			Shelf string `json:"shelf"`

			Category          string  `json:"category,omitempty"`
			Title             string  `json:"title,omitempty"`
			Subtitle          string  `json:"subtitle,omitempty"`
			Number            string  `json:"number,omitempty"`
			Country           string  `json:"country,omitempty"`
			Volume            string  `json:"volume,omitempty"`
			AlcoholPercentage float64 `json:"alcoholPercentage,omitempty"`
			Price             float64 `json:"price,omitempty"`
			Thumbnail         string  `json:"thumbnail,omitempty"`
			ImageURL          string  `json:"imageUrl,omitempty"`
		}

		response.Stock = stockStatus.Stock
		response.Shelf = stockStatus.Shelf

		if v, ok := product.Category(); ok {
			response.Category = v
		}

		if v, ok := product.Title(); ok {
			response.Title = v
		}

		if v, ok := product.Subtitle(); ok {
			response.Subtitle = v
		}

		if v, ok := product.Number(); ok {
			response.Number = v
		}

		if v, ok := product.Country(); ok {
			response.Country = v
		}

		if v, ok := product.VolumeText(); ok {
			response.Volume = v
		}

		if v, ok := product.AlcoholPercentage(); ok {
			response.AlcoholPercentage = v
		}

		if v, ok := product.Price(); ok {
			response.Price = v
		}

		if v, ok := product.Thumbnail(); ok {
			response.Thumbnail = base64.StdEncoding.EncodeToString(v)
		}

		if images, ok := product.Images(); ok {
			if len(images) > 0 {
				response.ImageURL = images[0].URL
			}
		}

		header := w.Header()
		header.Set("Content-Type", "application/json")
		header.Set("Cache-Control", "max-age=750") // Ask client to cache for 15min

		encoder := json.NewEncoder(w)
		_ = encoder.Encode(&response)
	})

	cors := http.NewCrossOriginProtection()
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: cors.Handler(mux),
	}

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Failed to serve", slog.Any("error", err))
		os.Exit(1)
	}
}
