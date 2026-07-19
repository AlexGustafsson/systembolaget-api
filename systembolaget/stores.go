package systembolaget

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Store represents a Systembolaget store.
type Store struct {
	SiteID            string              `json:"siteId"`
	Alias             string              `json:"alias"`
	StreetAddress     string              `json:"streetAddress"`
	DisplayName       string              `json:"displayName"`
	City              string              `json:"city"`
	County            string              `json:"county"`
	IsAgent           bool                `json:"isAgent"`
	IsBlocked         bool                `json:"isBlocked"`
	BlockedText       string              `json:"blockedText"`
	IsSvanenCertified bool                `json:"isSvanenCertified"`
	IsOpen            bool                `json:"isOpen"`
	IsTastingStore    bool                `json:"isTastingStore"`
	OpeningHours      []StoreOpeningHours `json:"openingHours"`
	Position          *StorePosition      `json:"position"`
	// Ignored fields:
	// PostalCode string `json:"postalCode"` // Always null
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

// GetStore fetches all available stores.
func (c *AuthenticatedClient) GetStores(ctx context.Context) ([]Store, error) {
	return c.SearchStores(ctx, "", true)
}

// SearchStores searches for available using a query.
// Query typically matches both name and location.
func (c *AuthenticatedClient) SearchStores(ctx context.Context, query string, includePredictions bool) ([]Store, error) {
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

	header := http.Header{}
	header.Set("Origin", "https://www.systembolaget.se")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Pragma", "no-cache")
	header.Set("Accept", "application/json")
	header.Set("Cache-Control", "no-cache")
	header.Set("Ocp-Apim-Subscription-Key", c.APIKey)

	if c.UserAgent != "" {
		header.Set("User-Agent", c.UserAgent)
	}

	req := (&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: header,
	}).Clone(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d - %s", res.StatusCode, res.Status)
	}

	var result struct {
		Stores []Store `json:"siteSearchResults"`
	}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return result.Stores, nil
}
