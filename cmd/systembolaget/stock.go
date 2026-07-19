package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

func ActionStock(ctx context.Context, cmd *cli.Command) error {
	storeID := cmd.Args().Get(0)
	productID := cmd.Args().Get(1)
	if storeID == "" || productID == "" {
		return cli.Exit("Missing store and/or product id", 1)
	}

	log := configureLogging(cmd)

	client, err := getClient(ctx, cmd, log)
	if err != nil {
		return err
	}

	var output io.Writer
	if cmd.String("output") == "" {
		output = os.Stdout
	} else {
		file, err := os.OpenFile(cmd.String("output"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Error("Failed to open output file", slog.Any("error", err))
			return err
		}
		defer file.Close()
		output = file
	}

	stores, err := client.GetStockStatus(ctx, storeID, productID)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stores)
}
