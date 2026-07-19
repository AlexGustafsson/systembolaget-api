package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alexgustafsson/systembolaget-api/v4/systembolaget"
	"github.com/urfave/cli/v2"
)

func configureLogging(ctx *cli.Context) *slog.Logger {
	verbose := ctx.Bool("verbose")
	if verbose {
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}

func getAPIKey(ctx *cli.Context, log *slog.Logger) (string, error) {
	if apiKey := ctx.String("api-key"); apiKey != "" {
		return apiKey, nil
	}

	log.Debug("Fetching API key")
	apiKey, err := systembolaget.GetAPIKey(systembolaget.SetLogger(ctx.Context, log))
	if err != nil {
		return "", fmt.Errorf("failed to get API key, please specify one")
	}

	return apiKey, nil
}
