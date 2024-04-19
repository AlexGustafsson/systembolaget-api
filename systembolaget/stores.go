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
	ID         string `json:"siteId"`
	Alias      string `json:"alias"`
	Address    string `json:"address"`
	PostalCode string `json:"postalCode"`
	City       string `json:"city"`
	Phone      string `json:"phone"`
	County     string `json:"county"`
	Position   struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	OpeningHours []StoreOpeningHours `json:"openingHours"`
}

type StoreOpeningHours struct {
	// Date is the date the entry is for.
	Date string `json:"date"`
	// OpenFrom specifies the hour that the store opens.
	OpenFrom string `json:"openFrom"`
	// OpenTo specifies the hour that the store closes.
	OpenTo string `json:"openTo"`
	// Reason optionally specifies the reason for deviating opening hours.
	Reason string `json:"reason"`
}

// Stores fetches available stores.
func (c *Client) Stores(ctx context.Context) ([]Store, error) {
	log := GetLogger(ctx)

	u := &url.URL{
		Scheme: "https",
		Host:   "api-extern.systembolaget.se",
		Path:   "/site/V2/Store",
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

	var result []Store
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&result); err != nil {
		log.Error("Failed to decode body", slog.Any("error", err))
		return nil, err
	}

	return result, nil
}
