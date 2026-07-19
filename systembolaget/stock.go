package systembolaget

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

type StockStatus struct {
	ProductID           string `json:"productId"`
	StoreID             string `json:"storeId"`
	Shelf               string `json:"shelf"`
	Stock               int    `json:"stock"`
	IsInStoreAssortment bool   `json:"isInStoreAssortment"`
}

func (c *Client) GetStockStatus(ctx context.Context, storeID string, productID string) (*StockStatus, error) {
	log := GetLogger(ctx)

	u := &url.URL{
		Scheme: "https",
		Host:   "api-extern.systembolaget.se",
		Path:   fmt.Sprintf("/sb-api-ecommerce/v1/stockbalance/store/%s/%s/", storeID, productID),
	}

	log = log.With(slog.String("url", u.String()))

	header := http.Header{}
	header.Set("Origin", "https://www.systembolaget.se")
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

	var status StockStatus
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&status); err != nil {
		log.Error("Failed to decode body", slog.Any("error", err))
		return nil, err
	}

	return &status, nil
}
