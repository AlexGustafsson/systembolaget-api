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

type StockBalance struct {
	ProductNumber string `json:"productId"`
	StoreID       string `json:"storeId"`
	StockLevel    int    `json:"stock"`
	Shelf         string `json:"shelf"`
}

func (c *Client) StockBalance(ctx context.Context, storeID string, productNumber string) (*StockBalance, error) {
	log := GetLogger(ctx)

	u := &url.URL{
		Scheme: "https",
		Host:   "api-extern.systembolaget.se",
		Path:   fmt.Sprintf("/sb-api-ecommerce/v1/stockbalance/store/%s/%s/", storeID, productNumber),
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
		log.Error("Request failed", slog.Any("error", err))
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

	var stockBalance StockBalance
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&stockBalance); err != nil {
		log.Error("Failed to decode body", slog.Any("error", err))
		return nil, err
	}

	return &stockBalance, nil
}
