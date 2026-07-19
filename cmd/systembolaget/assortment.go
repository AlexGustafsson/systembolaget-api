package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/alexgustafsson/systembolaget-api/v4/systembolaget"
	"github.com/urfave/cli/v2"
)

type JSONStream struct {
	out      io.Writer
	previous bool
}

func NewJSONStream(out io.Writer) *JSONStream {
	out.Write([]byte("[\n  "))
	return &JSONStream{
		out:      out,
		previous: false,
	}
}

func (s *JSONStream) Write(v any) error {
	if s.previous {
		if _, err := s.out.Write([]byte(",\n  ")); err != nil {
			return err
		}
	}

	buffer, err := json.MarshalIndent(v, "  ", "  ")
	if err != nil {
		return err
	}

	s.previous = true
	_, err = s.out.Write(bytes.TrimSpace(buffer))
	return err
}

func (s *JSONStream) Close() error {
	_, err := s.out.Write([]byte("\n]"))
	return err
}

func ActionAssortment(ctx *cli.Context) error {
	log := configureLogging(ctx)

	apiKey, err := getAPIKey(ctx, log)
	if err != nil {
		return err
	}

	delayBetweenPages := ctx.Duration("page-delay")
	limit := ctx.Int("limit")

	client := systembolaget.NewClient(apiKey)

	options := &systembolaget.SearchOptions{}

	if property, ok := ctx.Value("sort-by").(string); ok && property != "" {
		options.SortBy = systembolaget.SortProperty(property)
	}

	if ctx.String("sort") == "ascending" {
		options.SortDirection = systembolaget.SortDirectionAscending
	} else if ctx.String("sort") == "descending" {
		options.SortDirection = systembolaget.SortDirectionDescending
	}

	filters := []systembolaget.SearchFilter{}

	if store := ctx.String("store"); store != "" {
		filters = append(filters, systembolaget.FilterByStore(store))
	}

	if query := ctx.String("query"); query != "" {
		filters = append(filters, systembolaget.FilterByQuery(query))
	}

	if r, ok := ctx.Value("taste-clock-body").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByTasteClockBody(r.Minimum, r.Maximum))
	}

	if r, ok := ctx.Value("taste-clock-bitterness").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByTasteClockBitterness(r.Minimum, r.Maximum))
	}

	if r, ok := ctx.Value("taste-clock-sweetness").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByTasteClockSweetness(r.Minimum, r.Maximum))
	}

	if r, ok := ctx.Value("taste-clock-smokiness").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByTasteClockSmokiness(r.Minimum, r.Maximum))
	}

	for _, vintage := range ctx.IntSlice("vintage") {
		filters = append(filters, systembolaget.FilterByVintage(vintage))
	}

	if r, ok := ctx.Value("product-launch").(*Range[string]); ok && r != nil {
		min, err := time.Parse("2006-01-02", r.Minimum)
		if err != nil {
			return err
		}
		max, err := time.Parse("2006-01-02", r.Maximum)
		if err != nil {
			return err
		}
		filters = append(filters, systembolaget.FilterByProductLaunch(min, max))
	}

	if r, ok := ctx.Value("alcohol-percentage").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByAlcoholPercentage(r.Minimum, r.Maximum))
	}

	if r, ok := ctx.Value("sugar-content").(*Range[float32]); ok && r != nil {
		filters = append(filters, systembolaget.FilterBySugarContent(r.Minimum, r.Maximum))
	}

	for _, grapes := range ctx.StringSlice("grapes") {
		filters = append(filters, systembolaget.FilterByGrapes(grapes))
	}

	for _, match := range ctx.StringSlice("match") {
		filters = append(filters, systembolaget.FilterByMatch(match))
	}

	for _, assortment := range ctx.StringSlice("assortment") {
		filters = append(filters, systembolaget.FilterByAssortment(assortment))
	}

	for _, seal := range ctx.StringSlice("seal") {
		filters = append(filters, systembolaget.FilterBySeal(seal))
	}

	if r, ok := ctx.Value("volume").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByVolume(r.Minimum, r.Maximum))
	}

	if category := ctx.String("packaging-category"); category != "" {
		subcategory := ctx.String("packaging-subcategory")
		filters = append(filters, systembolaget.FilterByPackaging(category, subcategory))
	}

	if r, ok := ctx.Value("price").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByPrice(r.Minimum, r.Maximum))
	}

	for _, origin := range ctx.StringSlice("origin") {
		filters = append(filters, systembolaget.FilterByOrigin(origin))
	}

	if category := ctx.String("category"); category != "" {
		subcategory := ctx.String("subcategory")
		subsubcategories := ctx.StringSlice("subsubcategory")
		if len(subsubcategories) == 0 {
			filters = append(filters, systembolaget.FilterByCategory(category, subcategory, ""))
		} else {
			for _, subsubcategory := range subsubcategories {
				filters = append(filters, systembolaget.FilterByCategory(category, subcategory, subsubcategory))
			}
		}
	}

	log.Debug("Retrieving cursor")
	cursor := client.SearchWithCursor(options, filters...)

	var output io.Writer
	if ctx.String("output") == "" {
		output = os.Stdout
	} else {
		file, err := os.OpenFile(ctx.String("output"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Error("Failed to open output file", slog.Any("error", err))
			return err
		}
		defer file.Close()
		output = file
	}

	runCtx, cancel := context.WithCancel(ctx.Context)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	caught := 0
	go func() {
		for range signals {
			caught++
			if caught == 1 {
				slog.Info("Caught signal, exiting gracefully")
				cancel()
			} else {
				slog.Info("Caught signal, exiting now")
				os.Exit(1)
			}
		}
	}()

	out := NewJSONStream(output)
	defer out.Close()

	log.Debug("Fetching results")
	fetchedResults := 0
	for cursor.Next(systembolaget.SetLogger(runCtx, log), delayBetweenPages) {
		if err := out.Write(cursor.At()); err != nil {
			log.Error("Failed to write result", slog.Any("error", err))
			return err
		}

		fetchedResults++
		if limit > 0 && fetchedResults == limit {
			break
		}
	}

	totalResults := -1
	if currentPage := cursor.CurrentPage(); currentPage != nil {
		totalResults = currentPage.Metadata.FullAssortmentDocumentCount
	}

	if err := cursor.Error(); err != nil {
		log.Error("Failed to fetch next item", slog.Any("error", err), slog.Int("results", fetchedResults), slog.Int("resultsLimit", limit), slog.Int("totalResults", totalResults))
		return err
	}

	log.Debug("All results have been processed", slog.Int("results", fetchedResults), slog.Int("resultsLimit", limit), slog.Int("totalResults", totalResults))
	return nil
}
