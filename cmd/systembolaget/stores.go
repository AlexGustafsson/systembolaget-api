package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/urfave/cli/v3"
)

func ActionStores(ctx context.Context, cmd *cli.Command) error {
	log := getLogger(cmd)

	client, err := getClient(ctx, cmd, log)
	if err != nil {
		return err
	}

	stores, err := client.SearchStores(ctx, cmd.String("search"), cmd.String("search") == "")
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(os.Stdout)
	for _, store := range stores {
		err := encoder.Encode(store)
		if err != nil {
			return err
		}
	}
	return nil
}
