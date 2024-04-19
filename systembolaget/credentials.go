package systembolaget

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

var appSettingsPathRegex = regexp.MustCompile(`"appSettingsFilePath":"([^"]+)"`)
var appSettingsRegex = regexp.MustCompile(`(?s)Object.freeze\((\{.*\})\)`)

// GetAPIKey returns the API credentials used by the Systembolaget
// frontend.
func GetAPIKey(ctx context.Context) (string, error) {
	appSettings, err := getAppSettings(ctx)
	if err != nil {
		return "", err
	}

	apiKeyValue, ok := appSettings["ocpApimSubscriptionKey"]
	if !ok {
		return "", fmt.Errorf("malformed app settings")
	}
	apiKey, ok := apiKeyValue.(string)
	if !ok {
		return "", fmt.Errorf("malformed app settings")
	}

	return apiKey, nil
}

func getAppSettings(ctx context.Context) (map[string]any, error) {
	log := GetLogger(ctx)

	log.Debug("Fetching app settings script path")
	appSettingsScriptPath, err := getAppSettingsScriptPath(ctx)
	if err != nil {
		return nil, err
	}
	log = slog.With(slog.String("appSettingsScriptPath", appSettingsScriptPath))

	log.Debug("Fetching app settings")
	res, err := http.DefaultClient.Get("https://www.systembolaget.se/" + appSettingsScriptPath)
	if err != nil {
		log.Error("Request failed", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error("Got unexpected status code", slog.Int("statusCode", res.StatusCode), slog.String("status", res.Status))
		return nil, fmt.Errorf("unexpected status code: %d - %s", res.StatusCode, res.Status)
	}

	source, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to read body")
		return nil, err
	}

	match := appSettingsRegex.FindSubmatch(source)
	if match == nil {
		log.Error("Unable to find script path")
		return nil, fmt.Errorf("unable to identify appsettings script path")
	}

	var settings map[string]any
	if err := json.Unmarshal(match[1], &settings); err != nil {
		log.Error("Failed to decode body", slog.Any("error", err))
		return nil, err
	}

	return settings, nil
}

func getAppSettingsScriptPath(ctx context.Context) (string, error) {
	log := GetLogger(ctx)

	log.Debug("Fetching systembolaget.se")
	res, err := http.DefaultClient.Get("https://www.systembolaget.se")
	if err != nil {
		slog.Error("Request failed", slog.Any("error", err))
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("Got unexpected status code", slog.Int("statusCode", res.StatusCode), slog.String("status", res.Status))
		return "", fmt.Errorf("unexpected status code: %d - %s", res.StatusCode, res.Status)
	}

	source, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to read body")
		return "", err
	}

	match := appSettingsPathRegex.FindSubmatch(source)
	if match == nil {
		slog.Error("Unable to find script path")
		return "", fmt.Errorf("unable to identify appsettings script path")
	}

	// The path seems to always be prefixed with "~/"
	return strings.TrimPrefix(string(match[1]), "~/"), nil
}
