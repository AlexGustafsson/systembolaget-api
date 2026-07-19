package systembolaget

import (
	"context"
	"encoding/json"
	"fmt"
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

// GetStockStatus fetches the stock status of a product in a specific store.
func (c *AuthenticatedClient) GetStockStatus(ctx context.Context, storeID string, productID string) (*StockStatus, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "api-extern.systembolaget.se",
		Path:   fmt.Sprintf("/sb-api-ecommerce/v1/stockbalance/store/%s/%s/", storeID, productID),
	}

	header := make(http.Header)
	header.Set("Origin", "https://www.systembolaget.se")
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

	var status StockStatus
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}
