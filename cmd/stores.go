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

func ActionStores(ctx *cli.Context) error {
	log := configureLogging(ctx)

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

	stores, err := client.SearchStores(systembolaget.SetLogger(runCtx, log), ctx.String("search"), ctx.String("search") == "")
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stores)
}
