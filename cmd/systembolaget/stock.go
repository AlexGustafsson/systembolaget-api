package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/urfave/cli/v3"
)

func ActionStock(ctx context.Context, cmd *cli.Command) error {
	log := configureLogging(cmd)

	client, err := getClient(ctx, cmd, log)
	if err != nil {
		return err
	}

	stores, err := client.GetStockStatus(ctx, cmd.String("store-id"), cmd.String("product-id"))
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(os.Stdout)
	return encoder.Encode(stores)
}
