package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/alexgustafsson/systembolaget-api/v4/systembolaget"
	"github.com/urfave/cli/v3"
)

func ActionAssortment(ctx context.Context, cmd *cli.Command) error {
	log := getLogger(cmd)

	client, err := getClient(ctx, cmd, log)
	if err != nil {
		return err
	}

	delayBetweenPages := cmd.Duration("page-delay")
	limit := cmd.Int("limit")

	options := &systembolaget.SearchOptions{}

	if property, ok := ctx.Value("sort-by").(string); ok && property != "" {
		options.SortBy = systembolaget.SortProperty(property)
	}

	if cmd.String("sort") == "ascending" {
		options.SortDirection = systembolaget.SortDirectionAscending
	} else if cmd.String("sort") == "descending" {
		options.SortDirection = systembolaget.SortDirectionDescending
	}

	filters := []systembolaget.SearchFilter{}

	if store := cmd.String("store"); store != "" {
		filters = append(filters, systembolaget.FilterByStore(store))
	}

	if query := cmd.String("query"); query != "" {
		filters = append(filters, systembolaget.FilterByQuery(query))
	}

	if r, ok := cmd.Value("taste-clock-body").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByTasteClockBody(r.Minimum, r.Maximum))
	}

	if r, ok := cmd.Value("taste-clock-bitterness").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByTasteClockBitterness(r.Minimum, r.Maximum))
	}

	if r, ok := cmd.Value("taste-clock-sweetness").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByTasteClockSweetness(r.Minimum, r.Maximum))
	}

	if r, ok := cmd.Value("taste-clock-smokiness").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByTasteClockSmokiness(r.Minimum, r.Maximum))
	}

	for _, vintage := range cmd.IntSlice("vintage") {
		filters = append(filters, systembolaget.FilterByVintage(vintage))
	}

	if r, ok := cmd.Value("product-launch").(*Range[string]); ok && r != nil {
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

	if r, ok := cmd.Value("alcohol-percentage").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByAlcoholPercentage(r.Minimum, r.Maximum))
	}

	if r, ok := cmd.Value("sugar-content").(*Range[float32]); ok && r != nil {
		filters = append(filters, systembolaget.FilterBySugarContent(r.Minimum, r.Maximum))
	}

	for _, grapes := range cmd.StringSlice("grapes") {
		filters = append(filters, systembolaget.FilterByGrapes(grapes))
	}

	for _, match := range cmd.StringSlice("match") {
		filters = append(filters, systembolaget.FilterByMatch(match))
	}

	for _, assortment := range cmd.StringSlice("assortment") {
		filters = append(filters, systembolaget.FilterByAssortment(assortment))
	}

	for _, seal := range cmd.StringSlice("seal") {
		filters = append(filters, systembolaget.FilterBySeal(seal))
	}

	if r, ok := cmd.Value("volume").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByVolume(r.Minimum, r.Maximum))
	}

	if category := cmd.String("packaging-category"); category != "" {
		subcategory := cmd.String("packaging-subcategory")
		filters = append(filters, systembolaget.FilterByPackaging(category, subcategory))
	}

	if r, ok := cmd.Value("price").(*Range[int]); ok && r != nil {
		filters = append(filters, systembolaget.FilterByPrice(r.Minimum, r.Maximum))
	}

	for _, origin := range cmd.StringSlice("origin") {
		filters = append(filters, systembolaget.FilterByOrigin(origin))
	}

	if category := cmd.String("category"); category != "" {
		subcategory := cmd.String("subcategory")
		subsubcategories := cmd.StringSlice("subsubcategory")
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

	encoder := json.NewEncoder(os.Stdout)

	log.Debug("Fetching results")
	fetchedResults := 0
	for cursor.Next(ctx, delayBetweenPages) {
		if err := encoder.Encode(cursor.At()); err != nil {
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
