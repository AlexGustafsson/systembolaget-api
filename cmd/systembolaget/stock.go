package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"
)

func ActionStock(ctx *cli.Context) error {
	storeID := ctx.Args().Get(0)
	productID := ctx.Args().Get(1)
	if storeID == "" || productID == "" {
		return cli.Exit("Missing store and/or product id", 1)
	}

	log := configureLogging(ctx)

	client, err := getClient(ctx, log)
	if err != nil {
		return err
	}

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

	stores, err := client.GetStockStatus(runCtx, storeID, productID)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stores)
}
