package systembolaget

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// Store represents a Systembolaget store.
type Store struct {
	SiteID        string `json:"siteId"`
	Alias         string `json:"alias"`
	StreetAddress string `json:"streetAddress"`
	DisplayName   string `json:"displayName"`
	City          string `json:"city"`
	County        string `json:"county"`
	// NOTE: always null, exclude for now
	// PostalCode
	IsAgent        bool                `json:"isAgent"`
	IsBlocked      bool                `json:"isBlocked"`
	BlockedText    string              `json:"blockedText"`
	IsOpen         bool                `json:"isOpen"`
	IsTastingStore bool                `json:"isTastingStore"`
	OpeningHours   []StoreOpeningHours `json:"openingHours"`
}

type StoreOpeningHours struct {
	Date     string `json:"date"`
	OpenFrom string `json:"openFrom"`
	OpenTo   string `json:"openTo"`
	Reason   string `json:"reason"`
}

// Stores fetches available stores.
func (c *Client) Stores(ctx context.Context) ([]Store, error) {
	log := GetLogger(ctx)

	query := url.Values{}
	query.Set("includePredictions", "false")

	u := &url.URL{
		Scheme:   "https",
		Host:     "api-extern.systembolaget.se",
		Path:     "/sb-api-ecommerce/v1/sitesearch/site",
		RawQuery: query.Encode(),
	}

	log = log.With(slog.String("url", u.String()))

	header := http.Header{}
	header.Set("Origin", "https://www.systembolaget.se")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Pragma", "no-cache")
	header.Set("Accept", "application/json")
	header.Set("Accept-Encoding", "gzip")
	header.Set("Cache-Control", "no-cache")
	header.Set("Ocp-Apim-Subscription-Key", c.apiKey)

	if c.userAgent != "" {
		header.Set("User-Agent", c.userAgent)
	}

	req := (&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: header,
	}).Clone(ctx)

	log.Debug("Performing request")
	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("Request failed")
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Error("Got unexpected status code", slog.Int("statusCode", res.StatusCode), slog.String("status", res.Status))
		return nil, fmt.Errorf("unexpected status code: %d - %s", res.StatusCode, res.Status)
	}

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			log.Error("Failed to create gzip reader", slog.Any("error", err))
			return nil, err
		}
	default:
		reader = res.Body
	}
	defer reader.Close()

	var result struct {
		Stores []Store `json:"siteSearchResults"`
	}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&result); err != nil {
		log.Error("Failed to decode body", slog.Any("error", err))
		return nil, err
	}

	return result.Stores, nil
}
