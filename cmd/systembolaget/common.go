package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
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

func getClient(ctx *cli.Context, log *slog.Logger) (*systembolaget.AuthenticatedClient, error) {
	if apiKey := ctx.String("api-key"); apiKey != "" {
		return &systembolaget.AuthenticatedClient{
			APIKey: apiKey,
			Client: http.DefaultClient,
		}, nil
	}

	log.Debug("Fetching API key")
	client, err := systembolaget.DefaultClient.GetAuthenticatedClient(ctx.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key, please specify one: %w", err)
	}

	return client, nil
}

type JSONStream struct {
	out      io.Writer
	previous bool
}

func NewJSONStream(out io.Writer) *JSONStream {
	out.Write([]byte("[\n  "))
	return &JSONStream{
		out:      out,
		previous: false,
	}
}

func (s *JSONStream) Write(v any) error {
	if s.previous {
		if _, err := s.out.Write([]byte(",\n  ")); err != nil {
			return err
		}
	}

	buffer, err := json.MarshalIndent(v, "  ", "  ")
	if err != nil {
		return err
	}

	s.previous = true
	_, err = s.out.Write(bytes.TrimSpace(buffer))
	return err
}

func (s *JSONStream) Close() error {
	_, err := s.out.Write([]byte("\n]"))
	return err
}
