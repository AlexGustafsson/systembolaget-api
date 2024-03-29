package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/alexgustafsson/systembolaget-api/v2/systembolaget"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func ActionStores(ctx *cli.Context) error {
	log, err := configureLogging(ctx)
	if err != nil {
		return err
	}
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
			log.Fatal("Failed to open output file", zap.Error(err))
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
