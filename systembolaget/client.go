package systembolaget

import (
	"net/http"
	"net/url"
)

// Client is a Systembolaget API client.
type Client struct {
	apiKey     string
	httpClient *http.Client
	userAgent  string
}

// Option is a Client option.
type Option func(*Client)

// WithUserAgent specifies a user agent string to include in requests to the
// API.
func WithUserAgent(agent string) Option {
	return func(client *Client) {
		client.userAgent = agent
	}
}

// WithProxy specifies a proxy to use for requests.
func WithProxy(proxy string) Option {
	return func(client *Client) {
		client.httpClient.Transport = &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse(proxy)
			},
		}
	}
}

// NewClient creates a new client.
func NewClient(apiKey string, options ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(c)
	}

	return c
}
