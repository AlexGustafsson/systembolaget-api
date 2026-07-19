package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexgustafsson/systembolaget-api/v4/systembolaget"
	"github.com/urfave/cli/v3"
)

func getLogger(cmd *cli.Command) *slog.Logger {
	verbose := cmd.Bool("verbose")
	if verbose {
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}

func getClient(ctx context.Context, cmd *cli.Command, log *slog.Logger) (*systembolaget.AuthenticatedClient, error) {
	if apiKey := cmd.String("api-key"); apiKey != "" {
		return &systembolaget.AuthenticatedClient{
			APIKey: apiKey,
			Client: http.DefaultClient,
		}, nil
	}

	log.Debug("Fetching API key")
	client, err := systembolaget.DefaultClient.GetAuthenticatedClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key, please specify one: %w", err)
	}

	return client, nil
}
