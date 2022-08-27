package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alexgustafsson/systembolaget-api/v2/systembolaget"
	"github.com/urfave/cli/v2"
)

func ActionStores(ctx *cli.Context) error {
	apiKey, err := systembolaget.GetAPIKey(ctx.Context)
	if err != nil {
		return fmt.Errorf("failed to get API key, please specify one")
	}

	client := systembolaget.NewClient(apiKey)

	stores, err := client.Stores(ctx.Context)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stores)
}
