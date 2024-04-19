package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/alexgustafsson/systembolaget-api/v3/systembolaget"
	"github.com/urfave/cli/v2"
)

func ActionStores(ctx *cli.Context) error {
	log := configureLogging(ctx)
	ctxWithLogging := systembolaget.SetLogger(ctx.Context, log)

	apiKey, err := getAPIKey(ctx, log)
	if err != nil {
		return err
	}

	client := systembolaget.NewClient(apiKey)

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

	stores, err := client.Stores(ctxWithLogging)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stores)
}
