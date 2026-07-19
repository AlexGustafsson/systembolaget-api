package systembolaget

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

var apiTokenRegex = regexp.MustCompile(`NEXT_PUBLIC_API_KEY_APIM:"([^"]+)"`)

// <script src="/_next/static/chunks/131pjojsj1pi9.js" defer=""></script>
var chunkPathRegex = regexp.MustCompile(`src="(/_next/static/chunks/[^"]+\.js)"`)

// GetAPIKey returns the API credentials used by the Systembolaget
// frontend.
func (c *Client) GetAPIKey(ctx context.Context) (string, error) {
	chunkPaths, err := c.getChunkPaths(ctx)
	if err != nil {
		return "", err
	}

	for _, chunkPath := range chunkPaths {
		key, err := c.extractAPIKey(ctx, chunkPath)
		if err == nil {
			return key, nil
		}
	}

	return "", fmt.Errorf("unable to identify API token in any script chunk")
}

// GetAuthenticatedClient calls [Client.GetAPIKey] and returns an initialized
// [AuthenticatedClient] using the key and defaults from the client.
func (c *Client) GetAuthenticatedClient(ctx context.Context) (*AuthenticatedClient, error) {
	apiKey, err := c.GetAPIKey(ctx)
	if err != nil {
		return nil, err
	}

	return &AuthenticatedClient{
		APIKey:    apiKey,
		Client:    c.Client,
		UserAgent: c.UserAgent,
	}, nil
}

func (c *Client) getChunkPaths(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.systembolaget.se", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		slog.Error("Request failed", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("Got unexpected status code", slog.Int("statusCode", res.StatusCode), slog.String("status", res.Status))
		return nil, fmt.Errorf("unexpected status code: %d - %s", res.StatusCode, res.Status)
	}

	source, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to read body")
		return nil, err
	}

	matches := chunkPathRegex.FindAllSubmatch(source, -1)
	if len(matches) == 0 {
		slog.Error("Unable to find script chunks")
		return nil, fmt.Errorf("unable to identify script chunks")
	}

	paths := make([]string, 0, len(matches))
	for _, match := range matches {
		path := string(match[1])
		if strings.HasPrefix(path, "/") {
			path = "https://www.systembolaget.se" + path
		}
		paths = append(paths, path)
	}

	return paths, nil
}

func (c *Client) extractAPIKey(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d - %s", res.StatusCode, res.Status)
	}

	source, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	match := apiTokenRegex.FindSubmatch(source)
	if match == nil {
		return "", fmt.Errorf("api token not found in chunk")
	}

	return string(match[1]), nil
}
