package systembolaget

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

var apiTokenRegex = regexp.MustCompile(`NEXT_PUBLIC_OCP_APIM_KEY:"([^"]+)"`)

// <script src="https://sb-web-ecommerce-app.azureedge.net/_next/static/chunks/pages/_app-b8fd056cfd040021.js" defer=""></script>
var appBundlePathRegex = regexp.MustCompile(`<script src="([^"]+app-[^"]+.js)"`)

// GetAPIKey returns the API credentials used by the Systembolaget
// frontend.
func GetAPIKey(ctx context.Context) (string, error) {
	log := GetLogger(ctx)

	log.Debug("Fetching app settings script path")
	appBundlePath, err := getAppBundlePath(ctx)
	if err != nil {
		return "", err
	}
	log = slog.With(slog.String("appBundlePath", appBundlePath))

	log.Debug("Fetching app settings")
	res, err := http.DefaultClient.Get(appBundlePath)
	if err != nil {
		log.Error("Request failed", slog.Any("error", err))
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error("Got unexpected status code", slog.Int("statusCode", res.StatusCode), slog.String("status", res.Status))
		return "", fmt.Errorf("unexpected status code: %d - %s", res.StatusCode, res.Status)
	}

	source, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to read body")
		return "", err
	}

	match := apiTokenRegex.FindSubmatch(source)
	if match == nil {
		log.Error("Unable to find API token")
		return "", fmt.Errorf("unable to identify API token")
	}

	return string(match[1]), nil
}

func getAppBundlePath(ctx context.Context) (string, error) {
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

	match := appBundlePathRegex.FindSubmatch(source)
	if match == nil {
		slog.Error("Unable to find script path")
		return "", fmt.Errorf("unable to identify appsettings script path")
	}

	return string(match[1]), nil
}
