package systembolaget

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStores(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := NewClient(apiKey)

	stores, err := client.Stores(context.TODO())
	require.NoError(t, err)

	for _, store := range stores {
		if assert.NotNil(t, store) {
			assert.NotEmpty(t, store.ID)
		}
	}
}
