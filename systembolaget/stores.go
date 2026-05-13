package systembolaget

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
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
	IsAgent           bool                `json:"isAgent"`
	IsBlocked         bool                `json:"isBlocked"`
	BlockedText       string              `json:"blockedText"`
	IsSvanenCertified bool                `json:"isSvanenCertified"`
	IsOpen            bool                `json:"isOpen"`
	IsTastingStore    bool                `json:"isTastingStore"`
	OpeningHours      []StoreOpeningHours `json:"openingHours"`
	Position          *StorePosition      `json:"position"`
}

type StoreOpeningHours struct {
	Date     string `json:"date"`
	OpenFrom string `json:"openFrom"`
	OpenTo   string `json:"openTo"`
	Reason   string `json:"reason"`
}

type StorePosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Stores fetches available stores.
func (c *Client) Stores(ctx context.Context) ([]Store, error) {
	return c.SearchStores(ctx, "", true)
}

// SearchStores searches for available stores.
// Query typically matches both name and location.
func (c *Client) SearchStores(ctx context.Context, query string, includePredictions bool) ([]Store, error) {
	log := GetLogger(ctx)

	queryParams := url.Values{}
	queryParams.Set("includePredictions", strconv.FormatBool(includePredictions))
	if query != "" {
		queryParams.Set("q", query)
	}

	u := &url.URL{
		Scheme:   "https",
		Host:     "api-extern.systembolaget.se",
		Path:     "/sb-api-ecommerce/v1/sitesearch/site",
		RawQuery: queryParams.Encode(),
	}

	log = log.With(slog.String("url", u.String()))

	header := http.Header{}
	header.Set("Origin", "https://www.systembolaget.se")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Pragma", "no-cache")
	header.Set("Accept", "application/json")
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

	var result struct {
		Stores []Store `json:"siteSearchResults"`
	}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&result); err != nil {
		log.Error("Failed to decode body", slog.Any("error", err))
		return nil, err
	}

	return result.Stores, nil
}
