package main

import (
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
