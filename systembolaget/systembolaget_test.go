package systembolaget

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

// testClient is populated with an authenticated client to use for integration
// tests.
//
// Only set if running full tests.
var testClient *AuthenticatedClient

func TestMain(m *testing.M) {
	if !testing.Short() {
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			client, err := DefaultClient.GetAuthenticatedClient(ctx)
			cancel()
			if err != nil {
				panic(err)
			}

			testClient = client
		} else {
			testClient = &AuthenticatedClient{
				APIKey: apiKey,
				Client: http.DefaultClient,
			}
		}
	}

	os.Exit(m.Run())
}
