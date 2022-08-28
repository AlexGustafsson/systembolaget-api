package main

import (
	"fmt"

	"github.com/alexgustafsson/systembolaget-api/v2/systembolaget"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func configureLogging(ctx *cli.Context) (*zap.Logger, error) {
	var logConfig zap.Config
	verbose := ctx.Bool("verbose")
	if verbose {
		logConfig = zap.NewDevelopmentConfig()
	} else {
		logConfig = zap.NewProductionConfig()
	}

	logConfig.EncoderConfig.TimeKey = "time"
	logConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	return logConfig.Build()
}

func getAPIKey(ctx *cli.Context, log *zap.Logger) (string, error) {
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
