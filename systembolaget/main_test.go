package systembolaget

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
)

// apiKey is populated with the API key to use for integration tests.
// Only set if running full test.
var apiKey string

// TestMain handles fetching of API credentials if running a full test.
func TestMain(m *testing.M) {
	flag.Parse()

	if !testing.Short() {
		apiKey = os.Getenv("API_KEY")
		if apiKey == "" {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			var err error
			apiKey, err = GetAPIKey(ctx)
			cancel()

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	}

	os.Exit(m.Run())
}
