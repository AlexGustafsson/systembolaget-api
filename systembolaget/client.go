package systembolaget

import (
	"net/http"
)

var DefaultClient = &Client{
	Client: http.DefaultClient,
}

// Client is a Systembolaget API client.
type Client struct {
	Client    *http.Client
	UserAgent string
}

type AuthenticatedClient struct {
	APIKey    string
	Client    *http.Client
	UserAgent string
}
