package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"os/signal"

	"github.com/alexgustafsson/systembolaget-api/v4/systembolaget"
	"github.com/urfave/cli/v2"
)

func ActionStock(ctx *cli.Context) error {
	log := configureLogging(ctx)

	apiKey, err := getAPIKey(ctx, log)
	if err != nil {
		return err
	}

	storeID := ctx.String("store-id")
	if storeID == "" {
		log.Error("Store ID is required")
		return cli.Exit("store-id is required", 1)
	}

	productNumber := ctx.String("product-number")
	if productNumber == "" {
		log.Error("Product number is required")
		return cli.Exit("product-number is required", 1)
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

	stockBalance, err := client.StockBalance(systembolaget.SetLogger(runCtx, log), storeID, productNumber)
	if err != nil {
		log.Error("Failed to retrieve stock balance", slog.Any("error", err))
		return err
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(stockBalance); err != nil {
		log.Error("Failed to encode stock balance", slog.Any("error", err))
		return err
	}

	log.Debug("Stock balance retrieved successfully", slog.String("storeID", storeID), slog.String("productNumber", productNumber))
	return nil
}
