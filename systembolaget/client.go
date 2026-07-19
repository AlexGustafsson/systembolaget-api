package systembolaget

import (
	"net/http"
)

var DefaultClient = &Client{
	Client: http.DefaultClient,
}

// Client is a Systembolaget API client for unauthenticated methods.
type Client struct {
	Client    *http.Client
	UserAgent string
}

// AuthenticatedClient is a Systembolaget API client for authenticated methods.
//
// An [AuthenticatedClient] is typically retrieved by calling
// [Client.GetAuthenticatedClient].
type AuthenticatedClient struct {
	APIKey    string
	Client    *http.Client
	UserAgent string
}
